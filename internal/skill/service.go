package skill

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	embeddedskill "github.com/nazarkhatsko/mink/internal/skill/embedded"
)

var ErrNotInstalled = errors.New("not installed")

type Status int

const (
	StatusNotInstalled Status = iota
	StatusUpToDate
	StatusOutdated
)

type Service struct {
	target   string
	platform string
}

// NewFromFlags constructs a Service by resolving the target directory and platform from CLI flags.
func NewFromFlags(forClaude, forCodex, global bool) (*Service, error) {
	if !forClaude && !forCodex {
		return nil, fmt.Errorf("required flag: --claude or --codex")
	}
	if forClaude && forCodex {
		return nil, fmt.Errorf("flags --claude and --codex are mutually exclusive")
	}
	target, platform, err := resolveTarget(forClaude, forCodex, global)
	if err != nil {
		return nil, err
	}
	return &Service{target: target, platform: platform}, nil
}

func resolveTarget(forClaude, forCodex, global bool) (target, platform string, err error) {
	var subdir string
	switch {
	case forClaude:
		platform = "claude"
		subdir = filepath.Join(".claude", "skills")
	case forCodex:
		platform = "codex"
		subdir = filepath.Join(".codex", "skills")
	}
	if global {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", "", fmt.Errorf("get home directory: %w", err)
		}
		return filepath.Join(home, subdir), platform, nil
	}
	if _, err := os.Stat(subdir); os.IsNotExist(err) {
		return "", "", fmt.Errorf("%s not found in current directory; use --global to install globally", subdir)
	}
	return subdir, platform, nil
}

func (s *Service) Names() ([]string, error) {
	return embeddedskill.List(s.platform)
}

func (s *Service) Status(name string) (Status, error) {
	embedded, err := embeddedskill.GetVersion(s.platform, name)
	if err != nil {
		return 0, err
	}
	data, err := os.ReadFile(filepath.Join(s.target, name+".version"))
	if os.IsNotExist(err) {
		return StatusNotInstalled, nil
	}
	if err != nil {
		return 0, err
	}
	installed := strings.TrimSpace(string(data))
	if embedded == installed {
		return StatusUpToDate, nil
	}
	return StatusOutdated, nil
}

func (s *Service) Install(name string) (bool, error) {
	st, err := s.Status(name)
	if err != nil {
		return false, err
	}
	if st == StatusUpToDate {
		return false, nil
	}
	content, err := embeddedskill.Get(s.platform, name)
	if err != nil {
		return false, fmt.Errorf("skill %q not found: %w", name, err)
	}
	version, err := embeddedskill.GetVersion(s.platform, name)
	if err != nil {
		return false, fmt.Errorf("skill %q version not found: %w", name, err)
	}
	if err := os.MkdirAll(s.target, 0755); err != nil {
		return false, fmt.Errorf("create directory: %w", err)
	}
	if err := os.WriteFile(filepath.Join(s.target, name+".md"), content, 0644); err != nil {
		return false, fmt.Errorf("write skill: %w", err)
	}
	if err := os.WriteFile(filepath.Join(s.target, name+".version"), []byte(version+"\n"), 0644); err != nil {
		return false, fmt.Errorf("write skill version: %w", err)
	}
	return true, nil
}

func (s *Service) Uninstall(name string) error {
	mdPath := filepath.Join(s.target, name+".md")
	if err := os.Remove(mdPath); err != nil {
		if os.IsNotExist(err) {
			return ErrNotInstalled
		}
		return fmt.Errorf("remove skill: %w", err)
	}
	_ = os.Remove(filepath.Join(s.target, name+".version"))
	return nil
}

func (s *Service) Upgrade(name string) (bool, error) {
	st, err := s.Status(name)
	if err != nil {
		return false, err
	}
	switch st {
	case StatusNotInstalled:
		return false, ErrNotInstalled
	case StatusUpToDate:
		return false, nil
	default:
		_, err = s.Install(name)
		return true, err
	}
}
