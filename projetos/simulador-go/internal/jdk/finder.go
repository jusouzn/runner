package jdk

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// FindJava returns the path to the java binary.
// It checks ~/.hubsaude/jdk first, then falls back to PATH.
func FindJava() (string, error) {
	home, err := os.UserHomeDir()
	if err == nil {
		bin := "java"
		if runtime.GOOS == "windows" {
			bin = "java.exe"
		}
		managed := filepath.Join(home, ".hubsaude", "jdk", "bin", bin)
		if isFile(managed) {
			return managed, nil
		}
	}
	if path, err := exec.LookPath("java"); err == nil {
		return path, nil
	}
	return "", errors.New(
		"java não encontrado — execute `simulador obter` para provisionar o JRE automaticamente",
	)
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
