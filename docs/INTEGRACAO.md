# INTEGRACAO — Modo Servidor (HTTP) do `assinador.jar`

Este documento descreve o contrato de integração entre o CLI `assinatura` (Go) e o `assinador.jar` (Java) quando este último é executado em **modo servidor** (HTTP). O modo local (cold start) está descrito no `Design.md` (seção 3.1).

> Marco 4 (05/05/2026) do `plano_de_acao.md`.

---

## 1. Visão geral

| Item | Valor |
|---|---|
| Host | `localhost` |
| Porta padrão | `8080` (configurável via `--porta`) |
| Protocolo | HTTP/1.1 |
| Content-Type | `application/json; charset=utf-8` |
| Charset | UTF-8 |
| Versão atual | `0.2.0` |

O servidor é provisionado pelo próprio CLI quando necessário (warm start). Se já houver instância ativa na porta indicada, o CLI a reutiliza.

---

## 2. Como iniciar / parar o servidor

### 2.1. Via CLI `assinatura`

```bash
# inicia (na porta padrao 8080)
assinatura servidor iniciar

# inicia em porta especifica
assinatura servidor iniciar --porta=9090

# consulta estado
assinatura servidor status

# encerra
assinatura servidor parar
```

O CLI mantém um arquivo de estado em `~/.hubsaude/assinador.json` com o PID e a porta da instância iniciada.

### 2.2. Diretamente via Java

```bash
java -jar assinador.jar --modo=server --porta=8080
```

---

## 3. Endpoints

| Método | Caminho | Descrição |
|---|---|---|
| `GET`  | `/health`    | *Liveness* — usado pelo CLI para detectar instâncias ativas. |
| `POST` | `/sign`      | Cria assinatura simulada (alias: `/assinar`). |
| `POST` | `/validate`  | Valida assinatura simulada (alias: `/validar`). |
| `POST` | `/shutdown`  | Solicita encerramento ordenado do servidor. |

### 3.1. `GET /health`

**Resposta `200 OK`**
```json
{
  "status": "success",
  "data": { "versao": "0.2.0" }
}
```

### 3.2. `POST /sign` (e alias `/assinar`)

**Request body**
```json
{
  "pin": "1234"
}
```

**Resposta `200 OK`**
```json
{
  "status": "success",
  "operation": "sign",
  "data": {
    "signature": "MIIB...[base64 simulada]...",
    "algorithm": "SHA256withRSA"
  }
}
```

**Resposta `400 Bad Request`** (parâmetro inválido)
```json
{
  "status": "error",
  "error_code": "ERR_INVALID_PARAM",
  "message": "Campo 'pin' ausente ou invalido (minimo 4 digitos).",
  "details": []
}
```

### 3.3. `POST /validate` (e alias `/validar`)

**Request body**
```json
{
  "pin": "1234",
  "assinatura": "MIIB...base64..."
}
```

**Resposta `200 OK`**
```json
{
  "status": "success",
  "operation": "validate",
  "data": {
    "valid": true,
    "algorithm": "SHA256withRSA"
  }
}
```

### 3.4. `POST /shutdown`

**Resposta `200 OK`**
```json
{ "status": "success", "message": "servidor encerrando" }
```

O processo Java é encerrado logo após enviar a resposta.

---

## 4. Mapeamento de erros

O contrato segue o formato definido em `Design.md` (seção 4):

```json
{
  "status": "error",
  "error_code": "<CODIGO>",
  "message": "<descricao legivel>",
  "details": []
}
```

| Origem | Status HTTP | `error_code` | Quando ocorre |
|---|---|---|---|
| Assinador (HTTP) | 400 | `ERR_INVALID_PARAM` | Parâmetro ausente ou em formato inválido. |
| Assinador (HTTP) | 405 | `ERR_METHOD_NOT_ALLOWED` | Método HTTP diferente de `POST` em endpoints de operação. |
| Assinador (HTTP) | 500 | `ERR_INTERNO` | Falha não prevista durante o processamento. |
| CLI | — (stderr, exit≠0) | `ERR_INVALID_PARAM` | Faltam parâmetros obrigatórios no CLI (ex.: `--pin`). |
| CLI | — (stderr, exit≠0) | `ERR_SERVIDOR_INDISPONIVEL` | Não foi possível iniciar/atingir o servidor. |
| CLI | — (stderr, exit≠0) | `ERR_HTTP` | Falha de rede ou resposta inválida. |

O CLI propaga em `stdout` o JSON recebido do servidor e finaliza com **exit code ≠ 0** quando o status HTTP for `>= 400`.

---

## 5. Exemplos ponta a ponta

### 5.1. Modo HTTP (warm start) via CLI

```bash
# 1) primeira chamada inicia automaticamente o servidor
assinatura assinar --pin=1234

# 2) chamadas subsequentes reutilizam a instancia
assinatura assinar --pin=1234
assinatura validar --pin=1234 --assinatura=MIIB...

# 3) encerramento explicito
assinatura servidor parar
```

### 5.2. Modo local (cold start) via CLI

```bash
assinatura assinar --local --pin=1234
```

### 5.3. Chamada HTTP direta (curl)

```bash
curl -s http://localhost:8080/health
curl -s -X POST http://localhost:8080/sign \
  -H "Content-Type: application/json" \
  -d '{"pin":"1234"}'
```

---

## 6. Notas e limitações (Marco 4)

- O parser JSON do servidor é simplificado (`"campo":"valor"`); será substituído por parser real em iterações futuras.
- O servidor atende somente em `localhost`, sem TLS — o uso é estritamente local.
- Não há autenticação: o servidor é controlado pelo próprio usuário em sua máquina.
- A porta padrão é `8080`; pode ser alterada com `--porta=N` em ambos os lados (CLI e jar).
