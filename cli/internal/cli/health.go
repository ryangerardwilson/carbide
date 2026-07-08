package cli

import (
	"errors"
	"fmt"
)

type healthResult struct {
	check  string
	status string
	detail string
}

func (a app) commandHealth() error {
	results := a.projectHealthResults()
	results = append(results, healthResult{"runtime", "skip", "run carbide health runtime"})
	return a.renderHealthResults("app laws", results)
}

func (a app) commandHealthJSON() error {
	results := a.projectHealthResults()
	results = append(results, healthResult{"runtime", "skip", "run carbide health runtime"})
	return a.renderHealthResultsJSON("health", results, []string{"carbide health runtime"})
}

func (a app) commandHealthEnv() error {
	if !isFile("carbide.toml") {
		return errors.New("run this inside a Carbide project")
	}

	report, err := inspectEnvContract()
	if err != nil {
		return err
	}

	status := "ok"
	if len(report.missingRequired) > 0 || len(report.warnings) > 0 {
		status = "needs attention"
	}
	envFile := ".env not found; local defaults active"
	if report.envFileFound {
		envFile = ".env found"
	}

	r := newRenderer(a.stdout)
	r.Title("Carbide health", "environment contract")
	r.Rows(
		outputRow{"contract", projectConfigPath},
		outputRow{"env", envFile},
		outputRow{"status", status},
		outputRow{"required", fmt.Sprintf("%d missing", len(report.missingRequired))},
		outputRow{"secrets", fmt.Sprintf("%d declared", report.secretCount)},
		outputRow{"browser", fmt.Sprintf("%d exposed", report.browserCount)},
		outputRow{"framework", fmt.Sprintf("%d owned", report.frameworkCount)},
	)
	for _, name := range report.missingRequired {
		r.Row(outputRow{"missing", name})
	}
	for _, warning := range report.warnings {
		r.Row(outputRow{"warning", warning})
	}
	if len(report.missingRequired) > 0 {
		return fmt.Errorf("environment contract has %d missing required value(s)", len(report.missingRequired))
	}
	return nil
}

func (a app) commandHealthEnvJSON() error {
	if !isFile("carbide.toml") {
		return errors.New("run this inside a Carbide project")
	}

	report, err := inspectEnvContract()
	if err != nil {
		return err
	}
	env := envJSONFromReport(report)
	ok := len(report.missingRequired) == 0 && len(report.warnings) == 0
	out := a.baseJSONReport("health env")
	out.OK = ok
	out.Env = &env
	if !ok {
		out.Errors = append(out.Errors, envProblemDetails(report)...)
		out.Next = []string{"set missing required values in .env or shell env", "carbide health env"}
	}
	if err := a.writeJSON(out); err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("environment contract has %d missing required value(s)", len(report.missingRequired))
	}
	return nil
}

func (a app) commandHealthRuntime() error {
	results := a.projectHealthResults()
	if healthFailures(results) > 0 {
		results = append(results, healthResult{"runtime", "skip", "fix app laws first"})
		return a.renderHealthResults("runtime contract", results)
	}

	runtimeResults := a.runtimeHealthResults()
	results = append(results, runtimeResults...)
	return a.renderHealthResults("runtime contract", results)
}

func (a app) commandHealthRuntimeJSON() error {
	results := a.projectHealthResults()
	if healthFailures(results) > 0 {
		results = append(results, healthResult{"runtime", "skip", "fix app laws first"})
		return a.renderHealthResultsJSON("health runtime", results, []string{"fix app laws first", "carbide health runtime"})
	}

	results = append(results, a.runtimeHealthResults()...)
	return a.renderHealthResultsJSON("health runtime", results, nil)
}

func (a app) commandHealthFramework() error {
	results := a.frameworkHealthResults()
	return a.renderHealthResults("framework regressions", results)
}

func (a app) commandHealthFrameworkJSON() error {
	results := a.frameworkHealthResults()
	return a.renderHealthResultsJSON("health framework", results, nil)
}

func (a app) renderHealthResults(subtitle string, results []healthResult) error {
	rows := make([]tableRow, 0, len(results))
	for _, result := range results {
		rows = append(rows, tableRow{result.check, result.status, result.detail})
	}

	r := newRenderer(a.stdout)
	r.Title("Carbide health", subtitle)
	r.Table([]string{"check", "status", "detail"}, rows)

	failures := healthFailures(results)
	if failures > 0 {
		return fmt.Errorf("law compliance has %d failing check(s)", failures)
	}
	return nil
}

func (a app) renderHealthResultsJSON(command string, results []healthResult, next []string) error {
	failures := healthFailures(results)
	report := a.baseJSONReport(command)
	report.OK = failures == 0
	report.Checks = healthJSONChecks(results)
	report.Errors = healthFailureDetails(results)
	report.Next = next
	if failures > 0 && len(report.Next) == 0 {
		report.Next = []string{"fix failing checks", commandWithoutJSON(command)}
	}
	if err := a.writeJSON(report); err != nil {
		return err
	}
	if failures > 0 {
		return fmt.Errorf("law compliance has %d failing check(s)", failures)
	}
	return nil
}

func healthJSONChecks(results []healthResult) []healthJSONCheck {
	checks := make([]healthJSONCheck, 0, len(results))
	for _, result := range results {
		checks = append(checks, healthJSONCheck{
			Check:  result.check,
			Status: result.status,
			Detail: result.detail,
		})
	}
	return checks
}

func healthFailureDetails(results []healthResult) []string {
	var errors []string
	for _, result := range results {
		if result.status == "fail" {
			errors = append(errors, result.check+": "+result.detail)
		}
	}
	return errors
}

func commandWithoutJSON(command string) string {
	return "carbide " + command
}

func healthFailures(results []healthResult) int {
	count := 0
	for _, result := range results {
		if result.status == "fail" {
			count++
		}
	}
	return count
}

func healthOK(check string, detail string) healthResult {
	return healthResult{check: check, status: "ok", detail: detail}
}

func healthFail(check string, detail string) healthResult {
	return healthResult{check: check, status: "fail", detail: detail}
}

func healthWarn(check string, detail string) healthResult {
	return healthResult{check: check, status: "warn", detail: detail}
}

func healthSkip(check string, detail string) healthResult {
	return healthResult{check: check, status: "skip", detail: detail}
}
