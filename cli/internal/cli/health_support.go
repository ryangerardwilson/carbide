package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

func resolveFrameworkRoot(home string) (string, error) {
	starts := []string{}
	if cwd, err := os.Getwd(); err == nil {
		starts = append(starts, cwd)
	}
	if home != "" {
		starts = append(starts, home)
	}

	seen := map[string]bool{}
	for _, start := range starts {
		root, err := filepath.Abs(start)
		if err != nil {
			continue
		}
		for {
			if !seen[root] {
				seen[root] = true
				if isFrameworkSourceRoot(root) {
					return root, nil
				}
			}
			parent := filepath.Dir(root)
			if parent == root {
				break
			}
			root = parent
		}
	}
	return "", errors.New("run from a Carbide source checkout")
}

func isFrameworkSourceRoot(root string) bool {
	return isFile(filepath.Join(root, "cli", "go.mod")) &&
		isDir(filepath.Join(root, "tests")) &&
		isFile(filepath.Join(root, "scaffold", "carbide.toml"))
}

func missingFiles(paths []string) []string {
	var missing []string
	for _, path := range paths {
		if !isFile(path) {
			missing = append(missing, path)
		}
	}
	return missing
}

func missingDirs(paths []string) []string {
	var missing []string
	for _, path := range paths {
		if !isDir(path) {
			missing = append(missing, path)
		}
	}
	return missing
}

func existingFiles(paths []string) []string {
	var found []string
	for _, path := range paths {
		if isFile(path) {
			found = append(found, path)
		}
	}
	return found
}

func existingDirs(paths []string) []string {
	var found []string
	for _, path := range paths {
		if isDir(path) {
			found = append(found, path)
		}
	}
	return found
}

func missingNeedles(content string, needles []string) []string {
	var missing []string
	for _, needle := range needles {
		if !strings.Contains(content, needle) {
			missing = append(missing, needle)
		}
	}
	return missing
}

func readFileString(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}

func fileContains(path string, needle string) bool {
	return strings.Contains(readFileString(path), needle)
}

func fileLineCount(path string) int {
	content := strings.TrimSuffix(readFileString(path), "\n")
	if content == "" {
		return 0
	}
	return strings.Count(content, "\n") + 1
}

func lawLineLimitViolations(root string, limit int) ([]string, error) {
	var violations []string

	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		name := entry.Name()
		if entry.IsDir() {
			if shouldSkipLawLineLimitDir(name) {
				return filepath.SkipDir
			}
			return nil
		}
		if !entry.Type().IsRegular() {
			return nil
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if !isLawLineLimitFile(rel) {
			return nil
		}

		lines := fileLineCount(path)
		if lines > limit {
			violations = append(violations, fmt.Sprintf("%s (%d)", rel, lines))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Strings(violations)
	return violations, nil
}

func shouldSkipLawLineLimitDir(name string) bool {
	switch name {
	case ".git", ".carbide", ".audit", ".cli", ".bin", "node_modules", "vendor", "public", "dist", "build", "coverage":
		return true
	default:
		return false
	}
}

func isLawLineLimitFile(path string) bool {
	switch {
	case strings.HasSuffix(path, ".go"),
		strings.HasSuffix(path, ".ts"),
		strings.HasSuffix(path, ".tsx"),
		strings.HasSuffix(path, ".js"),
		strings.HasSuffix(path, ".jsx"),
		strings.HasSuffix(path, ".sh"),
		strings.HasSuffix(path, ".css"),
		strings.HasSuffix(path, ".sql"),
		strings.HasSuffix(path, ".toml"),
		strings.HasSuffix(path, ".yml"),
		strings.HasSuffix(path, ".yaml"),
		strings.HasSuffix(path, "Dockerfile"):
		return true
	default:
		return false
	}
}

func containedNeedles(content string, needles []string) []string {
	var found []string
	for _, needle := range needles {
		if strings.Contains(content, needle) {
			found = append(found, needle)
		}
	}
	return found
}

func floatingDockerReferences(paths []string) []string {
	var findings []string
	for _, path := range paths {
		content := readFileString(path)
		for _, ref := range dockerImageRefs(path, content) {
			if isFloatingImageRef(ref) {
				findings = append(findings, path+" "+ref)
			}
		}
	}
	return findings
}

func dockerImageRefs(path string, content string) []string {
	var refs []string
	if strings.HasSuffix(path, "Dockerfile") {
		scanner := bufio.NewScanner(strings.NewReader(content))
		for scanner.Scan() {
			fields := strings.Fields(scanner.Text())
			if len(fields) >= 2 && strings.EqualFold(fields[0], "FROM") {
				refs = append(refs, strings.Trim(fields[1], `"'`))
			}
		}
		return refs
	}

	imageLine := regexp.MustCompile(`(?m)^\s*image:\s*["']?([^"'\s]+)`)
	for _, match := range imageLine.FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			refs = append(refs, strings.Trim(match[1], `"'`))
		}
	}
	return refs
}

func isFloatingImageRef(ref string) bool {
	if ref == "" || strings.Contains(ref, "${") {
		return false
	}
	return !strings.Contains(ref, "@sha256:") || strings.Contains(ref, ":latest")
}

func packageVersionRangeFindings(path string) []string {
	return packageVersionRangeFindingsFor(path, []string{"react", "react-dom", "tailwindcss", "@tailwindcss/cli"})
}

func packageVersionRangeFindingsFor(path string, names []string) []string {
	content := readFileString(path)
	var findings []string
	for _, name := range names {
		version := packageVersion(content, name)
		if version == "" {
			findings = append(findings, name+" missing")
			continue
		}
		if isSemverRange(version) {
			findings = append(findings, name+" "+version)
		}
	}
	return findings
}

func packageVersion(content string, name string) string {
	pattern := regexp.MustCompile(`"` + regexp.QuoteMeta(name) + `"\s*:\s*"([^"]+)"`)
	match := pattern.FindStringSubmatch(content)
	if len(match) < 2 {
		return ""
	}
	return match[1]
}

func isSemverRange(version string) bool {
	trimmed := strings.TrimSpace(version)
	if trimmed == "" {
		return true
	}
	rangeMarkers := []string{"^", "~", ">", "<", "*", "x", "X", "latest", "||", " - "}
	for _, marker := range rangeMarkers {
		if strings.Contains(trimmed, marker) {
			return true
		}
	}
	return false
}

func unsupportedGoDirectiveFindings(paths []string) []string {
	var findings []string
	for _, path := range paths {
		version := goDirectiveVersion(path)
		if version != baselineGoModuleVersion {
			if version == "" {
				version = "missing"
			}
			findings = append(findings, path+" "+version)
		}
	}
	return findings
}

func goDirectiveVersion(path string) string {
	scanner := bufio.NewScanner(strings.NewReader(readFileString(path)))
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) == 2 && fields[0] == "go" {
			return fields[1]
		}
	}
	return ""
}

func containsString(values []string, value string) bool {
	for _, item := range values {
		if item == value {
			return true
		}
	}
	return false
}

func composeServiceNamesFromFile(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	var services []string
	seen := map[string]bool{}
	inServices := false
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		raw := scanner.Text()
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		if !strings.HasPrefix(raw, " ") && strings.HasSuffix(trimmed, ":") {
			inServices = trimmed == "services:"
			continue
		}
		if !inServices {
			continue
		}
		if strings.HasPrefix(raw, "  ") && !strings.HasPrefix(raw, "    ") && strings.HasSuffix(trimmed, ":") {
			service := strings.TrimSuffix(trimmed, ":")
			if service != "" && !seen[service] {
				seen[service] = true
				services = append(services, service)
			}
		}
	}
	return services
}

func rootDirsOutsideContract(allowed map[string]bool) []string {
	entries, err := os.ReadDir(".")
	if err != nil {
		return []string{err.Error()}
	}
	var extras []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		if !allowed[name] {
			extras = append(extras, name)
		}
	}
	return extras
}

func treeContains(root string, needle string) bool {
	return len(treeContainsAny(root, []string{needle})) > 0
}

func treeContainsAny(root string, needles []string) []string {
	found := map[string]bool{}
	_ = filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return nil
		}
		name := entry.Name()
		if entry.IsDir() {
			switch name {
			case ".git", ".carbide", ".cli", ".bin", "node_modules", "vendor":
				return filepath.SkipDir
			default:
				return nil
			}
		}
		if !entry.Type().IsRegular() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		text := string(data)
		for _, needle := range needles {
			if strings.Contains(text, needle) {
				found[needle] = true
			}
		}
		return nil
	})

	var hits []string
	for _, needle := range needles {
		if found[needle] {
			hits = append(hits, needle)
		}
	}
	return hits
}

func anyPathWithExtension(root string, extensions ...string) bool {
	found := false
	_ = filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil || found {
			return nil
		}
		if entry.IsDir() {
			return nil
		}
		for _, extension := range extensions {
			if strings.HasSuffix(path, extension) {
				found = true
				return nil
			}
		}
		return nil
	})
	return found
}

func composeHasRunningServices(compose composeCommand, env []string) bool {
	snapshots, err := composeServiceSnapshots(compose, env)
	if err != nil {
		return false
	}
	for _, snapshot := range snapshots {
		if strings.EqualFold(strings.TrimSpace(snapshot.State), "running") {
			return true
		}
	}
	return false
}

func publishedWebPort(compose composeCommand, env []string) int {
	snapshots, err := composeServiceSnapshots(compose, env)
	if err != nil {
		return 0
	}
	web, ok := snapshots["web"]
	if !ok {
		return 0
	}
	for _, publisher := range web.Publishers {
		if publisher.PublishedPort > 0 && publisher.TargetPort == 8080 {
			return publisher.PublishedPort
		}
	}
	for _, publisher := range web.Publishers {
		if publisher.PublishedPort > 0 {
			return publisher.PublishedPort
		}
	}
	return 0
}

func runtimePortFromEnv() int {
	if value := os.Getenv("CARBIDE_HTTP_PORT"); value != "" {
		if port, err := validatePort(value); err == nil {
			return port
		}
	}
	return 8080
}

func waitForHTTP(client *http.Client, endpoint string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	var lastErr error
	for time.Now().Before(deadline) {
		resp, err := client.Get(endpoint)
		if err == nil {
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				return nil
			}
			lastErr = fmt.Errorf("%s returned %s", endpoint, resp.Status)
		} else {
			lastErr = err
		}
		time.Sleep(time.Second)
	}
	if lastErr == nil {
		lastErr = fmt.Errorf("%s did not become ready", endpoint)
	}
	return lastErr
}

func httpGetContains(client *http.Client, endpoint string, needle string) error {
	resp, err := client.Get(endpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("%s returned %s", endpoint, resp.Status)
	}
	if !strings.Contains(string(body), needle) {
		return fmt.Errorf("%s did not contain %q", endpoint, needle)
	}
	return nil
}

func httpPostFormContains(client *http.Client, endpoint string, values url.Values, needle string) error {
	if values == nil {
		values = url.Values{}
	}
	resp, err := client.PostForm(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("%s returned %s", endpoint, resp.Status)
	}
	if !strings.Contains(string(body), needle) {
		return fmt.Errorf("%s did not contain %q", endpoint, needle)
	}
	return nil
}

func runFrameworkGoTests(home string) error {
	if _, err := exec.LookPath("go"); err == nil {
		_, err := commandOutput(filepath.Join(home, "cli"), "go", "test", "./...")
		return err
	}
	if _, err := exec.LookPath("docker"); err != nil {
		return errors.New("Go is not installed and Docker is unavailable for Go test fallback")
	}
	args := []string{
		"run",
		"--rm",
		"--user", fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid()),
		"-e", "HOME=/tmp",
		"-e", "GOCACHE=/tmp/gocache",
		"-e", "GOMODCACHE=/tmp/gomodcache",
		"-v", home + ":" + home,
		"-w", filepath.Join(home, "cli"),
		"golang:1.25-bookworm",
		"bash",
		"-lc",
		"export PATH=/usr/local/go/bin:$PATH; go test ./...",
	}
	_, err := commandOutput("", "docker", args...)
	return err
}

func frameworkHealthCommandEnv(home string) ([]string, func(), error) {
	env := setEnv(os.Environ(), "CARBIDE_HOME", home)
	if _, err := exec.LookPath("go"); err == nil {
		return env, func() {}, nil
	}
	dockerPath, err := exec.LookPath("docker")
	if err != nil {
		return nil, func() {}, errors.New("Go is not installed and Docker is unavailable for framework checks")
	}

	dir, err := os.MkdirTemp("", "carbide-health-go-")
	if err != nil {
		return nil, func() {}, err
	}
	cleanup := func() {
		_ = os.RemoveAll(dir)
	}
	wrapper := fmt.Sprintf(`#!/usr/bin/env bash
set -euo pipefail
root="${CARBIDE_HOME:-$(pwd -P)}"
workdir="$(pwd -P)"
exec %q run --rm \
  --user "$(id -u):$(id -g)" \
  -e HOME=/tmp \
  -e GOCACHE=/tmp/gocache \
  -e GOMODCACHE=/tmp/gomodcache \
  -v "$root:$root" \
  -w "$workdir" \
  golang:1.25-bookworm \
  /usr/local/go/bin/go "$@"
`, dockerPath)
	if err := os.WriteFile(filepath.Join(dir, "go"), []byte(wrapper), 0755); err != nil {
		cleanup()
		return nil, func() {}, err
	}
	env = setEnv(env, "PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	return env, cleanup, nil
}

func firstLine(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if line, _, ok := strings.Cut(value, "\n"); ok {
		return strings.TrimSpace(line)
	}
	return value
}
