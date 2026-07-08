package cli

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func (a app) commandAuditFlow(parts []string) error {
	switch {
	case len(parts) == 0:
		_, err := a.commandAuditStage()
		return err
	case len(parts) == 1 && parts[0] == "resolve":
		auditResult, err := a.commandAuditStage()
		if err != nil {
			return err
		}
		if !auditResult.reportsReady {
			return nil
		}
		_, err = a.commandResolveStage()
		return err
	case len(parts) == 2 && parts[0] == "resolve" && parts[1] == "fix":
		auditResult, err := a.commandAuditStage()
		if err != nil {
			return err
		}
		if !auditResult.reportsReady {
			return nil
		}
		resolveResult, err := a.commandResolveStage()
		if err != nil {
			return err
		}
		if !resolveResult.ready {
			return nil
		}
		return a.commandFix()
	default:
		return errors.New("usage: carbide audit [resolve [fix]]")
	}
}

func (a app) commandResolveFlow(parts []string) error {
	switch {
	case len(parts) == 0:
		_, err := a.commandResolveStage()
		return err
	case len(parts) == 1 && parts[0] == "fix":
		resolveResult, err := a.commandResolveStage()
		if err != nil {
			return err
		}
		if !resolveResult.ready {
			return nil
		}
		return a.commandFix()
	default:
		return errors.New("usage: carbide resolve [fix]")
	}
}

func (a app) commandAuditStage() (auditStageResult, error) {
	result := auditStageResult{
		root:       auditRootDir,
		reportDir:  filepath.Join(auditRootDir, auditReportDirName),
		starterDir: filepath.Join(auditRootDir, auditStarterDirName),
	}
	if !isFile(projectConfigPath) {
		return result, errors.New("carbide audit requires carbide.toml")
	}

	pwd, err := os.Getwd()
	if err != nil {
		return result, err
	}
	name, slug, err := currentProjectIdentity(pwd)
	if err != nil {
		return result, err
	}

	if err := os.RemoveAll(result.root); err != nil {
		return result, err
	}
	if err := os.MkdirAll(result.reportDir, 0755); err != nil {
		return result, err
	}
	if err := a.copyScaffold(result.starterDir, name, slug); err != nil {
		return result, err
	}

	specs := auditSpecs()
	result.reportCount = len(specs)
	for _, spec := range specs {
		if err := os.WriteFile(filepath.Join(result.reportDir, spec.fileName), []byte(pendingAuditReport(spec)), 0644); err != nil {
			return result, err
		}
	}

	if !canRunCodexAutomation(a.stdout) {
		newRenderer(a.stdout).Message(
			"Carbide audit",
			"workspace prepared",
			outputRow{"path", result.root},
			outputRow{"starter", result.starterDir},
			outputRow{"reports", fmt.Sprintf("%d pending files", result.reportCount)},
			outputRow{"next", "run carbide resolve in a terminal with codex installed"},
		)
		return result, nil
	}

	failures := a.runAuditReportAgents(pwd, result.root, specs)
	if len(failures) == 0 {
		result.reportsReady = true
		newRenderer(a.stdout).Message(
			"Carbide audit",
			"reports created",
			outputRow{"path", result.root},
			outputRow{"starter", result.starterDir},
			outputRow{"reports", fmt.Sprintf("%d markdown files", result.reportCount)},
			outputRow{"next", "carbide resolve"},
		)
		return result, nil
	}

	newRenderer(a.stdout).Message(
		"Carbide audit",
		"reports partially generated",
		outputRow{"path", result.root},
		outputRow{"reports", fmt.Sprintf("%d attempted, %d failed", result.reportCount, len(failures))},
		outputRow{"next", "rerun carbide audit in a terminal with codex installed"},
	)
	return result, nil
}

func (a app) commandResolveStage() (resolveStageResult, error) {
	result := resolveStageResult{planPath: filepath.Join(auditRootDir, auditPlanFileName)}
	if !isFile(projectConfigPath) {
		return result, errors.New("carbide resolve requires carbide.toml")
	}
	if !isDir(auditRootDir) {
		return result, errors.New("run carbide audit first")
	}

	pwd, err := os.Getwd()
	if err != nil {
		return result, err
	}
	specs := auditSpecs()
	reportDir := filepath.Join(auditRootDir, auditReportDirName)
	if missing := missingAuditReports(reportDir, specs); len(missing) > 0 {
		return result, fmt.Errorf("run carbide audit first; missing %s", strings.Join(missing, ", "))
	}

	if auditReportsPending(reportDir, specs) {
		if canRunCodexAutomation(a.stdout) {
			failures := a.runAuditReportAgents(pwd, auditRootDir, specs)
			if len(failures) > 0 || auditReportsPending(reportDir, specs) {
				if err := os.WriteFile(result.planPath, []byte(pendingResolvePlan("audit reports are still pending")), 0644); err != nil {
					return result, err
				}
				newRenderer(a.stdout).Message(
					"Carbide resolve",
					"reports pending",
					outputRow{"plan", result.planPath},
					outputRow{"next", "rerun carbide resolve after report agents finish cleanly"},
				)
				return result, nil
			}
		} else {
			if err := os.WriteFile(result.planPath, []byte(pendingResolvePlan("codex is required to resolve pending audit reports automatically")), 0644); err != nil {
				return result, err
			}
			newRenderer(a.stdout).Message(
				"Carbide resolve",
				"plan stub created",
				outputRow{"plan", result.planPath},
				outputRow{"next", "run carbide resolve in a terminal with codex installed"},
			)
			return result, nil
		}
	}

	if !canRunCodexAutomation(a.stdout) {
		if err := os.WriteFile(result.planPath, []byte(pendingResolvePlan("codex is required to synthesize the audit reports into a plan automatically")), 0644); err != nil {
			return result, err
		}
		newRenderer(a.stdout).Message(
			"Carbide resolve",
			"plan stub created",
			outputRow{"plan", result.planPath},
			outputRow{"next", "run carbide resolve in a terminal with codex installed"},
		)
		return result, nil
	}

	draft, err := a.generateResolvePlan(pwd, auditRootDir, nil)
	if err != nil {
		return result, err
	}
	status, questions := parseResolvePlan(draft)
	if status == "needs-clarification" && len(questions) > 0 && canPromptForClarifications(a.stdout) {
		answers, answered := promptForClarifications(a.stdout, questions)
		if answered {
			draft, err = a.generateResolvePlan(pwd, auditRootDir, answers)
			if err != nil {
				return result, err
			}
			status, questions = parseResolvePlan(draft)
		}
	}

	if err := os.WriteFile(result.planPath, []byte(draft), 0644); err != nil {
		return result, err
	}

	if status == "ready" {
		result.ready = true
		newRenderer(a.stdout).Message(
			"Carbide resolve",
			"plan created",
			outputRow{"plan", result.planPath},
			outputRow{"next", "carbide fix"},
		)
		return result, nil
	}

	newRenderer(a.stdout).Message(
		"Carbide resolve",
		"clarification needed",
		outputRow{"plan", result.planPath},
		outputRow{"questions", strconv.Itoa(len(questions))},
		outputRow{"next", "answer the questions and rerun carbide resolve"},
	)
	return result, nil
}

func (a app) commandFix() error {
	if !isFile(projectConfigPath) {
		return errors.New("carbide fix requires carbide.toml")
	}
	planPath := filepath.Join(auditRootDir, auditPlanFileName)
	if !isFile(planPath) {
		return errors.New("run carbide resolve first")
	}
	status, questions := parseResolvePlan(readFileString(planPath))
	if status != "ready" {
		if len(questions) > 0 {
			return errors.New("carbide fix requires a ready .audit/plan.md; resolve still needs clarification")
		}
		return errors.New("carbide fix requires a ready .audit/plan.md")
	}
	if !canRunCodexAutomation(a.stdout) {
		return errors.New("codex is required to implement the audit plan automatically")
	}

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	summaryPath := filepath.Join(auditRootDir, "fix.md")
	newRenderer(a.stdout).Message(
		"Carbide fix",
		"launching Codex",
		outputRow{"plan", planPath},
		outputRow{"summary", summaryPath},
	)
	summary, err := a.runFixCodex(pwd, auditRootDir)
	if err != nil {
		return err
	}
	if err := os.WriteFile(summaryPath, []byte(summary), 0644); err != nil {
		return err
	}
	newRenderer(a.stdout).Message(
		"Carbide fix",
		"plan implemented",
		outputRow{"summary", summaryPath},
		outputRow{"next", "run carbide health"},
	)
	return nil
}

func (a app) runAuditReportAgents(projectPath string, auditRoot string, specs []auditSpec) []string {
	var wg sync.WaitGroup
	var mu sync.Mutex
	failures := make([]string, 0)

	for _, spec := range specs {
		spec := spec
		wg.Add(1)
		go func() {
			defer wg.Done()
			reportPath := filepath.Join(auditRoot, auditReportDirName, spec.fileName)
			if err := runCodexExecLastMessage(projectPath, reportPath, auditReportPrompt(auditRoot, spec), "read-only"); err != nil {
				mu.Lock()
				failures = append(failures, fmt.Sprintf("%s: %v", spec.ref, err))
				mu.Unlock()
				appendAutomationFailure(reportPath, err)
			}
		}()
	}

	wg.Wait()
	sort.Strings(failures)
	return failures
}

func (a app) generateResolvePlan(projectPath string, auditRoot string, answers []auditClarification) (string, error) {
	tmp, err := os.CreateTemp("", "carbide-resolve-*.md")
	if err != nil {
		return "", err
	}
	tmpPath := tmp.Name()
	if err := tmp.Close(); err != nil {
		return "", err
	}
	defer os.Remove(tmpPath)

	if err := runCodexExecLastMessage(projectPath, tmpPath, resolveCodexPrompt(auditRoot, answers), "read-only"); err != nil {
		return "", err
	}
	return readFileString(tmpPath), nil
}

func (a app) runFixCodex(projectPath string, auditRoot string) (string, error) {
	tmp, err := os.CreateTemp("", "carbide-fix-*.md")
	if err != nil {
		return "", err
	}
	tmpPath := tmp.Name()
	if err := tmp.Close(); err != nil {
		return "", err
	}
	defer os.Remove(tmpPath)

	if err := runCodexExecLastMessage(projectPath, tmpPath, fixCodexPrompt(auditRoot), "workspace-write"); err != nil {
		return "", err
	}
	return readFileString(tmpPath), nil
}

func runCodexExecLastMessage(projectPath string, outputPath string, prompt string, sandbox string) error {
	args := []string{
		"exec",
		"--cd", projectPath,
		"--sandbox", sandbox,
		"--color", "never",
		"--ephemeral",
		"--output-last-message", outputPath,
		prompt,
	}
	cmd := exec.Command("codex", args...)
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		detail := strings.TrimSpace(output.String())
		if detail == "" {
			return err
		}
		return fmt.Errorf("%w: %s", err, detail)
	}
	return nil
}

func auditSpecs() []auditSpec {
	return []auditSpec{
		{ref: "Law 1", title: "One App Repo", fileName: "law-1-one-app-repo.md", description: "The app stays in one repo."},
		{ref: "Law 2", title: "Root Runtime Directories", fileName: "law-2-root-runtime-directories.md", description: "The root runtime directories are web/, api/, and db/."},
		{ref: "Law 3", title: "Checked-In Runtime Contracts", fileName: "law-3-checked-in-runtime-contracts.md", description: "The checked-in contracts are carbide.toml and docker-compose.yml."},
		{ref: "Law 4", title: "Same-Origin Browser Flow", fileName: "law-4-same-origin-browser-flow.md", description: "Browser traffic stays same-origin: web proxies /api to api."},
		{ref: "Law 5", title: "Postgres Is Required", fileName: "law-5-postgres-is-required.md", description: "Postgres is the required durable database."},
		{ref: "Law 6", title: "Preview Before Apply", fileName: "law-6-preview-before-apply.md", description: "Deploy keeps the preview-before-apply rule."},
		{ref: "Law 7", title: "Secrets Are Never Printed", fileName: "law-7-secrets-are-never-printed.md", description: "Secrets never appear in logs, docs, or command output."},
		{ref: "Taste 1", title: "Starter Stack", fileName: "taste-1-starter-stack.md", description: "The current preferred starter stack."},
		{ref: "Taste 2", title: "Runtime Pins", fileName: "taste-2-runtime-pins.md", description: "The current runtime pin and baseline posture."},
		{ref: "Taste 3", title: "Starter Product Surface", fileName: "taste-3-starter-product-surface.md", description: "The current starter auth, dashboard, and first-run surface."},
		{ref: "Taste 4", title: "Frontend Organization", fileName: "taste-4-frontend-organization.md", description: "The current React, Tailwind, and component-organization taste."},
		{ref: "Taste 5", title: "Docs And Examples", fileName: "taste-5-docs-and-examples.md", description: "The current docs and example-writing taste."},
		{ref: "Taste 6", title: "CLI And Audit Reporting", fileName: "taste-6-cli-and-audit-reporting.md", description: "The current CLI copy, audit flow, and operator-facing reporting taste."},
	}
}

func pendingAuditReport(spec auditSpec) string {
	rows := []string{
		"status: pending",
		"",
		fmt.Sprintf("# %s. %s", spec.ref, spec.title),
		"",
		spec.description,
		"",
		"This report has not been generated yet.",
		"",
		"Next:",
		"- Run `carbide audit` in an interactive terminal with `codex` installed.",
	}
	return strings.Join(rows, "\n")
}

func pendingResolvePlan(reason string) string {
	rows := []string{
		"status: pending",
		"",
		"# Carbide Resolve Plan",
		"",
		"## State",
		"",
		reason,
		"",
		"## Next",
		"",
		"- Run `carbide audit` if the report set is incomplete.",
		"- Run `carbide resolve` in an interactive terminal with `codex` installed.",
		"- Run `carbide fix` only after this file switches to `status: ready`.",
	}
	return strings.Join(rows, "\n")
}

func parseResolvePlan(content string) (string, []string) {
	status := ""
	questions := make([]string, 0)
	inQuestions := false

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if status == "" && strings.HasPrefix(strings.ToLower(line), "status:") {
			status = strings.TrimSpace(strings.TrimPrefix(strings.ToLower(line), "status:"))
			continue
		}
		if strings.EqualFold(line, "## Questions") {
			inQuestions = true
			continue
		}
		if strings.HasPrefix(line, "## ") {
			inQuestions = false
			continue
		}
		if inQuestions && strings.HasPrefix(line, "- ") {
			question := strings.TrimSpace(strings.TrimPrefix(line, "- "))
			if question != "" {
				questions = append(questions, question)
			}
		}
	}

	if status == "" {
		status = "pending"
	}
	return status, questions
}

func missingAuditReports(reportDir string, specs []auditSpec) []string {
	missing := make([]string, 0)
	for _, spec := range specs {
		if !isFile(filepath.Join(reportDir, spec.fileName)) {
			missing = append(missing, spec.fileName)
		}
	}
	return missing
}

func auditReportsPending(reportDir string, specs []auditSpec) bool {
	for _, spec := range specs {
		content := readFileString(filepath.Join(reportDir, spec.fileName))
		if !strings.HasPrefix(strings.ToLower(strings.TrimSpace(content)), "status: complete") {
			return true
		}
	}
	return false
}

func canRunCodexAutomation(w io.Writer) bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("CARBIDE_AUDIT_AUTOMATION"))) {
	case "0", "false", "off":
		return false
	}
	_, err := exec.LookPath("codex")
	if err != nil {
		return false
	}
	return true
}

func canPromptForClarifications(w io.Writer) bool {
	return isTerminalOutput(w) && isTerminalFile(os.Stdin)
}

func promptForClarifications(w io.Writer, questions []string) ([]auditClarification, bool) {
	if len(questions) == 0 {
		return nil, false
	}
	reader := bufio.NewReader(os.Stdin)
	answers := make([]auditClarification, 0, len(questions))
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Carbide resolve needs clarification:")
	for i, question := range questions {
		fmt.Fprintf(w, "%d. %s\n", i+1, question)
		fmt.Fprint(w, "> ")
		answer, err := reader.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, false
		}
		answer = strings.TrimSpace(answer)
		if answer == "" {
			return nil, false
		}
		answers = append(answers, auditClarification{question: question, answer: answer})
		if errors.Is(err, io.EOF) {
			break
		}
	}
	return answers, len(answers) == len(questions)
}

func appendAutomationFailure(path string, err error) {
	lines := []string{
		"status: failed",
		"",
		"# Automation Failure",
		"",
		fmt.Sprintf("- error: %v", err),
		"- rerun the command in an interactive terminal with `codex` installed.",
	}
	_ = os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0644)
}

func currentProjectIdentity(projectPath string) (string, string, error) {
	metadata, err := readProjectMetadata(filepath.Join(projectPath, projectConfigPath))
	if err != nil {
		return "", "", err
	}
	name := strings.TrimSpace(metadata.name)
	if name == "" || isTemplatePlaceholder(name) {
		name = projectDisplayName(filepath.Base(projectPath))
	}
	slug := normalizeComposeProjectName(metadata.slug)
	if slug == "" {
		slug = projectSlug(filepath.Base(projectPath))
	}
	if slug == "" {
		slug = "carbide-app"
	}
	return name, slug, nil
}

func auditReportPrompt(auditRoot string, spec auditSpec) string {
	rows := []string{
		"Audit one Carbide contract slice and write only the requested markdown report.",
		"",
		fmt.Sprintf("Focus only on `%s. %s`.", spec.ref, spec.title),
		fmt.Sprintf("Report path: `%s`", filepath.ToSlash(filepath.Join(auditRoot, auditReportDirName, spec.fileName))),
		fmt.Sprintf("Starter reference: `%s`", filepath.ToSlash(filepath.Join(auditRoot, auditStarterDirName))),
		"",
		"Read the current project and compare it against Carbide's current starter taste only for this one topic.",
		"Do not edit project files. Only produce the report markdown.",
		"",
		"Output format:",
		"status: complete",
		"",
		fmt.Sprintf("# %s. %s", spec.ref, spec.title),
		"",
		"## Verdict",
		"- compliant | drifted | not-applicable",
		"",
		"## Evidence",
		"- concise repo evidence with file paths",
		"",
		"## Gaps",
		"- concrete missing or drifted behavior",
		"",
		"## Recommended Changes",
		"- practical next steps",
		"",
		"## Questions",
		"- include only if a real user clarification is required",
	}
	return strings.Join(rows, "\n")
}

func resolveCodexPrompt(auditRoot string, answers []auditClarification) string {
	rows := []string{
		"Resolve the Carbide audit reports into one implementation plan.",
		"",
		fmt.Sprintf("Read every markdown file under `%s`.", filepath.ToSlash(filepath.Join(auditRoot, auditReportDirName))),
		fmt.Sprintf("Use `%s` only as a starter reference, not as framework-managed truth.", filepath.ToSlash(filepath.Join(auditRoot, auditStarterDirName))),
		"",
		"Do not edit project files. Only write the plan as markdown.",
		"",
		"Output format:",
		"status: ready | needs-clarification",
		"",
		"# Carbide Resolve Plan",
		"",
		"## Summary",
		"- one short paragraph",
		"",
		"## Laws",
		"- Law N: keep/fix plus brief rationale",
		"",
		"## Taste",
		"- Taste N: adopt/ignore plus brief rationale",
		"",
		"## Implementation Steps",
		"1. concrete edit sequence",
		"2. concrete verification sequence",
		"",
		"## Questions",
		"- include only if user clarification is materially required",
		"",
		"## Done When",
		"- concrete completion checks",
	}
	if len(answers) > 0 {
		rows = append(rows, "", "Clarifications already provided by the user:")
		for _, answer := range answers {
			rows = append(rows, fmt.Sprintf("- %s => %s", answer.question, answer.answer))
		}
	}
	return strings.Join(rows, "\n")
}

func fixCodexPrompt(auditRoot string) string {
	rows := []string{
		"Implement the latest Carbide resolve plan.",
		"",
		fmt.Sprintf("Read `%s` first.", filepath.ToSlash(filepath.Join(auditRoot, auditPlanFileName))),
		fmt.Sprintf("Use `%s` and `%s` as supporting context.", filepath.ToSlash(filepath.Join(auditRoot, auditReportDirName)), filepath.ToSlash(filepath.Join(auditRoot, auditStarterDirName))),
		"",
		"Rules:",
		"- edit project files intentionally",
		"- do not rewrite the starter into the app blindly",
		"- do not delete .audit inputs",
		"- keep the repo inside Carbide law compliance",
		"- run `carbide health` at the end",
		"- run `carbide health runtime` when runtime behavior or containers changed",
		"- run any relevant app-specific build or test commands",
		"",
		"Return a short markdown summary with:",
		"- files changed",
		"- laws fixed",
		"- taste changes adopted",
		"- checks run",
		"- remaining risk",
	}
	return strings.Join(rows, "\n")
}
