package config

import (
	_ "embed"
)

//go:embed resources/workspace.yaml
var defaultWorkspaceConfig string
