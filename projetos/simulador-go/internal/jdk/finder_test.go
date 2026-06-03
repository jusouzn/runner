package jdk

import (
	"testing"
)

func TestParseVersionOutput(t *testing.T) {
	tests := []struct {
		name    string
		output  string
		want    int
		wantErr bool
	}{
		{
			name:   "openjdk 21",
			output: `openjdk version "21.0.3" 2024-04-16 LTS`,
			want:   21,
		},
		{
			name:   "java 21",
			output: `java version "21.0.2" 2024-01-16`,
			want:   21,
		},
		{
			name:   "openjdk 17",
			output: `openjdk version "17.0.11" 2024-04-16`,
			want:   17,
		},
		{
			name:   "java 11",
			output: `java version "11.0.23" 2024-04-16`,
			want:   11,
		},
		{
			name:   "java 8 (esquema 1.x)",
			output: `java version "1.8.0_412" 2024-04-16`,
			want:   8,
		},
		{
			name:   "openjdk 22",
			output: `openjdk version "22" 2024-03-19`,
			want:   22,
		},
		{
			name:    "saída inválida",
			output:  `not a java version string`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseVersionOutput(tt.output)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseVersionOutput() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("parseVersionOutput() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestMinJavaVersion(t *testing.T) {
	if MinJavaVersion != 21 {
		t.Errorf("MinJavaVersion = %d, want 21", MinJavaVersion)
	}
}

func TestIsCachedJDKValid_NotExists(t *testing.T) {
	valid := isCachedJDKValid("/tmp/nao-existe-jdk-" + t.Name())
	if valid {
		t.Error("isCachedJDKValid() = true para diretório inexistente, want false")
	}
}
