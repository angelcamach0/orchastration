package config

import (
	"errors"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Logging LoggingConfig `toml:"logging"`
	Hash    HashConfig    `toml:"hash"`
	Jobs    map[string]JobConfig `toml:"jobs"`
	Tasks   map[string]TaskConfig `toml:"tasks"`
	Agents         map[string]AgentConfig         `toml:"agents"`
	Orchestrations map[string]OrchestrationConfig `toml:"orchestrations"`
}

type LoggingConfig struct {
	Level string `toml:"level"`
}

type HashConfig struct {
	Algorithm string `toml:"algorithm"`
}

type JobConfig struct {
	Description    string            `toml:"description"`
	Command        []string          `toml:"command"`
	WorkingDir     string            `toml:"working_dir"`
	TimeoutSeconds int               `toml:"timeout_seconds"`
	Env            map[string]string `toml:"env"`
}

type TaskConfig struct {
	Description string   `toml:"description"`
	Repo        string   `toml:"repo"`
	WorkingDir  string   `toml:"working_dir"`
	Command     []string `toml:"command"`
	Outputs     []string `toml:"outputs"`
	Documents   []string `toml:"documents"`
	Status      string   `toml:"status"`
}

type AgentConfig struct{}

type OrchestrationConfig struct {
	Agents      []string `toml:"agents"`
	Steps       [][]string `toml:"steps"`
	Description string   `toml:"description"`
}

func Default() Config {
	return Config{
		Logging: LoggingConfig{Level: "info"},
		Hash:    HashConfig{Algorithm: "sha256"},
	}
}

func Load(path string) (Config, error) {
	cfg := Default()

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}
		return cfg, err
	}

	if err := toml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}

	if cfg.Logging.Level == "" {
		cfg.Logging.Level = "info"
	}
	if cfg.Hash.Algorithm == "" {
		cfg.Hash.Algorithm = "sha256"
	}

	return cfg, nil
}
