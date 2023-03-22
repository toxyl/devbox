package config

import (
	"os"
	"path/filepath"

	"github.com/toxyl/devbox/devip"
	"gopkg.in/yaml.v2"
)

type WorkspaceDevbox struct {
	Name   string `mapstructure:"name"`
	Delay  int64  `mapstructure:"delay"`  // defines how long to wait before starting this devbox
	Image  string `mapstructure:"image"`  // usually this would the tarball eith the latest image
	Config Config `mapstructure:"config"` // a Config object used to configure the workspace
}

type Workspace struct {
	Path     string            `mapstructure:"path"`
	IPs      []string          `mapstructure:"ips"`
	Devboxes []WorkspaceDevbox `mapstructure:"devboxes"`
}

func (w *Workspace) AddIP(ip string) {
	for _, wip := range w.IPs {
		if wip == ip {
			return
		}
	}
	devip.Add(ip)
	w.IPs = append(w.IPs, ip)
}

func (w *Workspace) RemoveIP(ip string) {
	for i, wip := range w.IPs {
		if wip == ip {
			devip.Remove(ip)
			w.IPs = append(w.IPs[:i], w.IPs[i+1:]...)
			return
		}
	}
}

func (w *Workspace) Add(name, image string, startDelay int64, c Config) {
	wd := WorkspaceDevbox{
		Config: c,
		Image:  image,
		Name:   name,
		Delay:  startDelay,
	}
	for i, d := range w.Devboxes {
		if d.Name == wd.Name {
			w.Devboxes[i] = wd
			return
		}
	}
	w.Devboxes = append(w.Devboxes, wd)
}

func (w *Workspace) Remove(name string) {
	for i, d := range w.Devboxes {
		if d.Name == name {
			w.Devboxes = append(w.Devboxes[:i], w.Devboxes[i+1:]...)
			return
		}
	}
}

func (w *Workspace) Save(file string) error {
	yamlConfig, err := yaml.Marshal(w)
	if err != nil {
		return err
	}

	err = os.WriteFile(file, yamlConfig, 0644)
	if err != nil {
		return err
	}

	return nil
}

func NewWorkspace(path string) *Workspace {
	w := &Workspace{
		Path: path,
		IPs:  []string{},
	}
	return w
}

func OpenWorkspace(file string) (*Workspace, error) {
	file, err := filepath.Abs(file)
	if err != nil {
		return nil, err
	}
	dir := filepath.Dir(file)
	if !fileExists(dir) {
		_ = os.MkdirAll(dir, 0755)
		_ = os.WriteFile(file, []byte(defaultWorkspaceConfig), 0644)
	}
	w := NewWorkspace(dir)
	err = parseWithViper(file, w)

	return w, err
}
