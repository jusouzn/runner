package jdk

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
)

const MinJavaVersion = 21

var javaVersionRe = regexp.MustCompile(`(?:java|openjdk) version "(\d+)[\.\-]?`)

// FindJava returns the path to a java binary with version >= MinJavaVersion.
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
			if v, err := javaVersion(managed); err == nil && v >= MinJavaVersion {
				return managed, nil
			}
		}
	}
	if path, err := exec.LookPath("java"); err == nil {
		if v, err := javaVersion(path); err == nil && v >= MinJavaVersion {
			return path, nil
		}
		return "", fmt.Errorf(
			"java encontrado em %s, mas versão incompatível (mínimo: %d) — execute `assinatura setup`",
			path, MinJavaVersion,
		)
	}
	return "", errors.New(
		"java não encontrado — instale o Java 21 ou execute `assinatura setup` para provisionamento automático",
	)
}

// ParseJavaVersion runs java -version on the given binary and returns the major version number.
func ParseJavaVersion(javaPath string) (int, error) {
	return javaVersion(javaPath)
}

func javaVersion(javaPath string) (int, error) {
	out, err := exec.Command(javaPath, "-version").CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("erro ao executar java -version: %w", err)
	}
	return parseVersionOutput(string(out))
}

func parseVersionOutput(output string) (int, error) {
	m := javaVersionRe.FindStringSubmatch(output)
	if m == nil {
		return 0, fmt.Errorf("não foi possível identificar versão do Java na saída: %q", output)
	}
	major, err := strconv.Atoi(m[1])
	if err != nil {
		return 0, fmt.Errorf("versão Java inválida: %q", m[1])
	}
	// Versões antigas usam o esquema "1.8" → major é 1, segunda parte é 8.
	if major == 1 {
		re2 := regexp.MustCompile(`"1\.(\d+)`)
		if m2 := re2.FindStringSubmatch(output); m2 != nil {
			major, _ = strconv.Atoi(m2[1])
		}
	}
	return major, nil
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
