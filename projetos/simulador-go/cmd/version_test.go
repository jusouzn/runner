package cmd

import (
	"strings"
	"testing"
)

func TestVersionString_ComCommit(t *testing.T) {
	defer restaurar(Version, Commit, Date)
	Version, Commit, Date = "1.2.3", "abc1234", "2026-06-17"

	got := versionString()
	for _, sub := range []string{"v1.2.3", "abc1234", "2026-06-17"} {
		if !strings.Contains(got, sub) {
			t.Errorf("versionString()=%q não contém %q", got, sub)
		}
	}
}

func TestVersionString_SemCommit(t *testing.T) {
	defer restaurar(Version, Commit, Date)
	Version, Commit, Date = "0.1.0-dev", "", ""

	got := versionString()
	if !strings.Contains(got, "v0.1.0-dev") {
		t.Errorf("versionString()=%q deveria conter a versão", got)
	}
	if strings.Contains(got, "(") {
		t.Errorf("sem commit não deveria haver parênteses: %q", got)
	}
}

func restaurar(v, c, d string) {
	Version, Commit, Date = v, c, d
}
