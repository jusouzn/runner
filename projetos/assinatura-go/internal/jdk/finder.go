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
		"java não encontrado — instale o Java 21 ou execute `assinatura setup` para provisionamento automático",
	)
}

// FindJar returns the path to assinador.jar.
// It checks ~/.hubsaude/, the directory alongside the binary, and the current directory.
func FindJar() (string, error) {
	home, err := os.UserHomeDir()
	if err == nil {
		managed := filepath.Join(home, ".hubsaude", "assinador.jar")
		if _, err := os.Stat(managed); err == nil {
			return managed, nil
		}
	}
	if exe, err := os.Executable(); err == nil {
		next := filepath.Join(filepath.Dir(exe), "assinador.jar")
		if _, err := os.Stat(next); err == nil {
			return next, nil
		}
	}
	if _, err := os.Stat("assinador.jar"); err == nil {
		return "assinador.jar", nil
	}
	return "", errors.New(
		"assinador.jar não encontrado — coloque-o em ~/.hubsaude/ ou no mesmo diretório do executável",
	)
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
