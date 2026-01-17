package platform

import (
	"os"
	"path/filepath"
)

func DefaultConfigPath(appName string) string {
	dir, err := os.UserConfigDir()
	if err != nil {
		home, _ := os.UserHomeDir()
		if home == "" {
			return appName + ".toml"
		}
		return filepath.Join(home, ".config", appName, "config.toml")
	}
	return filepath.Join(dir, appName, "config.toml")
}

func DefaultLogPath(appName string) string {
	dir, err := os.UserCacheDir()
	if err != nil {
		home, _ := os.UserHomeDir()
		if home == "" {
			return appName + ".log"
		}
		return filepath.Join(home, ".cache", appName, appName+".log")
	}
	return filepath.Join(dir, appName, appName+".log")
}

func DefaultStateDir(appName string) string {
	dir, err := os.UserCacheDir()
	if err != nil {
		home, _ := os.UserHomeDir()
		if home == "" {
			return filepath.Join(".", "state")
		}
		return filepath.Join(home, ".cache", appName, "state")
	}
	return filepath.Join(dir, appName, "state")
}
