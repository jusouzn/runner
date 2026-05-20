package invoker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// HTTP sends a request to a running assinador.jar server.
// baseURL should be "http://localhost:<port>".
func HTTP(baseURL string, p Params) (string, error) {
	var endpoint string
	switch p.Operacao {
	case "assinar":
		endpoint = "/sign"
	case "validar":
		endpoint = "/validate"
	default:
		return "", fmt.Errorf("operação inválida: %s", p.Operacao)
	}

	payload := map[string]string{
		"conteudo": p.Conteudo,
		"pin":      p.Pin,
	}
	if p.Algoritmo != "" {
		payload["algoritmo"] = p.Algoritmo
	}
	if p.Assinatura != "" {
		payload["assinatura"] = p.Assinatura
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(baseURL+endpoint, "application/json", bytes.NewReader(body)) //nolint:gosec
	if err != nil {
		return "", fmt.Errorf("erro ao conectar ao servidor assinador: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("erro ao ler resposta: %w", err)
	}

	if resp.StatusCode >= 400 {
		return "", extractError(*bytes.NewBuffer(respBytes))
	}
	return strings.TrimSpace(string(respBytes)), nil
}
