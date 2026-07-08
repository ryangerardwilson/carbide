package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func newRenderer(out io.Writer) renderer {
	interactive := isTerminalOutput(out)
	termWidth := 0
	if interactive {
		termWidth = terminalColumns(out)
		if termWidth == 0 {
			termWidth = terminalColumnsFromEnv()
		}
	}
	return renderer{
		out:         out,
		interactive: interactive,
		styled:      interactive && os.Getenv("NO_COLOR") == "",
		termWidth:   termWidth,
	}
}

func renderError(out io.Writer, err error) {
	newRenderer(out).Message(
		"Carbide",
		"command failed",
		outputRow{"error", err.Error()},
		outputRow{"help", "carbide help"},
	)
}

func (r renderer) Message(title string, subtitle string, rows ...outputRow) {
	r.Title(title, subtitle)
	r.Rows(rows...)
}

func (r renderer) Title(title string, subtitle string) {
	if r.styled {
		fmt.Fprintf(r.out, "%s\n", r.paint("1;38;5;81", title))
		if subtitle != "" {
			fmt.Fprintf(r.out, "%s\n", r.paint("2;38;5;245", subtitle))
		}
	} else {
		fmt.Fprintln(r.out, title)
		if subtitle != "" {
			fmt.Fprintln(r.out, subtitle)
		}
	}
	fmt.Fprintln(r.out)
}

func (r renderer) Section(title string, subtitle string) {
	fmt.Fprintln(r.out)
	if r.styled {
		fmt.Fprintf(r.out, "%s\n", r.paint("1;38;5;245", title))
		if subtitle != "" {
			fmt.Fprintf(r.out, "%s\n", r.paint("2;38;5;245", subtitle))
		}
	} else {
		fmt.Fprintln(r.out, title)
		if subtitle != "" {
			fmt.Fprintln(r.out, subtitle)
		}
	}
	fmt.Fprintln(r.out)
}

func (r renderer) Rows(rows ...outputRow) {
	width := rowKeyWidth(rows)
	for _, row := range rows {
		r.writeRow(row, width)
	}
}

func (r renderer) Table(headers []string, rows []tableRow) {
	widths := make([]int, len(headers))
	for index, header := range headers {
		widths[index] = len(header)
	}
	for _, row := range rows {
		for index := range headers {
			value := ""
			if index < len(row) {
				value = row[index]
			}
			if len(value) > widths[index] {
				widths[index] = len(value)
			}
		}
	}

	writeCells := func(cells []string, header bool) {
		for index := range headers {
			value := ""
			if index < len(cells) {
				value = cells[index]
			}
			if index > 0 {
				fmt.Fprint(r.out, "  ")
			}
			padded := value
			if index < len(headers)-1 {
				padded += strings.Repeat(" ", widths[index]-len(value))
			}
			if header {
				padded = r.paint("2;38;5;245", padded)
			}
			fmt.Fprint(r.out, padded)
		}
		fmt.Fprintln(r.out)
	}

	writeCells(headers, true)
	for _, row := range rows {
		writeCells([]string(row), false)
	}
}

func (r renderer) CommandList(sections []helpCommandSection) {
	fmt.Fprintln(r.out, r.formatHelpHeading("Usage:"))
	fmt.Fprintln(r.out, "  carbide <command> [arguments]")
	fmt.Fprintln(r.out)
	fmt.Fprintln(r.out, r.formatHelpHeading("Available commands:"))

	width := helpCommandWidth(sections)
	for _, section := range sections {
		if section.name != "" {
			fmt.Fprintln(r.out, r.formatHelpGroup(section.name))
		}
		for _, row := range section.rows {
			r.writeHelpCommand(row, width)
		}
	}
}

func (r renderer) writeHelpCommand(row outputRow, width int) {
	valueWidth := helpOutputMaxWidth - 4 - width
	if valueWidth < 16 {
		valueWidth = 16
	}
	lines := wrapHelpText(row.value, valueWidth)
	if len(lines) == 0 {
		lines = []string{""}
	}
	if r.styled {
		key := r.paint("38;5;245", row.key)
		fmt.Fprintf(r.out, "  %s%s  %s\n", key, strings.Repeat(" ", width-len(row.key)), lines[0])
	} else {
		fmt.Fprintf(r.out, "  %-*s  %s\n", width, row.key, lines[0])
	}
	padding := strings.Repeat(" ", width+4)
	for _, line := range lines[1:] {
		fmt.Fprintf(r.out, "%s%s\n", padding, line)
	}
}

func (r renderer) formatHelpHeading(value string) string {
	return r.paint("1;38;5;245", value)
}

func (r renderer) formatHelpGroup(value string) string {
	return r.paint("1;38;5;245", value)
}

func (r renderer) Row(row outputRow) {
	r.writeRow(row, len(row.key))
}

func (r renderer) Blank() {
	fmt.Fprintln(r.out)
}

func (r renderer) Logo(logo string) {
	for index, line := range logoLines(logo) {
		fmt.Fprintln(r.out, r.formatLogoLine(index, line))
	}
	fmt.Fprintln(r.out)
}

func (r renderer) AnimateLogo(logo string) {
	lines := logoLines(logo)
	if len(lines) == 0 {
		return
	}

	width := maxLineWidth(lines)
	chompFrames := width + (len(lines)-1)*2 + 1
	for frame := 0; frame <= chompFrames; frame++ {
		if frame > 0 {
			fmt.Fprintf(r.out, "\033[%dA", len(lines))
		}
		for index, line := range lines {
			position := frame - index*2
			fmt.Fprintf(r.out, "\r\033[K%s\n", r.formatLogoPacmanLine(line, position, frame+index))
		}
		if frame < chompFrames {
			time.Sleep(9 * time.Millisecond)
		}
	}
	fmt.Fprintln(r.out)
}

func (r renderer) writeRow(row outputRow, width int) {
	lines := strings.Split(row.value, "\n")
	if len(lines) > 1 {
		r.writeSingleLine(outputRow{row.key, lines[0]}, width)
		for _, line := range lines[1:] {
			r.writeSingleLine(outputRow{"", line}, width)
		}
		return
	}
	r.writeSingleLine(row, width)
}

func (r renderer) writeSingleLine(row outputRow, width int) {
	if row.key == "" {
		fmt.Fprintf(r.out, "%*s  %s\n", width, "", r.formatValue(row))
		return
	}
	key := r.formatKey(row.key)
	if r.styled {
		fmt.Fprintf(r.out, "%s%s  %s\n", key, strings.Repeat(" ", width-len(row.key)), r.formatValue(row))
		return
	}
	fmt.Fprintf(r.out, "%-*s  %s\n", width, row.key, row.value)
}

func (r renderer) Log(service string, message string) {
	r.LogAt(time.Now(), service, message)
}

func (r renderer) LogEntry(entry structuredLogEntry) {
	r.LogAt(entryTimestamp(entry), entry.Service, entry.Message)
}

func (r renderer) LogAt(timestamp time.Time, service string, message string) {
	label := service
	if label == "" {
		label = "log"
	}
	if timestamp.IsZero() {
		timestamp = time.Now()
	}
	stamp := timestamp.Local().Format("15:04:05")
	width := 9
	if len(label) > width {
		width = len(label)
	}
	if r.styled {
		fmt.Fprintf(
			r.out,
			"%s  %s%s  %s\n",
			r.paint("2;38;5;245", stamp),
			r.formatService(label),
			strings.Repeat(" ", width-len(label)),
			message,
		)
		return
	}
	fmt.Fprintf(r.out, "%s  %-*s  %s\n", stamp, width, label, message)
}

func (r renderer) RunServiceProgress(
	services []string,
	poll func() map[string]composeServiceStatus,
	work func() error,
) error {
	return r.runServiceProgress("start", services, poll, work)
}

func (r renderer) RunServiceStopProgress(
	services []string,
	poll func() map[string]composeServiceStatus,
	work func() error,
) error {
	return r.runServiceProgress("stop", services, poll, work)
}

func (r renderer) runServiceProgress(
	mode string,
	services []string,
	poll func() map[string]composeServiceStatus,
	work func() error,
) error {
	if !r.interactive {
		return work()
	}

	done := make(chan error, 1)
	go func() {
		done <- work()
	}()

	if len(services) == 0 {
		services = []string{"containers"}
	}
	statuses := map[string]composeServiceStatus{}
	step := 0
	r.writeServiceProgress(mode, services, statuses, step)
	ticker := time.NewTicker(140 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case err := <-done:
			if err != nil {
				for _, service := range services {
					statuses[service] = composeServiceStatus{service: service, state: "failed"}
				}
				r.rewriteServiceProgress(mode, services, statuses, step)
				return err
			}
			for _, service := range services {
				statuses[service] = composeServiceStatus{service: service, state: progressDoneState(mode), health: "healthy"}
			}
			r.rewriteServiceProgress(mode, services, statuses, step)
			return nil
		case <-ticker.C:
			step++
			if current := poll(); len(current) > 0 {
				statuses = current
			}
			r.rewriteServiceProgress(mode, services, statuses, step)
		}
	}
}

func (r renderer) rewriteServiceProgress(mode string, services []string, statuses map[string]composeServiceStatus, step int) {
	fmt.Fprintf(r.out, "\033[%dA", len(services))
	r.writeServiceProgress(mode, services, statuses, step)
}

func (r renderer) writeServiceProgress(mode string, services []string, statuses map[string]composeServiceStatus, step int) {
	width := rowTextWidth(services)
	frameWidth := r.serviceProgressFrameWidth(width)
	for index, service := range services {
		status := statuses[service]
		state := progressState(mode, status)
		frame := serviceProgressFrame(frameWidth, step+index, state)
		stateLabel := padRight(state, progressStateColumnWidth)
		fmt.Fprintf(
			r.out,
			"\r\033[K%s%s  %s %s\n",
			r.formatService(service),
			strings.Repeat(" ", width-len(service)),
			r.paint(serviceProgressColor(state), frame),
			r.paint("2;38;5;245", stateLabel),
		)
	}
}

func (r renderer) serviceProgressFrameWidth(serviceWidth int) int {
	termWidth := r.currentTerminalWidth()
	frameWidth := termWidth - serviceWidth - progressStateColumnWidth - 5
	if frameWidth < minimumProgressFrameWidth {
		return minimumProgressFrameWidth
	}
	return frameWidth
}

func (r renderer) currentTerminalWidth() int {
	if r.interactive {
		if width := terminalColumns(r.out); width > 0 {
			return width
		}
	}
	if r.termWidth > 0 {
		return r.termWidth
	}
	return defaultTerminalWidth
}

func progressDoneState(mode string) string {
	if mode == "stop" {
		return "stopped"
	}
	return "running"
}

func progressState(mode string, status composeServiceStatus) string {
	if mode == "stop" {
		return serviceStopProgressState(status)
	}
	return serviceProgressState(status)
}

func serviceProgressFrame(width int, step int, state string) string {
	if width < 1 {
		width = 1
	}
	switch state {
	case "ready":
		return "[" + strings.Repeat("#", width) + "]"
	case "stopped":
		return "[" + strings.Repeat(" ", width) + "]"
	case "failed":
		return "[" + strings.Repeat("!", width) + "]"
	case "stopping":
		return reversePacmanFrame(width, step)
	default:
		return pacmanFrame(width, step)
	}
}

func pacmanFrame(width int, step int) string {
	position := step % width
	if position < 0 {
		position = 0
	}
	var b strings.Builder
	b.WriteByte('[')
	for i := 0; i < width; i++ {
		switch {
		case i == position:
			b.WriteByte(pacmanMouth(step, "right"))
		case i < position:
			b.WriteByte('-')
		case isCandyPosition(i - position):
			b.WriteByte('o')
		default:
			b.WriteByte(' ')
		}
	}
	b.WriteByte(']')
	return b.String()
}

func reversePacmanFrame(width int, step int) string {
	position := width - 1 - (step % width)
	if position < 0 {
		position = 0
	}
	var b strings.Builder
	b.WriteByte('[')
	for i := 0; i < width; i++ {
		switch {
		case i == position:
			b.WriteByte(pacmanMouth(step, "left"))
		case i > position:
			b.WriteByte('-')
		case isCandyPosition(position - i):
			b.WriteByte('o')
		default:
			b.WriteByte(' ')
		}
	}
	b.WriteByte(']')
	return b.String()
}

func isCandyPosition(distance int) bool {
	return distance > 1 && (distance-2)%3 == 0
}

func pacmanMouth(step int, direction string) byte {
	open := step%2 == 0
	if direction == "left" {
		if open {
			return 'D'
		}
		return 'd'
	}
	if open {
		return 'C'
	}
	return 'c'
}

func serviceProgressState(status composeServiceStatus) string {
	state := strings.ToLower(status.state)
	health := strings.ToLower(status.health)
	if state == "failed" || state == "exited" || state == "dead" {
		return "failed"
	}
	if state == "running" && (health == "" || health == "healthy") {
		return "ready"
	}
	if health == "healthy" {
		return "ready"
	}
	return "starting"
}

func serviceStopProgressState(status composeServiceStatus) string {
	state := strings.ToLower(status.state)
	if state == "failed" || state == "dead" {
		return "failed"
	}
	if state == "stopped" {
		return "stopped"
	}
	return "stopping"
}

func serviceProgressColor(state string) string {
	switch state {
	case "ready":
		return "38;5;114"
	case "stopped":
		return "2;38;5;245"
	case "failed":
		return "38;5;203"
	default:
		return "38;5;81"
	}
}

func rowTextWidth(values []string) int {
	width := 0
	for _, value := range values {
		if len(value) > width {
			width = len(value)
		}
	}
	return width
}

func helpCommandWidth(sections []helpCommandSection) int {
	width := 0
	for _, section := range sections {
		for _, row := range section.rows {
			if len(row.key) > width {
				width = len(row.key)
			}
		}
	}
	return width
}

func wrapHelpText(value string, width int) []string {
	if width <= 0 || len(value) <= width {
		return []string{value}
	}

	var lines []string
	for _, paragraph := range strings.Split(value, "\n") {
		if strings.TrimSpace(paragraph) == "" {
			lines = append(lines, "")
			continue
		}

		words := strings.Fields(paragraph)
		current := words[0]
		for len(current) > width {
			lines = append(lines, current[:width])
			current = current[width:]
		}
		for _, word := range words[1:] {
			if len(current)+1+len(word) <= width {
				current += " " + word
				continue
			}
			lines = append(lines, current)
			current = word
			for len(current) > width {
				lines = append(lines, current[:width])
				current = current[width:]
			}
		}
		lines = append(lines, current)
	}
	if len(lines) == 0 {
		return []string{""}
	}
	return lines
}

func padRight(value string, width int) string {
	if len(value) >= width {
		return value
	}
	return value + strings.Repeat(" ", width-len(value))
}

func clamp(value int, minValue int, maxValue int) int {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}

func (r renderer) formatKey(key string) string {
	return r.paint("2;38;5;245", key)
}

func (r renderer) formatService(service string) string {
	switch service {
	case "web":
		return r.paint("38;5;81", service)
	case "api":
		return r.paint("38;5;114", service)
	case "db":
		return r.paint("38;5;222", service)
	case "watch":
		return r.paint("38;5;147", service)
	default:
		return r.paint("38;5;245", service)
	}
}

func (r renderer) formatLogoLine(_ int, line string) string {
	if !r.styled {
		return line
	}

	var out strings.Builder
	activeColor := ""
	for i := 0; i < len(line); i++ {
		color := logoGlyphColor(line[i])
		if color != activeColor {
			if activeColor != "" {
				out.WriteString("\033[0m")
			}
			if color != "" {
				out.WriteString("\033[")
				out.WriteString(color)
				out.WriteString("m")
			}
			activeColor = color
		}
		out.WriteByte(line[i])
	}
	if activeColor != "" {
		out.WriteString("\033[0m")
	}
	return out.String()
}

func (r renderer) formatLogoPacmanLine(line string, position int, step int) string {
	width := len(line)
	if position >= width {
		return r.formatLogoLine(0, line)
	}

	var out strings.Builder
	if position > 0 {
		out.WriteString(r.formatLogoLine(0, visiblePrefix(line, position)))
	}
	start := clamp(position, 0, width)
	for column := start; column < width; column++ {
		switch {
		case column == position:
			out.WriteString(r.formatLogoChomper(pacmanMouth(step, "right")))
		case isCandyPosition(column - position):
			out.WriteString(r.formatLogoPellet())
		default:
			out.WriteByte(' ')
		}
	}
	return out.String()
}

func (r renderer) formatLogoChomper(ch byte) string {
	return r.paint("1;38;5;226", string(ch))
}

func (r renderer) formatLogoPellet() string {
	return r.paint("2;38;5;220", "o")
}

func logoGlyphColor(ch byte) string {
	switch ch {
	case '_':
		return "2;38;5;245"
	case 'o', 'O', '0':
		return "38;5;220"
	default:
		return ""
	}
}

func (r renderer) formatValue(row outputRow) string {
	if !r.styled {
		return row.value
	}
	switch row.key {
	case "app", "api":
		return r.paint("38;5;81", row.value)
	case "error":
		return r.paint("38;5;203", row.value)
	default:
		return row.value
	}
}

func (r renderer) paint(code string, value string) string {
	if !r.styled {
		return value
	}
	return "\033[" + code + "m" + value + "\033[0m"
}

func rowKeyWidth(rows []outputRow) int {
	width := 0
	for _, row := range rows {
		if len(row.key) > width {
			width = len(row.key)
		}
	}
	return width
}

func logoLines(logo string) []string {
	logo = strings.TrimRight(logo, "\n")
	if strings.TrimSpace(logo) == "" {
		return nil
	}
	return strings.Split(logo, "\n")
}

func maxLineWidth(lines []string) int {
	width := 0
	for _, line := range lines {
		if len(line) > width {
			width = len(line)
		}
	}
	return width
}

func visiblePrefix(value string, width int) string {
	if width <= 0 {
		return ""
	}
	if width >= len(value) {
		return value
	}
	return value[:width]
}

func streamWatchOutput(input io.Reader, r renderer, logSink *devLogSink, stream string, wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(input)
	scanner.Buffer(make([]byte, 1024), 1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || line == "Watch enabled" {
			continue
		}
		entry := newStructuredLogEntry("compose-watch", stream, "watch", line)
		logSink.WriteEntry(entry)
		r.LogEntry(entry)
	}
}

func streamLogOutput(input io.Reader, r renderer, logSink *devLogSink, stream string, wg *sync.WaitGroup) {
	streamLogOutputWithQuery(input, r, logSink, stream, logQuery{}, wg)
}

func streamLogOutputWithQuery(input io.Reader, r renderer, logSink *devLogSink, stream string, query logQuery, wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(input)
	scanner.Buffer(make([]byte, 1024), 1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		service, message := parseComposeLogLine(line)
		entry := newStructuredLogEntry("compose-log", stream, service, message)
		logSink.WriteEntry(entry)
		if logEntryMatchesQuery(entry, query) {
			r.LogEntry(entry)
		}
	}
}

func newStructuredLogEntry(source string, stream string, service string, message string) structuredLogEntry {
	return structuredLogEntry{
		Time:    time.Now().UTC().Format(time.RFC3339Nano),
		Source:  source,
		Stream:  stream,
		Service: service,
		Message: message,
	}
}

func openDevLogSink(path string) (*devLogSink, error) {
	return openDevLogSinkWithFlags(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
}

func openAppendDevLogSink(path string) (*devLogSink, error) {
	return openDevLogSinkWithFlags(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND)
}

func openDevLogSinkWithFlags(path string, flags int) (*devLogSink, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}
	file, err := os.OpenFile(path, flags, 0644)
	if err != nil {
		return nil, err
	}
	return &devLogSink{file: file, encoder: json.NewEncoder(file)}, nil
}

func (s *devLogSink) Close() error {
	if s == nil || s.file == nil {
		return nil
	}
	return s.file.Close()
}

func (s *devLogSink) Write(source string, stream string, service string, message string) {
	s.WriteEntry(newStructuredLogEntry(source, stream, service, message))
}

func (s *devLogSink) WriteEntry(entry structuredLogEntry) {
	if s == nil || s.encoder == nil || strings.TrimSpace(entry.Message) == "" {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	_ = s.encoder.Encode(entry)
}
