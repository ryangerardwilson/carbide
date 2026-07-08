package cli

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func findCompose() (composeCommand, error) {
	if _, err := commandOutput("", "docker", "compose", "version"); err == nil {
		help, _ := commandOutput("", "docker", "compose", "up", "--help")
		logHelp, _ := commandOutput("", "docker", "compose", "logs", "--help")
		return composeCommand{name: "docker", base: []string{"compose"}, help: help, logHelp: logHelp}, nil
	}
	if _, err := commandOutput("", "docker-compose", "version"); err == nil {
		help, _ := commandOutput("", "docker-compose", "up", "--help")
		logHelp, _ := commandOutput("", "docker-compose", "logs", "--help")
		return composeCommand{name: "docker-compose", help: help, logHelp: logHelp}, nil
	}
	return composeCommand{}, errors.New("Docker Compose is required for carbide run dev")
}

func composeEnv(env []string) []string {
	if isFile(composeFilePath) {
		env = setEnv(env, "COMPOSE_FILE", composeFilePath)
	} else if isFile(legacyComposeFilePath) {
		env = setEnv(env, "COMPOSE_FILE", legacyComposeFilePath)
	}
	if isFile(projectConfigPath) {
		env = setEnv(env, "COMPOSE_PROJECT_NAME", composeProjectName())
	}
	return env
}

func (c composeCommand) args(extra ...string) []string {
	args := make([]string, 0, len(c.base)+len(extra))
	args = append(args, c.base...)
	args = append(args, extra...)
	return args
}

func (c composeCommand) supports(option string) bool {
	return strings.Contains(c.help, option)
}

func (c composeCommand) logsSupports(option string) bool {
	return strings.Contains(c.logHelp, option)
}

func composeServices(compose composeCommand, env []string) []string {
	output, err := runComposeCaptured(compose, env, "config", "--services")
	if err != nil {
		return defaultComposeServices()
	}
	seen := map[string]bool{}
	var services []string
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		service := strings.TrimSpace(scanner.Text())
		if service == "" || seen[service] {
			continue
		}
		seen[service] = true
		services = append(services, service)
	}
	if len(services) == 0 {
		return defaultComposeServices()
	}
	return services
}

func defaultComposeServices() []string {
	return []string{"web", "api", "db"}
}

func composeServiceStatuses(compose composeCommand, env []string) map[string]composeServiceStatus {
	snapshots, err := composeServiceSnapshots(compose, env)
	if err != nil {
		return nil
	}
	return composeStatusesFromSnapshots(snapshots)
}

func composeServiceSnapshots(compose composeCommand, env []string) (map[string]composeServiceSnapshot, error) {
	output, err := runComposeCaptured(compose, env, "ps", "--format", "json")
	if err != nil {
		return nil, err
	}
	return parseComposeServiceSnapshots(output)
}

func parseComposeServiceStatuses(output string) (map[string]composeServiceStatus, error) {
	snapshots, err := parseComposeServiceSnapshots(output)
	if err != nil {
		return nil, err
	}
	return composeStatusesFromSnapshots(snapshots), nil
}

func parseComposeServiceSnapshots(output string) (map[string]composeServiceSnapshot, error) {
	parseRecords := func(records []composeServiceSnapshot) map[string]composeServiceSnapshot {
		snapshots := map[string]composeServiceSnapshot{}
		for _, record := range records {
			service := strings.TrimSpace(record.Service)
			if service == "" {
				continue
			}
			record.Service = service
			record.Name = strings.TrimSpace(record.Name)
			record.State = strings.TrimSpace(record.State)
			record.Health = strings.TrimSpace(record.Health)
			record.Status = strings.TrimSpace(record.Status)
			record.Ports = strings.TrimSpace(record.Ports)
			snapshots[service] = record
		}
		return snapshots
	}

	var records []composeServiceSnapshot
	if err := json.Unmarshal([]byte(output), &records); err == nil {
		return parseRecords(records), nil
	}

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var record composeServiceSnapshot
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return parseRecords(records), nil
}

func composeStatusesFromSnapshots(snapshots map[string]composeServiceSnapshot) map[string]composeServiceStatus {
	statuses := map[string]composeServiceStatus{}
	for service, snapshot := range snapshots {
		statuses[service] = composeServiceStatus{
			service: service,
			state:   snapshot.State,
			health:  snapshot.Health,
		}
	}
	return statuses
}

func composeStatusRow(snapshot composeServiceSnapshot) tableRow {
	return tableRow{
		snapshot.Service,
		statusValue(snapshot.Name),
		composePublishedPorts(snapshot),
		composeInternalPorts(snapshot),
		composeServiceStatusText(snapshot),
	}
}

func statusValue(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "-"
	}
	return value
}

func composePublishedPorts(snapshot composeServiceSnapshot) string {
	ports := composePublishedPortValues(snapshot)
	if len(ports) == 0 {
		return "-"
	}
	return strings.Join(ports, ", ")
}

func composePublishedPortValues(snapshot composeServiceSnapshot) []string {
	seen := map[string]bool{}
	var ports []string
	for _, publisher := range snapshot.Publishers {
		if publisher.PublishedPort <= 0 {
			continue
		}
		host := strings.TrimSpace(publisher.URL)
		if host == "" || host == "0.0.0.0" || host == "::" {
			host = "localhost"
		}
		value := fmt.Sprintf("%s:%d", host, publisher.PublishedPort)
		if !seen[value] {
			seen[value] = true
			ports = append(ports, value)
		}
	}
	return ports
}

func composeInternalPorts(snapshot composeServiceSnapshot) string {
	ports := composeInternalPortValues(snapshot)
	if len(ports) == 0 && strings.TrimSpace(snapshot.Ports) != "" {
		return strings.TrimSpace(snapshot.Ports)
	}
	if len(ports) == 0 {
		return "-"
	}
	return strings.Join(ports, ", ")
}

func composeInternalPortValues(snapshot composeServiceSnapshot) []string {
	seen := map[string]bool{}
	var ports []string
	for _, publisher := range snapshot.Publishers {
		if publisher.TargetPort <= 0 {
			continue
		}
		protocol := strings.TrimSpace(publisher.Protocol)
		if protocol == "" {
			protocol = "tcp"
		}
		value := fmt.Sprintf("%d/%s", publisher.TargetPort, protocol)
		if !seen[value] {
			seen[value] = true
			ports = append(ports, value)
		}
	}
	if len(ports) == 0 {
		raw := strings.TrimSpace(snapshot.Ports)
		if raw != "" {
			return []string{raw}
		}
	}
	return ports
}

func composeServiceStatusText(snapshot composeServiceSnapshot) string {
	state := strings.TrimSpace(strings.ToLower(snapshot.State))
	health := strings.TrimSpace(strings.ToLower(snapshot.Health))
	if state == "" {
		return statusValue(snapshot.Status)
	}
	if health != "" {
		return fmt.Sprintf("%s (%s)", state, health)
	}
	return state
}

func composeLogsArgs(compose composeCommand) []string {
	args := []string{"logs", "-f", "--tail", "80"}
	if compose.logsSupports("--no-color") {
		args = append(args, "--no-color")
	}
	return args
}

func parseComposeLogLine(line string) (string, string) {
	line = stripANSI(strings.TrimSpace(line))
	parts := strings.SplitN(line, "|", 2)
	if len(parts) != 2 {
		return "log", line
	}
	service := normalizeServiceName(strings.TrimSpace(parts[0]))
	message := strings.TrimSpace(parts[1])
	return service, message
}

func normalizeServiceName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "log"
	}
	parts := strings.Split(name, "-")
	if len(parts) > 1 {
		last := parts[len(parts)-1]
		if _, err := strconv.Atoi(last); err == nil {
			name = strings.Join(parts[:len(parts)-1], "-")
		}
	}
	if idx := strings.LastIndex(name, "-"); idx >= 0 {
		name = name[idx+1:]
	}
	switch name {
	case "web", "api", "db":
		return name
	default:
		return name
	}
}

func stripANSI(value string) string {
	var out strings.Builder
	inEscape := false
	inCSI := false
	for i := 0; i < len(value); i++ {
		ch := value[i]
		if inEscape {
			if !inCSI && ch == '[' {
				inCSI = true
				continue
			}
			if inCSI {
				if ch >= '@' && ch <= '~' {
					inEscape = false
					inCSI = false
				}
				continue
			}
			if ch >= 0x30 && ch <= 0x7e {
				inEscape = false
			}
			continue
		}
		if ch == 0x1b {
			inEscape = true
			continue
		}
		out.WriteByte(ch)
	}
	return out.String()
}

func composeUpDetached(compose composeCommand, env []string) error {
	args := []string{"up", "-d", "--build", "--remove-orphans"}
	if compose.supports("--quiet-build") {
		args = append(args, "--quiet-build")
	}
	if compose.supports("--quiet-pull") {
		args = append(args, "--quiet-pull")
	}
	if compose.supports("--wait") {
		args = append(args, "--wait", "--wait-timeout", "120")
	}
	output, err := runComposeCaptured(compose, env, args...)
	if err != nil {
		return fmt.Errorf("Docker Compose start failed: %w\n%s", err, strings.TrimSpace(output))
	}
	return nil
}

func composeDown(compose composeCommand, env []string) error {
	output, err := runComposeCaptured(compose, env, "down", "--remove-orphans")
	if err != nil {
		return fmt.Errorf("Docker Compose stop failed: %w\n%s", err, strings.TrimSpace(output))
	}
	return nil
}
