package jdk

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Provision downloads and extracts JRE 21 (Eclipse Temurin) to ~/.hubsaude/jdk.
// If a valid JRE >= 21 is already provisioned there, it skips the download.
func Provision() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("não foi possível determinar o diretório home: %w", err)
	}
	destDir := filepath.Join(home, ".hubsaude", "jdk")

	// Cache check: skip download if already provisioned and valid.
	if isCachedJDKValid(destDir) {
		fmt.Fprintln(os.Stderr, "JRE 21 já provisionado em ~/.hubsaude/jdk — nenhum download necessário.")
		return nil
	}

	url, isZip := jreURL()
	fmt.Fprintf(os.Stderr, "Baixando JRE 21 (Eclipse Temurin)...\n")

	tmp, err := os.CreateTemp("", "jre-download-*")
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo temporário: %w", err)
	}
	defer os.Remove(tmp.Name())

	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return fmt.Errorf("erro ao baixar JRE: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download falhou com status HTTP %d", resp.StatusCode)
	}

	if _, err := io.Copy(tmp, resp.Body); err != nil {
		tmp.Close()
		return fmt.Errorf("erro ao salvar arquivo: %w", err)
	}
	tmp.Close()

	fmt.Fprintln(os.Stderr, "Extraindo JRE...")
	os.RemoveAll(destDir)

	if isZip {
		return extractZip(tmp.Name(), destDir)
	}
	return extractTarGz(tmp.Name(), destDir)
}

// isCachedJDKValid returns true if ~/.hubsaude/jdk/bin/java exists and reports version >= MinJavaVersion.
func isCachedJDKValid(destDir string) bool {
	bin := "java"
	if runtime.GOOS == "windows" {
		bin = "java.exe"
	}
	javaBin := filepath.Join(destDir, "bin", bin)
	if !isFile(javaBin) {
		return false
	}
	v, err := javaVersion(javaBin)
	return err == nil && v >= MinJavaVersion
}

func jreURL() (url string, isZip bool) {
	arch := "x64"
	switch runtime.GOOS {
	case "windows":
		return fmt.Sprintf(
			"https://api.adoptium.net/v3/binary/latest/21/ga/windows/%s/jre/hotspot/normal/eclipse",
			arch), true
	case "darwin":
		return fmt.Sprintf(
			"https://api.adoptium.net/v3/binary/latest/21/ga/mac/%s/jre/hotspot/normal/eclipse",
			arch), false
	default:
		return fmt.Sprintf(
			"https://api.adoptium.net/v3/binary/latest/21/ga/linux/%s/jre/hotspot/normal/eclipse",
			arch), false
	}
}

func extractTarGz(src, dest string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Strip top-level directory (e.g. "jdk-21.0.3+9/bin/java" → "bin/java").
		parts := strings.SplitN(hdr.Name, "/", 2)
		if len(parts) < 2 || parts[1] == "" {
			continue
		}

		target, err := safeJoin(dest, parts[1])
		if err != nil {
			continue // skip unsafe paths
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			os.MkdirAll(target, os.FileMode(hdr.Mode))
		case tar.TypeReg:
			os.MkdirAll(filepath.Dir(target), 0755)
			out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}
			_, copyErr := io.Copy(out, tr)
			out.Close()
			if copyErr != nil {
				return copyErr
			}
		case tar.TypeSymlink:
			os.Remove(target)
			os.Symlink(hdr.Linkname, target) //nolint:errcheck
		}
	}
	return nil
}

func extractZip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		// Strip top-level directory.
		parts := strings.SplitN(f.Name, "/", 2)
		if len(parts) < 2 || parts[1] == "" {
			continue
		}

		target, err := safeJoin(dest, parts[1])
		if err != nil {
			continue
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(target, f.Mode())
			continue
		}

		os.MkdirAll(filepath.Dir(target), 0755)
		rc, err := f.Open()
		if err != nil {
			return err
		}

		out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
		if err != nil {
			rc.Close()
			return err
		}
		_, copyErr := io.Copy(out, rc)
		out.Close()
		rc.Close()
		if copyErr != nil {
			return copyErr
		}
	}
	return nil
}

// safeJoin joins base and rel and verifies the result stays inside base (zip-slip guard).
func safeJoin(base, rel string) (string, error) {
	target := filepath.Join(base, filepath.FromSlash(rel))
	absBase, err := filepath.Abs(base)
	if err != nil {
		return "", err
	}
	absTarget, err := filepath.Abs(target)
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(absTarget, absBase+string(os.PathSeparator)) && absTarget != absBase {
		return "", fmt.Errorf("path traversal: %s", rel)
	}
	return target, nil
}
