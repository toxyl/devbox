package bashcompletion

import (
	"bytes"
	_ "embed"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"text/template"
)

type Completion struct {
	Variadic    bool
	Completions string
	Compopt     []string
	Shopt       []string
}

//go:embed resources/bashcompletion.gotemplate
var completionTemplate string

// generate generates a Bash completion script for the given appName and commands
func generate(appName string, devboxDir string, commands []string, commandData map[string][]Completion) ([]byte, error) {
	tmpl, err := template.New("completion").Parse(completionTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse completion template: %v", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, struct {
		AppName     string
		Commands    []string
		CommandData map[string][]Completion
		DevboxDir   string
	}{
		AppName:     appName,
		Commands:    commands,
		CommandData: commandData,
		DevboxDir:   devboxDir,
	}); err != nil {
		return nil, fmt.Errorf("failed to execute completion template: %v", err)
	}

	return buf.Bytes(), nil
}

// install installs the given Bash completion script
func install(script []byte, appName string) error {
	completionDir := filepath.Join("/", "etc", "bash_completion.d")
	scriptPath := filepath.Join(completionDir, appName)
	if err := ioutil.WriteFile(scriptPath, script, 0644); err != nil {
		return fmt.Errorf("failed to write completion script: %v", err)
	}
	return nil
}

// Make generates a Bash completion script for the given appName and commands and installs it
func Make(appName string, devboxDir string, commands []string, commandData map[string][]Completion) error {
	script, err := generate(appName, devboxDir, commands, commandData)
	if err != nil {
		return err
	}

	if err := install(script, appName); err != nil {
		return err
	}

	return nil
}
