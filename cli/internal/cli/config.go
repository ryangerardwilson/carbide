package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func inspectEnvContract() (envContractReport, error) {
	schema, err := readEnvSchema(projectConfigPath)
	if err != nil {
		return envContractReport{}, err
	}
	dotenv, envFileFound, err := readDotenv(".env")
	if err != nil {
		return envContractReport{}, err
	}

	report := envContractReport{
		schema:       schema,
		envFileFound: envFileFound,
	}
	seen := map[string]bool{}
	for _, variable := range schema.Variables {
		name := strings.TrimSpace(variable.Name)
		if name == "" {
			report.warnings = append(report.warnings, "contract contains an unnamed variable")
			continue
		}
		if seen[name] {
			report.warnings = append(report.warnings, fmt.Sprintf("%s is declared more than once", name))
			continue
		}
		seen[name] = true
		if variable.Secret {
			report.secretCount++
		}
		if variable.BrowserExposed {
			report.browserCount++
		}
		if variable.FrameworkOwned {
			report.frameworkCount++
		}
		if variable.Secret && variable.BrowserExposed {
			report.warnings = append(report.warnings, fmt.Sprintf("%s is secret and browser-exposed", name))
		}
		if variable.Required && envValue(name, variable.LocalDefault, dotenv) == "" {
			report.missingRequired = append(report.missingRequired, name)
		}
	}
	return report, nil
}

func readEnvSchema(path string) (envSchema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return envSchema{}, fmt.Errorf("missing %s", path)
		}
		return envSchema{}, err
	}
	schema, err := parseEnvContractTOML(string(data))
	if err != nil {
		return envSchema{}, fmt.Errorf("invalid %s: %w", path, err)
	}
	if schema.Version == 0 {
		return envSchema{}, fmt.Errorf("%s is missing env.contract_version", path)
	}
	if len(schema.Variables) == 0 {
		return envSchema{}, fmt.Errorf("%s declares no env variables", path)
	}
	return schema, nil
}

func parseEnvContractTOML(content string) (envSchema, error) {
	var schema envSchema
	var variables []envVariable
	var current *envVariable
	seen := map[string]bool{}
	section := ""

	flushVariable := func() error {
		if current == nil {
			return nil
		}
		current.Name = strings.TrimSpace(current.Name)
		if current.Name == "" {
			return errors.New("env variable table has an empty name")
		}
		if seen[current.Name] {
			return fmt.Errorf("%s is declared more than once", current.Name)
		}
		seen[current.Name] = true
		variables = append(variables, *current)
		current = nil
		return nil
	}

	scanner := bufio.NewScanner(strings.NewReader(content))
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(stripTomlComment(scanner.Text()))
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			if err := flushVariable(); err != nil {
				return envSchema{}, err
			}
			section = strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(line, "["), "]"))
			if strings.HasPrefix(section, "env.variables.") {
				name := strings.TrimSpace(strings.TrimPrefix(section, "env.variables."))
				name = strings.Trim(name, `"`)
				current = &envVariable{Name: name}
			}
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		switch {
		case section == "env" && key == "contract_version":
			version, err := strconv.Atoi(value)
			if err != nil {
				return envSchema{}, fmt.Errorf("line %d has invalid env.contract_version", lineNumber)
			}
			schema.Version = version
		case strings.HasPrefix(section, "env.variables."):
			if current == nil {
				return envSchema{}, fmt.Errorf("line %d is outside an env variable table", lineNumber)
			}
			if err := assignEnvVariableField(current, key, value, lineNumber); err != nil {
				return envSchema{}, err
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return envSchema{}, err
	}
	if err := flushVariable(); err != nil {
		return envSchema{}, err
	}
	schema.Variables = variables
	return schema, nil
}

func assignEnvVariableField(variable *envVariable, key string, value string, lineNumber int) error {
	switch key {
	case "service":
		variable.Service = parseTomlString(value)
	case "required":
		variable.Required = parseTomlBool(value)
	case "secret":
		variable.Secret = parseTomlBool(value)
	case "browser_exposed":
		variable.BrowserExposed = parseTomlBool(value)
	case "framework_owned":
		variable.FrameworkOwned = parseTomlBool(value)
	case "local_default":
		variable.LocalDefault = parseTomlString(value)
	case "description":
		variable.Description = parseTomlString(value)
	default:
		return fmt.Errorf("line %d has unknown env variable field %q", lineNumber, key)
	}
	return nil
}

func parseTomlBool(value string) bool {
	return strings.EqualFold(strings.TrimSpace(value), "true")
}

func parseTomlString(value string) string {
	return unquoteEnvValue(strings.TrimSpace(value))
}

func stripTomlComment(line string) string {
	inString := false
	escaped := false
	for i, r := range line {
		if escaped {
			escaped = false
			continue
		}
		if r == '\\' && inString {
			escaped = true
			continue
		}
		if r == '"' {
			inString = !inString
			continue
		}
		if r == '#' && !inString {
			return line[:i]
		}
	}
	return line
}

func readDotenv(path string) (map[string]string, bool, error) {
	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return map[string]string{}, false, nil
		}
		return nil, false, err
	}
	defer file.Close()

	values := map[string]string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		values[key] = unquoteEnvValue(strings.TrimSpace(value))
	}
	if err := scanner.Err(); err != nil {
		return nil, false, err
	}
	return values, true, nil
}

func envValue(name string, localDefault string, dotenv map[string]string) string {
	if value, ok := dotenv[name]; ok && strings.TrimSpace(value) != "" {
		return strings.TrimSpace(value)
	}
	if value := strings.TrimSpace(os.Getenv(name)); value != "" {
		return value
	}
	return strings.TrimSpace(localDefault)
}

func unquoteEnvValue(value string) string {
	if len(value) < 2 {
		return value
	}
	if (value[0] == '"' && value[len(value)-1] == '"') ||
		(value[0] == '\'' && value[len(value)-1] == '\'') {
		return value[1 : len(value)-1]
	}
	return value
}

func ensureDeployTarget(target string) error {
	if !regexp.MustCompile(`^[a-z][a-z0-9-]*$`).MatchString(target) {
		return errors.New("deploy target must use lowercase letters, numbers, and dashes")
	}
	return nil
}

func loadDeployTarget(name string) (deployTarget, bool, error) {
	targets, err := loadDeployTargets()
	if err != nil {
		return deployTarget{}, false, err
	}
	target, found := targets[name]
	return target, found, nil
}

func loadDeployTargets() (map[string]deployTarget, error) {
	content, err := os.ReadFile(projectConfigPath)
	if err != nil {
		return nil, err
	}

	targets := map[string]deployTarget{}
	section := ""
	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(stripTomlComment(scanner.Text()))
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section = strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(line, "["), "]"))
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		if !strings.HasPrefix(section, "deploy.targets.") {
			continue
		}
		targetName := strings.TrimSpace(strings.TrimPrefix(section, "deploy.targets."))
		if targetName == "" || strings.Contains(targetName, ".") {
			continue
		}
		target := targets[targetName]
		target.Name = targetName
		if err := assignDeployTargetField(&target, key, value, lineNumber); err != nil {
			return nil, err
		}
		targets[targetName] = target
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return targets, nil
}

func assignDeployTargetField(target *deployTarget, key string, value string, lineNumber int) error {
	switch key {
	case "script":
		target.Script = parseTomlString(value)
	case "description":
		target.Description = parseTomlString(value)
	default:
		_ = lineNumber
	}
	return nil
}

func validateDeployTarget(target deployTarget) error {
	if err := ensureDeployTarget(target.Name); err != nil {
		return err
	}
	if strings.TrimSpace(target.Script) == "" {
		return fmt.Errorf("deploy target %s is missing script", target.Name)
	}
	if strings.ContainsAny(target.Script, "\r\n") {
		return fmt.Errorf("deploy target %s script is invalid", target.Name)
	}
	root, err := filepath.Abs(".")
	if err != nil {
		return err
	}
	_, err = deployScriptPath(root, target)
	return err
}

func deployScriptPath(root string, target deployTarget) (string, error) {
	clean := filepath.Clean(strings.TrimSpace(target.Script))
	if clean == "." || clean == "" {
		return "", fmt.Errorf("deploy target %s is missing script", target.Name)
	}
	if filepath.IsAbs(clean) {
		return "", fmt.Errorf("deploy target %s script must stay inside the project", target.Name)
	}
	absolute := filepath.Join(root, clean)
	absolute, err := filepath.Abs(absolute)
	if err != nil {
		return "", err
	}
	relative, err := filepath.Rel(root, absolute)
	if err != nil {
		return "", err
	}
	if relative == ".." || strings.HasPrefix(relative, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("deploy target %s script must stay inside the project", target.Name)
	}
	info, err := os.Stat(absolute)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("deploy target %s script does not exist: %s", target.Name, clean)
		}
		return "", err
	}
	if info.IsDir() {
		return "", fmt.Errorf("deploy target %s script is a directory: %s", target.Name, clean)
	}
	return absolute, nil
}

func deployScriptEnv(root string, target deployTarget) []string {
	metadata, err := readProjectMetadata(projectConfigPath)
	if err != nil {
		return []string{
			fmt.Sprintf("CARBIDE_DEPLOY_TARGET=%s", target.Name),
			fmt.Sprintf("CARBIDE_PROJECT_ROOT=%s", root),
		}
	}
	projectName := strings.TrimSpace(metadata.name)
	projectSlug := strings.TrimSpace(metadata.slug)
	values := []string{
		fmt.Sprintf("CARBIDE_DEPLOY_TARGET=%s", target.Name),
		fmt.Sprintf("CARBIDE_PROJECT_ROOT=%s", root),
	}
	if projectName != "" && !isTemplatePlaceholder(projectName) {
		values = append(values, fmt.Sprintf("CARBIDE_PROJECT_NAME=%s", projectName))
	}
	if projectSlug != "" && !isTemplatePlaceholder(projectSlug) {
		values = append(values, fmt.Sprintf("CARBIDE_PROJECT_SLUG=%s", projectSlug))
	}
	return values
}
