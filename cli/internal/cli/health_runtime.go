package cli

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"time"
)

func (a app) runtimeHealthResults() []healthResult {
	if !isFile(projectConfigPath) {
		return []healthResult{healthFail("runtime", "run this inside a Carbide project")}
	}
	profile := projectProfile()

	compose, err := findCompose()
	if err != nil {
		return []healthResult{healthFail("runtime", err.Error())}
	}

	env := setEnv(os.Environ(), "COMPOSE_MENU", "false")
	env = composeEnv(env)
	alreadyRunning := composeHasRunningServices(compose, env)
	port := 0
	if alreadyRunning {
		port = publishedWebPort(compose, env)
		if port == 0 {
			port = runtimePortFromEnv()
		}
	} else {
		selected, err := chooseDevPort(os.Getenv("CARBIDE_HTTP_PORT"))
		if err != nil {
			return []healthResult{healthFail("runtime", err.Error())}
		}
		port = selected
		env = setEnv(env, "CARBIDE_HTTP_PORT", strconv.Itoa(port))
	}

	results := []healthResult{}
	if _, err := runComposeCaptured(compose, env, "config"); err != nil {
		return append(results, healthFail("compose config", err.Error()))
	}
	results = append(results, healthOK("compose config", "valid"))

	startedByHealth := !alreadyRunning
	if startedByHealth {
		if err := composeUpDetached(compose, env); err != nil {
			results = append(results, healthFail("stack start", err.Error()))
			return results
		}
		results = append(results, healthOK("stack start", fmt.Sprintf("localhost:%d", port)))
	} else {
		results = append(results, healthOK("stack start", "already running"))
	}

	cleanupNeeded := startedByHealth
	if cleanupNeeded {
		defer func() {
			if cleanupNeeded {
				_ = composeDown(compose, env)
			}
		}()
	}

	client := &http.Client{Timeout: 10 * time.Second}
	if err := waitForHTTP(client, fmt.Sprintf("http://localhost:%d/health", port), 60*time.Second); err != nil {
		results = append(results, healthFail("health", err.Error()))
		return results
	}
	results = append(results, healthOK("health", "/health"))

	baseURL := fmt.Sprintf("http://localhost:%d", port)
	if profile == "docs" {
		if err := httpGetContains(client, baseURL+"/api/version", `"name":"Carbide Docs"`); err != nil {
			results = append(results, healthFail("version api", err.Error()))
			return results
		}
		results = append(results, healthOK("version api", "/api/version"))

		if startedByHealth {
			if err := composeDown(compose, env); err != nil {
				results = append(results, healthWarn("cleanup", err.Error()))
			} else {
				results = append(results, healthOK("cleanup", "stopped health stack"))
			}
			cleanupNeeded = false
		}
		return results
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		results = append(results, healthFail("auth flow", err.Error()))
		return results
	}
	client.Jar = jar
	if err := httpGetContains(client, baseURL+"/api/me", `"authenticated":false`); err != nil {
		results = append(results, healthFail("anonymous", err.Error()))
		return results
	}
	results = append(results, healthOK("anonymous", "/api/me"))

	email := fmt.Sprintf("health-%d@carbide.local", time.Now().UnixNano())
	if err := httpPostFormContains(client, baseURL+"/api/register", url.Values{"email": {email}, "password": {"password"}}, `"ok":true`); err != nil {
		results = append(results, healthFail("register", err.Error()))
		return results
	}
	results = append(results, healthOK("register", "first-user flow"))

	if err := httpGetContains(client, baseURL+"/api/dashboard", email); err != nil {
		results = append(results, healthFail("dashboard api", err.Error()))
		return results
	}
	if err := httpGetContains(client, baseURL+"/dashboard", `<div id="root"></div>`); err != nil {
		results = append(results, healthFail("dashboard web", err.Error()))
		return results
	}
	results = append(results, healthOK("dashboard", "api and web shell"))

	if err := httpPostFormContains(client, baseURL+"/api/logout", nil, `"ok":true`); err != nil {
		results = append(results, healthFail("logout", err.Error()))
		return results
	}
	if err := httpGetContains(client, baseURL+"/api/me", `"authenticated":false`); err != nil {
		results = append(results, healthFail("logout", err.Error()))
		return results
	}
	results = append(results, healthOK("logout", "session cleared"))

	if startedByHealth {
		if err := composeDown(compose, env); err != nil {
			results = append(results, healthWarn("cleanup", err.Error()))
		} else {
			results = append(results, healthOK("cleanup", "stopped health stack"))
		}
		cleanupNeeded = false
	}
	return results
}

func (a app) frameworkHealthResults() []healthResult {
	root, err := resolveFrameworkRoot(a.home)
	if err != nil {
		return []healthResult{healthFail("framework", err.Error())}
	}
	env, cleanup, err := frameworkHealthCommandEnv(root)
	if err != nil {
		return []healthResult{healthFail("framework", err.Error())}
	}
	defer cleanup()

	type frameworkCheck struct {
		name string
		run  func() error
	}
	checks := []frameworkCheck{
		{
			name: "shell syntax",
			run: func() error {
				_, err := commandOutputEnv(
					root,
					env,
					"bash",
					"-n",
					"tests/contract/audit_versions.sh",
					"tests/contract/check_line_limits.sh",
					"tests/contract/check_repo_contract.sh",
					"tests/scaffold/cli_scaffold.sh",
					"tests/smoke/starter_docker_flow.sh",
					"tests/smoke/docs_for_agents_http.sh",
					"cli/bin/carbide",
					"cli/install.sh",
				)
				return err
			},
		},
		{name: "Go CLI tests", run: func() error { return runFrameworkGoTests(root) }},
		{name: "repo contract", run: func() error {
			_, err := commandOutputEnv(root, env, "bash", "tests/contract/check_repo_contract.sh")
			return err
		}},
		{name: "CLI scaffold", run: func() error {
			_, err := commandOutputEnv(root, env, "bash", "tests/scaffold/cli_scaffold.sh")
			return err
		}},
		{name: "Docker smoke", run: func() error {
			_, err := commandOutputEnv(root, env, "bash", "tests/smoke/starter_docker_flow.sh")
			return err
		}},
		{name: "agent guide HTTP smoke", run: func() error {
			_, err := commandOutputEnv(root, env, "bash", "tests/smoke/docs_for_agents_http.sh")
			return err
		}},
	}

	results := make([]healthResult, 0, len(checks))
	for _, check := range checks {
		if err := check.run(); err != nil {
			results = append(results, healthFail(check.name, firstLine(err.Error())))
			continue
		}
		results = append(results, healthOK(check.name, "passed"))
	}
	return results
}
