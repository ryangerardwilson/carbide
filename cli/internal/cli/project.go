package cli

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

func resolveHome() (string, error) {
	if home := os.Getenv("CARBIDE_HOME"); home != "" {
		return filepath.Abs(home)
	}

	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		return "", err
	}

	dir := filepath.Dir(exe)
	switch filepath.Base(dir) {
	case "bin":
		parent := filepath.Dir(dir)
		switch filepath.Base(parent) {
		case ".cli", "cli":
			return filepath.Dir(parent), nil
		default:
			return parent, nil
		}
	case ".bin":
		return filepath.Dir(dir), nil
	default:
		return filepath.Dir(dir), nil
	}
}

func shouldStyleOutput(w io.Writer) bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	return isTerminalOutput(w)
}

func shouldLaunchAuditCodex(w io.Writer) bool {
	return isTerminalOutput(w) && isTerminalFile(os.Stdin)
}

func isTerminalOutput(w io.Writer) bool {
	file, ok := w.(*os.File)
	if !ok {
		return false
	}
	return isTerminalFile(file)
}

func isTerminalFile(file *os.File) bool {
	if file == nil {
		return false
	}
	info, err := file.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

type terminalWindowSize struct {
	rows    uint16
	columns uint16
	xpixels uint16
	ypixels uint16
}

func terminalColumns(w io.Writer) int {
	file, ok := w.(*os.File)
	if !ok {
		return 0
	}
	var size terminalWindowSize
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		file.Fd(),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&size)),
	)
	if errno != 0 || size.columns == 0 {
		return 0
	}
	return int(size.columns)
}

func terminalColumnsFromEnv() int {
	columns, err := strconv.Atoi(os.Getenv("COLUMNS"))
	if err != nil || columns <= 0 {
		return 0
	}
	return columns
}

func ensureProjectName(name string) error {
	if name == "" || strings.HasPrefix(name, ".") || strings.ContainsAny(name, `/\`) {
		return errors.New("project name must be a simple directory name")
	}
	matched, err := regexp.MatchString(`^[A-Za-z0-9._ -]+$`, name)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("project name may contain only letters, numbers, spaces, dots, underscores, and dashes")
	}
	return nil
}

func projectDisplayName(input string) string {
	words := projectWords(input)
	if len(words) == 0 {
		return "My Carbide App"
	}

	for i, word := range words {
		words[i] = titleWord(word)
	}
	return strings.Join(words, " ")
}

func projectWords(input string) []string {
	parts := strings.FieldsFunc(input, func(r rune) bool {
		return r == ' ' || r == '-' || r == '_' || r == '.'
	})
	words := make([]string, 0, len(parts))
	for _, part := range parts {
		word := strings.TrimSpace(part)
		if word != "" {
			words = append(words, word)
		}
	}
	return words
}

func titleWord(word string) string {
	lower := strings.ToLower(word)
	runes := []rune(lower)
	if len(runes) == 0 {
		return ""
	}
	runes[0] = []rune(strings.ToUpper(string(runes[0])))[0]
	return string(runes)
}

func projectSlug(input string) string {
	var b strings.Builder
	lastDash := false
	for _, r := range strings.ToLower(input) {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			b.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash {
			b.WriteByte('-')
			lastDash = true
		}
	}
	return strings.Trim(b.String(), "-")
}

type projectMetadata struct {
	name    string
	slug    string
	profile string
}

func projectProfile() string {
	data, err := os.ReadFile(projectConfigPath)
	if err != nil {
		return "app"
	}

	section := ""
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(stripTomlComment(scanner.Text()))
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section = strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(line, "["), "]"))
			continue
		}
		if section != "" {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok || strings.TrimSpace(key) != "profile" {
			continue
		}
		if profile := strings.TrimSpace(parseTomlString(value)); profile != "" {
			return profile
		}
	}
	return "app"
}

func readProjectMetadata(path string) (projectMetadata, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return projectMetadata{}, err
	}

	var metadata projectMetadata
	section := ""
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(stripTomlComment(scanner.Text()))
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section = strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(line, "["), "]"))
			continue
		}
		if section != "" {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		switch strings.TrimSpace(key) {
		case "name":
			metadata.name = parseTomlString(value)
		case "slug":
			metadata.slug = parseTomlString(value)
		case "profile":
			metadata.profile = parseTomlString(value)
		}
	}
	if err := scanner.Err(); err != nil {
		return projectMetadata{}, err
	}
	return metadata, nil
}

func composeProjectName() string {
	metadata, err := readProjectMetadata(projectConfigPath)
	if err == nil {
		if slug := normalizeComposeProjectName(metadata.slug); slug != "" {
			return slug
		}
	}

	pwd, err := os.Getwd()
	if err == nil {
		if slug := normalizeComposeProjectName(filepath.Base(pwd)); slug != "" {
			return slug
		}
	}
	return "carbide-app"
}

func normalizeComposeProjectName(value string) string {
	value = strings.TrimSpace(value)
	if value == "" || isTemplatePlaceholder(value) {
		return ""
	}
	return projectSlug(value)
}

func isTemplatePlaceholder(value string) bool {
	return strings.Contains(value, "__PROJECT_")
}

func displayCommit(home string) string {
	if commit != "" {
		return commit
	}
	if home != "" {
		return gitShortHead(home)
	}
	return "unknown"
}

func (a app) copyScaffold(target string, name string, slug string) error {
	source := filepath.Join(a.home, "scaffold")
	if !isDir(source) {
		return fmt.Errorf("missing scaffold source: %s", source)
	}
	return copyScaffoldPart(source, target, name, slug)
}

func copyScaffoldPart(source string, target string, name string, slug string) error {
	return filepath.WalkDir(source, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		rel, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return os.MkdirAll(target, 0755)
		}
		if skipGeneratedScaffoldPath(rel, entry.IsDir()) {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		dest := filepath.Join(target, rel)
		info, err := entry.Info()
		if err != nil {
			return err
		}

		if entry.IsDir() {
			return os.MkdirAll(dest, info.Mode().Perm())
		}
		if entry.Type()&os.ModeSymlink != 0 {
			link, err := os.Readlink(path)
			if err != nil {
				return err
			}
			return os.Symlink(link, dest)
		}
		if !entry.Type().IsRegular() {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		content = bytes.ReplaceAll(content, []byte("__PROJECT_NAME__"), []byte(name))
		content = bytes.ReplaceAll(content, []byte("__PROJECT_SLUG__"), []byte(slug))
		return os.WriteFile(dest, content, info.Mode().Perm())
	})
}

func skipGeneratedScaffoldPath(rel string, isDir bool) bool {
	rel = filepath.ToSlash(rel)
	generatedDirs := map[string]bool{
		".carbide":         true,
		"web/node_modules": true,
		"web/public":       true,
	}
	if isDir && generatedDirs[rel] {
		return true
	}
	generatedFiles := map[string]bool{
		".env":                 true,
		"web/src/tailwind.css": true,
	}
	return generatedFiles[rel]
}

func isCurrentDirEmpty() (bool, error) {
	entries, err := os.ReadDir(".")
	if err != nil {
		return false, err
	}
	return len(entries) == 0, nil
}
