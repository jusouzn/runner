package release

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// releaseURL é a URL fixa do release.json no repositório da disciplina.
// Atualizar se o repositório mudar.
const releaseURL = "https://raw.githubusercontent.com/kyriosdata/simulador/main/release.json"

// Info representa o conteúdo do release.json.
type Info struct {
	Jar struct {
		URL     string `json:"url"`
		Version string `json:"version"`
	} `json:"jar"`
	JRE struct {
		WindowsX64 string `json:"windows_x64"`
		LinuxX64   string `json:"linux_x64"`
		MacX64     string `json:"mac_x64"`
	} `json:"jre"`
}

// Fetch busca o release.json do repositório remoto.
func Fetch() (*Info, error) {
	resp, err := http.Get(releaseURL) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar release.json: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("release.json retornou status %d", resp.StatusCode)
	}

	var info Info
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("erro ao decodificar release.json: %w", err)
	}
	return &info, nil
}

// LocalVersion lê a versão do simulador.jar instalada localmente.
// Retorna string vazia se não houver versão instalada.
func LocalVersion() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	data, err := os.ReadFile(filepath.Join(home, ".hubsaude", "simulador-version.txt"))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

// JarPath retorna o caminho local do simulador.jar.
func JarPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	p := filepath.Join(home, ".hubsaude", "simulador.jar")
	if _, err := os.Stat(p); err != nil {
		return "", fmt.Errorf("simulador.jar não encontrado — execute `simulador obter`")
	}
	return p, nil
}

// Download baixa o simulador.jar da URL informada e salva em ~/.hubsaude/simulador.jar.
func Download(info *Info) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("não foi possível determinar o diretório home: %w", err)
	}
	dir := filepath.Join(home, ".hubsaude")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório %s: %w", dir, err)
	}

	fmt.Fprintf(os.Stderr, "Baixando simulador.jar v%s...\n", info.Jar.Version)

	resp, err := http.Get(info.Jar.URL) //nolint:gosec
	if err != nil {
		return fmt.Errorf("erro ao baixar simulador.jar: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download falhou com status HTTP %d", resp.StatusCode)
	}

	dest := filepath.Join(dir, "simulador.jar")
	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo: %w", err)
	}
	if _, err := io.Copy(f, resp.Body); err != nil {
		f.Close()
		return fmt.Errorf("erro ao salvar simulador.jar: %w", err)
	}
	f.Close()

	// Salva a versão instalada
	os.WriteFile(filepath.Join(dir, "simulador-version.txt"), //nolint:errcheck
		[]byte(info.Jar.Version), 0644)

	fmt.Fprintf(os.Stderr, "simulador.jar v%s instalado em %s\n", info.Jar.Version, dest)
	return nil
}

// JREURLForPlatform retorna a URL do JRE para a plataforma atual a partir do Info.
// isZip indica se o arquivo é zip (Windows) ou tar.gz.
func JREURLForPlatform(info *Info) (url string, isZip bool) {
	switch runtime.GOOS {
	case "windows":
		return info.JRE.WindowsX64, true
	case "darwin":
		return info.JRE.MacX64, false
	default:
		return info.JRE.LinuxX64, false
	}
}
