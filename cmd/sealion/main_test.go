package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func TestProjectSlug(t *testing.T) {
	tests := map[string]string{
		"Demo":           "demo",
		"my_app.test":    "my-app-test",
		"  Weird Name  ": "weird-name",
		"already--clean": "already-clean",
		"___":            "",
	}

	for input, want := range tests {
		if got := projectSlug(input); got != want {
			t.Fatalf("projectSlug(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestEnsureProjectName(t *testing.T) {
	valid := []string{"demo", "demo_app", "demo-app", "demo.app", "Demo1"}
	for _, name := range valid {
		if err := ensureProjectName(name); err != nil {
			t.Fatalf("ensureProjectName(%q) returned %v", name, err)
		}
	}

	invalid := []string{"", ".hidden", "two words", "nested/app", "bad*name"}
	for _, name := range invalid {
		if err := ensureProjectName(name); err == nil {
			t.Fatalf("ensureProjectName(%q) should fail", name)
		}
	}
}

func TestValidatePort(t *testing.T) {
	for _, value := range []string{"", "0", "65536", "abc"} {
		if _, err := validatePort(value); err == nil {
			t.Fatalf("validatePort(%q) should fail", value)
		}
	}

	got, err := validatePort("8080")
	if err != nil {
		t.Fatalf("validatePort returned %v", err)
	}
	if got != 8080 {
		t.Fatalf("validatePort returned %d, want 8080", got)
	}
}

func TestRendererPlainOutput(t *testing.T) {
	var out bytes.Buffer
	newRenderer(&out).Message(
		"Sealion",
		"project created",
		outputRow{"path", "/tmp/demo"},
		outputRow{"next", "cd demo"},
		outputRow{"", "sealion run dev"},
	)

	want := "Sealion\nproject created\n\npath  /tmp/demo\nnext  cd demo\n      sealion run dev\n"
	if out.String() != want {
		t.Fatalf("renderer output = %q, want %q", out.String(), want)
	}
}

func TestRendererIndentsMultilineValues(t *testing.T) {
	var out bytes.Buffer
	newRenderer(&out).Rows(outputRow{"error", "first line\nsecond line"})

	want := "error  first line\n       second line\n"
	if out.String() != want {
		t.Fatalf("renderer output = %q, want %q", out.String(), want)
	}
}

func TestStreamWatchOutputFiltersNoise(t *testing.T) {
	var out bytes.Buffer
	var wg sync.WaitGroup
	wg.Add(1)
	streamWatchOutput(strings.NewReader("Watch enabled\n\nrebuilt backend\n"), newRenderer(&out), nil, "stdout", &wg)
	wg.Wait()

	want := "watch      rebuilt backend\n"
	if out.String() != want {
		t.Fatalf("watch output = %q, want %q", out.String(), want)
	}
}

func TestStreamLogOutputParsesComposeServices(t *testing.T) {
	var out bytes.Buffer
	var wg sync.WaitGroup
	wg.Add(1)
	streamLogOutput(
		strings.NewReader("backend-1  | GET /health\nfrontend-1 | listening\ndemo-db-1 | ready\n"),
		newRenderer(&out),
		nil,
		"stdout",
		&wg,
	)
	wg.Wait()

	want := "backend    GET /health\nfrontend   listening\ndb         ready\n"
	if out.String() != want {
		t.Fatalf("log output = %q, want %q", out.String(), want)
	}
}

func TestStreamLogOutputWritesStructuredJSON(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "dev.jsonl")
	sink, err := openDevLogSink(logPath)
	if err != nil {
		t.Fatalf("openDevLogSink returned %v", err)
	}

	var out bytes.Buffer
	var wg sync.WaitGroup
	wg.Add(1)
	streamLogOutput(strings.NewReader("backend-1 | GET /health\n"), newRenderer(&out), sink, "stdout", &wg)
	wg.Wait()
	if err := sink.Close(); err != nil {
		t.Fatalf("Close returned %v", err)
	}

	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("ReadFile returned %v", err)
	}
	text := string(data)
	for _, want := range []string{
		`"source":"compose-log"`,
		`"stream":"stdout"`,
		`"service":"backend"`,
		`"message":"GET /health"`,
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("structured log %q missing %s", text, want)
		}
	}
}

func TestParseLogQuery(t *testing.T) {
	query, err := parseLogQuery([]string{"service", "backend", "containing", "health", "limit", "5", "json"})
	if err != nil {
		t.Fatalf("parseLogQuery returned %v", err)
	}
	if query.service != "backend" || query.contains != "health" || query.limit != 5 || !query.json {
		t.Fatalf("query = %#v", query)
	}
}

func TestFilterAndLimitLogEntries(t *testing.T) {
	entries := []structuredLogEntry{
		{Service: "frontend", Message: "listening"},
		{Service: "backend", Message: "GET /health"},
		{Service: "backend", Message: "POST /api/login"},
	}

	filtered := filterLogEntries(entries, logQuery{service: "backend", contains: "api"})
	if len(filtered) != 1 || filtered[0].Message != "POST /api/login" {
		t.Fatalf("filtered = %#v", filtered)
	}

	limited := limitLogEntries(entries, 2)
	if len(limited) != 2 || limited[0].Message != "GET /health" {
		t.Fatalf("limited = %#v", limited)
	}
}
