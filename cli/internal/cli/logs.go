package cli

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

func (a app) commandLogs(args []string) error {
	if !isFile("carbide.toml") {
		return errors.New("run this inside a Carbide project")
	}
	query, err := parseLogQuery(args)
	if err != nil {
		return err
	}
	entries, err := readStructuredLogEntries(devLogPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errors.New("no dev logs found; run carbide run dev first")
		}
		return err
	}
	entries = filterLogEntries(entries, query)
	entries = limitLogEntries(entries, query.limit)

	if query.json {
		encoder := json.NewEncoder(a.stdout)
		for _, entry := range entries {
			if err := encoder.Encode(entry); err != nil {
				return err
			}
		}
		return nil
	}

	r := newRenderer(a.stdout)
	for _, entry := range entries {
		r.LogEntry(entry)
	}
	return nil
}

func (a app) commandFollowLogs(args []string) error {
	if !isFile("carbide.toml") {
		return errors.New("run this inside a Carbide project")
	}
	query, err := parseLogQuery(args)
	if err != nil {
		return err
	}
	if query.json {
		return errors.New("carbide follow logs does not support json")
	}
	if query.limit != 80 {
		return errors.New("carbide follow logs does not support limit")
	}

	compose, err := findCompose()
	if err != nil {
		return err
	}
	env := setEnv(os.Environ(), "COMPOSE_MENU", "false")
	env = composeEnv(env)
	logSink, err := openAppendDevLogSink(devLogPath)
	if err != nil {
		return err
	}
	defer logSink.Close()

	var streams sync.WaitGroup
	results := make(chan processResult, 1)
	process, err := a.startComposeStream(
		"logs",
		compose,
		env,
		composeLogsArgs(compose),
		func(input io.Reader, r renderer, sink *devLogSink, stream string, wg *sync.WaitGroup) {
			streamLogOutputWithQuery(input, r, sink, stream, query, wg)
		},
		logSink,
		&streams,
		results,
	)
	if err != nil {
		return err
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signals)

	select {
	case sig := <-signals:
		stopProcesses([]runningProcess{process}, sig)
		waitForProcesses(1, []runningProcess{process}, results, 5*time.Second)
	case result := <-results:
		if result.err != nil {
			streams.Wait()
			return fmt.Errorf("Docker Compose %s failed: %w", result.name, result.err)
		}
	}

	streams.Wait()
	return nil
}

func entryTimestamp(entry structuredLogEntry) time.Time {
	timestamp, err := time.Parse(time.RFC3339Nano, entry.Time)
	if err != nil {
		return time.Now()
	}
	return timestamp
}

func parseLogQuery(args []string) (logQuery, error) {
	query := logQuery{limit: 80}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "service":
			i++
			if i >= len(args) || args[i] == "" {
				return query, errors.New("usage: carbide logs [service <name>] [containing <text>] [limit <count>] [json]")
			}
			query.service = args[i]
		case "containing":
			i++
			if i >= len(args) || args[i] == "" {
				return query, errors.New("usage: carbide logs [service <name>] [containing <text>] [limit <count>] [json]")
			}
			query.contains = args[i]
		case "limit":
			i++
			if i >= len(args) {
				return query, errors.New("usage: carbide logs [service <name>] [containing <text>] [limit <count>] [json]")
			}
			limit, err := strconv.Atoi(args[i])
			if err != nil || limit < 1 {
				return query, errors.New("log limit must be a positive number")
			}
			query.limit = limit
		case "json":
			query.json = true
		default:
			return query, fmt.Errorf("unknown logs option: %s", args[i])
		}
	}
	return query, nil
}

func readStructuredLogEntries(path string) ([]structuredLogEntry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []structuredLogEntry
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1024), 1024*1024)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var entry structuredLogEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			return nil, fmt.Errorf("invalid structured log line %d: %w", lineNumber, err)
		}
		entries = append(entries, entry)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return entries, nil
}

func filterLogEntries(entries []structuredLogEntry, query logQuery) []structuredLogEntry {
	filtered := make([]structuredLogEntry, 0, len(entries))
	for _, entry := range entries {
		if logEntryMatchesQuery(entry, query) {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

func logEntryMatchesQuery(entry structuredLogEntry, query logQuery) bool {
	if query.service != "" && entry.Service != query.service {
		return false
	}
	if query.contains != "" && !strings.Contains(strings.ToLower(entry.Message), strings.ToLower(query.contains)) {
		return false
	}
	return true
}

func limitLogEntries(entries []structuredLogEntry, limit int) []structuredLogEntry {
	if limit < 1 || len(entries) <= limit {
		return entries
	}
	return entries[len(entries)-limit:]
}

func carbideLogo() string {
	return defaultLogoText
}
