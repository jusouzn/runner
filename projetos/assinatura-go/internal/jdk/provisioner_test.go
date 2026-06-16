package jdk

import (
	"runtime"
	"strings"
	"testing"
)

func TestJreURL_PlataformaAtual(t *testing.T) {
	url, isZip := jreURL()

	if !strings.HasPrefix(url, "https://api.adoptium.net/v3/binary/latest/21/ga/") {
		t.Errorf("jreURL() = %q, esperado endpoint GA do Java 21 do Adoptium", url)
	}
	if !strings.Contains(url, "/jre/hotspot/") {
		t.Errorf("jreURL() = %q, esperado pacote JRE/HotSpot", url)
	}

	// Windows usa .zip; demais usam .tar.gz.
	wantZip := runtime.GOOS == "windows"
	if isZip != wantZip {
		t.Errorf("jreURL() isZip = %v para GOOS=%s, want %v", isZip, runtime.GOOS, wantZip)
	}

	var wantOS string
	switch runtime.GOOS {
	case "windows":
		wantOS = "/windows/"
	case "darwin":
		wantOS = "/mac/"
	default:
		wantOS = "/linux/"
	}
	if !strings.Contains(url, wantOS) {
		t.Errorf("jreURL() = %q, esperado conter %q para GOOS=%s", url, wantOS, runtime.GOOS)
	}
}

func TestSafeJoin_RejeitaPathTraversal(t *testing.T) {
	base := t.TempDir()

	maliciosos := []string{
		"../escapou.txt",
		"..\\escapou.txt",
		"../../etc/passwd",
		"sub/../../fora.txt",
	}
	for _, rel := range maliciosos {
		if _, err := safeJoin(base, rel); err == nil {
			t.Errorf("safeJoin(%q) não retornou erro — zip-slip não bloqueado", rel)
		}
	}
}

func TestSafeJoin_AceitaCaminhoInterno(t *testing.T) {
	base := t.TempDir()

	validos := []string{
		"bin/java",
		"lib/modules",
		"release",
	}
	for _, rel := range validos {
		got, err := safeJoin(base, rel)
		if err != nil {
			t.Errorf("safeJoin(%q) retornou erro inesperado: %v", rel, err)
			continue
		}
		if !strings.HasPrefix(got, base) {
			t.Errorf("safeJoin(%q) = %q, esperado dentro de %q", rel, got, base)
		}
	}
}
