package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func SetCommit(value string) {
	if value != "" {
		commit = value
	}
}

func Main() {
	home, err := resolveHome()
	if err != nil {
		renderError(os.Stderr, err)
		os.Exit(1)
	}

	a := app{
		home:   home,
		stdout: os.Stdout,
		stderr: os.Stderr,
	}

	if err := a.run(os.Args[1:]); err != nil {
		renderError(os.Stderr, err)
		os.Exit(1)
	}
}

func (a app) run(args []string) error {
	if len(args) == 0 {
		a.printCommandList()
		return nil
	}

	switch args[0] {
	case "help", "-h", "--help":
		if len(args) != 1 {
			return errors.New("usage: carbide help")
		}
		a.printHelp()
		return nil
	case "version":
		if len(args) != 1 {
			return errors.New("usage: carbide version")
		}
		return a.commandVersion()
	case "upgrade":
		if len(args) != 1 {
			return errors.New("usage: carbide upgrade")
		}
		return a.commandUpgrade()
	case "new":
		if len(args) < 2 {
			return errors.New("usage: carbide new <project-name>")
		}
		return a.commandNew(strings.Join(args[1:], " "))
	case "init":
		if len(args) != 1 {
			return errors.New("usage: carbide init")
		}
		return a.commandInit()
	case "health":
		if len(args) == 1 {
			return a.commandHealth()
		}
		if len(args) == 2 && args[1] == "json" {
			return a.commandHealthJSON()
		}
		if len(args) == 2 && args[1] == "env" {
			return a.commandHealthEnv()
		}
		if len(args) == 3 && args[1] == "env" && args[2] == "json" {
			return a.commandHealthEnvJSON()
		}
		if len(args) == 2 && args[1] == "runtime" {
			return a.commandHealthRuntime()
		}
		if len(args) == 3 && args[1] == "runtime" && args[2] == "json" {
			return a.commandHealthRuntimeJSON()
		}
		if len(args) == 2 && args[1] == "framework" {
			return a.commandHealthFramework()
		}
		if len(args) == 3 && args[1] == "framework" && args[2] == "json" {
			return a.commandHealthFrameworkJSON()
		}
		return errors.New("usage: carbide health [json|env [json]|runtime [json]|framework [json]]")
	case "audit":
		return a.commandAuditFlow(args[1:])
	case "resolve":
		return a.commandResolveFlow(args[1:])
	case "fix":
		if len(args) != 1 {
			return errors.New("usage: carbide fix")
		}
		return a.commandFix()
	case "project":
		if len(args) == 2 && args[1] == "migrate" {
			return a.commandAuditFlow(nil)
		}
		return errors.New("usage: carbide audit")
	case "deploy":
		if len(args) == 2 {
			return a.commandDeploy(args[1])
		}
		return errors.New("usage: carbide deploy <target>")
	case "run":
		if len(args) == 2 && args[1] == "dev" {
			return a.commandRunDev()
		}
		return errors.New("usage: carbide run dev")
	case "status":
		if len(args) == 1 {
			return a.commandStatus()
		}
		if len(args) == 2 && args[1] == "json" {
			return a.commandStatusJSON()
		}
		return errors.New("usage: carbide status [json]")
	case "urls":
		if len(args) == 1 {
			return a.commandURLs(false)
		}
		if len(args) == 2 && args[1] == "json" {
			return a.commandURLs(true)
		}
		return errors.New("usage: carbide urls [json]")
	case "clean":
		if len(args) == 2 && args[1] == "dev" {
			return a.commandCleanDev()
		}
		return errors.New("usage: carbide clean dev")
	case "stop":
		if len(args) == 2 && args[1] == "dev" {
			return a.commandStopDev()
		}
		return errors.New("usage: carbide stop dev")
	case "follow":
		if len(args) >= 2 && args[1] == "logs" {
			return a.commandFollowLogs(args[2:])
		}
		return errors.New("usage: carbide follow logs [service <name>] [containing <text>]")
	case "logs":
		return a.commandLogs(args[1:])
	default:
		return fmt.Errorf("unknown command: %s", args[0])
	}
}

func (a app) printCommandList() {
	r := newRenderer(a.stdout)
	logo := carbideLogo()
	if r.interactive {
		r.AnimateLogo(logo)
	} else {
		r.Logo(logo)
	}
	text := fmt.Sprintf(commandListText, version)
	if r.styled {
		fmt.Fprint(a.stdout, r.paint("38;5;245", text))
		return
	}
	fmt.Fprint(a.stdout, text)
}

func (a app) printHelp() {
	r := newRenderer(a.stdout)
	r.CommandList([]helpCommandSection{
		{
			rows: []outputRow{
				{"audit", "audit laws and taste into .audit/report"},
				{"clean dev", "normalize local dev state"},
				{"deploy prod", "run checked-in deploy script"},
				{"fix", "implement the latest .audit plan"},
				{"health", "show law compliance"},
				{"health json", "show law compliance as JSON"},
				{"health env", "validate env contract"},
				{"health env json", "validate env contract as JSON"},
				{"health framework", "run framework regressions"},
				{"health framework json", "run framework regressions as JSON"},
				{"health runtime", "run Docker runtime checks"},
				{"health runtime json", "run Docker runtime checks as JSON"},
				{"help", "show this help"},
				{"init", "init current directory"},
				{"logs", "query saved logs"},
				{"new <project-name>", "create project directory"},
				{"resolve", "turn audit reports into a plan"},
				{"resolve fix", "plan and implement the latest audit"},
				{"status", "show containers and ports"},
				{"status json", "show containers and ports as JSON"},
				{"upgrade", "upgrade CLI from GitHub"},
				{"urls", "print local app and API URLs"},
				{"urls json", "print local app and API URLs as JSON"},
				{"version", "print installed version"},
			},
		},
		{
			name: "follow",
			rows: []outputRow{
				{"follow logs", "stream live logs"},
				{"follow logs service api", "stream one service"},
			},
		},
		{
			name: "logs",
			rows: []outputRow{
				{"logs containing \"/api/login\" json", "query logs as JSON"},
			},
		},
		{
			name: "run",
			rows: []outputRow{
				{"run dev", "start Docker dev stack"},
			},
		},
		{
			name: "stop",
			rows: []outputRow{
				{"stop dev", "stop dev containers"},
			},
		},
	})
}

func (a app) commandVersion() error {
	r := newRenderer(a.stdout)
	if commit != "" {
		r.Title("Carbide", "installed CLI")
		r.Rows(
			outputRow{"version", version},
			outputRow{"commit", commit},
		)
		return nil
	} else if head := gitShortHead(a.home); head != "" {
		r.Title("Carbide", "installed CLI")
		r.Rows(
			outputRow{"version", version},
			outputRow{"commit", head},
		)
		return nil
	}
	r.Title("Carbide", "installed CLI")
	r.Rows(outputRow{"version", version})
	return nil
}

func (a app) commandNew(name string) error {
	if err := ensureProjectName(name); err != nil {
		return err
	}

	slug := projectSlug(name)
	if slug == "" {
		slug = "carbide-app"
	}
	displayName := projectDisplayName(name)

	target, err := filepath.Abs(filepath.Join(".", slug))
	if err != nil {
		return err
	}
	if _, err := os.Stat(target); err == nil {
		return fmt.Errorf("%s already exists", slug)
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if err := a.copyScaffold(target, displayName, slug); err != nil {
		return err
	}

	newRenderer(a.stdout).Message(
		"Carbide",
		"project created",
		outputRow{"path", target},
		outputRow{"next", fmt.Sprintf("cd %s", slug)},
		outputRow{"", "carbide run dev"},
	)
	return nil
}

func (a app) commandInit() error {
	empty, err := isCurrentDirEmpty()
	if err != nil {
		return err
	}
	if !empty {
		return errors.New("carbide init requires an empty directory")
	}

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	name := filepath.Base(pwd)
	if err := ensureProjectName(name); err != nil {
		return err
	}

	slug := projectSlug(name)
	if slug == "" {
		slug = "carbide-app"
	}
	displayName := projectDisplayName(name)
	if err := a.copyScaffold(pwd, displayName, slug); err != nil {
		return err
	}

	newRenderer(a.stdout).Message(
		"Carbide",
		"project initialized",
		outputRow{"path", pwd},
		outputRow{"next", "carbide run dev"},
	)
	return nil
}
