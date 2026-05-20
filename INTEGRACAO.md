# Integração — `assinatura` ↔ `assinador.jar`

Este documento descreve o contrato técnico de integração entre o CLI **`assinatura`** (Go)
e a aplicação **`assinador.jar`** (Java), incluindo o protocolo HTTP, os payloads JSON
trocados e exemplos práticos de uso.

> Os fluxos lógicos são idênticos nos dois modos suportados (Local/CLI e HTTP). Apenas
> o mecanismo de transporte muda. Veja [design.md](design.md) seções 4 e 5 para a
> visão arquitetural e [especificacao.md](especificacao.md) seção 6 para o escopo
> funcional da integração.

---

## 1. Modos de invocação

| Modo | Como o `assinatura` chama o `assinador.jar` | Quando usar |
|------|---------------------------------------------|-------------|
| **Local (cold start)** | `java -jar assinador.jar --modo=local …` | Execuções pontuais / scripts |
| **HTTP (warm start)**  | `POST http://localhost:<porta>/sign` ou `/validate` | Múltiplas chamadas, menor latência |

Em ambos os casos, o **formato dos parâmetros lógicos é o mesmo** (`conteudo`, `pin`,
`algoritmo`, `assinatura`); muda apenas como eles são transmitidos (flags de CLI vs.
campos JSON).

---

## 2. Contrato HTTP (modo servidor)

### 2.1. Inicialização do servidor

Comando executado pelo `assinatura` ao iniciar o servidor:

```bash
java -jar assinador.jar --modo=servidor --porta=8080
```

- **Porta padrão:** `8080`
- **Bind address:** `localhost` (apenas loopback — não escuta em interfaces externas)
- **Sinal de prontidão:** ao subir, escreve em `stderr`:
  ```json
  {"status":"running","port":8080}
  ```
- **Executor:** thread virtual por requisição (JDK 21).

### 2.2. Endpoints

| Método | Caminho      | Descrição                                  | Códigos esperados      |
|--------|--------------|--------------------------------------------|------------------------|
| `GET`  | `/health`    | Health check (liveness)                    | `200`, `405`           |
| `POST` | `/sign`      | Cria assinatura simulada                   | `200`, `400`, `405`    |
| `POST` | `/validate`  | Valida assinatura simulada                 | `200`, `400`, `405`    |

Todas as respostas são `application/json; charset=utf-8`.

### 2.3. Health check — `GET /health`

**Requisição**
```http
GET /health HTTP/1.1
Host: localhost:8080
```

**Resposta 200 OK**
```json
{"status":"ok"}
```

O `assinatura` usa este endpoint em [`server.EnsureRunning`](projetos/assinatura-go/internal/server/manager.go#L26) para detectar se há uma instância já em execução antes de iniciar uma nova.

---

## 3. Operação `sign` — `POST /sign`

### 3.1. Requisição

```http
POST /sign HTTP/1.1
Host: localhost:8080
Content-Type: application/json

{
  "conteudo": "texto a ser assinado",
  "pin": "1234",
  "algoritmo": "SHA256withRSA"
}
```

| Campo       | Tipo   | Obrigatório | Restrições                                                                 |
|-------------|--------|-------------|----------------------------------------------------------------------------|
| `conteudo`  | string | Sim         | Não vazio; até **10.000 caracteres**                                       |
| `pin`       | string | Sim         | Apenas dígitos; **mínimo 4** caracteres                                    |
| `algoritmo` | string | Não         | Um de: `SHA256withRSA` (default), `SHA512withRSA`, `SHA1withRSA`           |

### 3.2. Resposta 200 OK

```json
{
  "status": "success",
  "operation": "sign",
  "data": {
    "signature": "MIIB...[base64-simulado-SHA256withRSA]...AAAA==",
    "algorithm": "SHA256withRSA"
  }
}
```

### 3.3. Resposta 400 Bad Request

```json
{
  "status": "error",
  "error_code": "ERR_INVALID_PARAM",
  "message": "O parâmetro --pin possui tamanho incorreto (deve ter no mínimo 4 dígitos)",
  "details": []
}
```

### 3.4. Exemplo com `curl`

```bash
curl -sS -X POST http://localhost:8080/sign \
  -H "Content-Type: application/json" \
  -d '{"conteudo":"olá mundo","pin":"1234","algoritmo":"SHA256withRSA"}'
```

---

## 4. Operação `validate` — `POST /validate`

### 4.1. Requisição

```http
POST /validate HTTP/1.1
Host: localhost:8080
Content-Type: application/json

{
  "conteudo": "texto original",
  "assinatura": "MIIB...AAAA==",
  "pin": "1234"
}
```

| Campo        | Tipo   | Obrigatório | Restrições                                  |
|--------------|--------|-------------|---------------------------------------------|
| `conteudo`   | string | Sim         | Não vazio; até 10.000 caracteres            |
| `pin`        | string | Sim         | Apenas dígitos; mínimo 4 caracteres         |
| `assinatura` | string | Sim         | Não vazio                                   |

### 4.2. Resposta 200 OK

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

> ⚠️ Trata-se de **simulação**. A resposta de validação é pré-determinada (`valid: true`)
> quando os parâmetros são válidos, conforme [especificacao.md §8.1](especificacao.md).

### 4.3. Resposta 400 Bad Request

```json
{
  "status": "error",
  "error_code": "ERR_INVALID_PARAM",
  "message": "O parâmetro --assinatura é obrigatório",
  "details": []
}
```

### 4.4. Exemplo com `curl`

```bash
curl -sS -X POST http://localhost:8080/validate \
  -H "Content-Type: application/json" \
  -d '{"conteudo":"olá mundo","assinatura":"MIIB...AAAA==","pin":"1234"}'
```

---

## 5. Modo Local (CLI / cold start)

Quando o `assinatura` é chamado com `--local`, ele executa diretamente o JAR via
`java -jar`. Os parâmetros são passados como flags `--chave=valor`.

### 5.1. Comando montado

```bash
java -jar assinador.jar \
  --modo=local \
  --operacao=assinar \
  --pin=1234 \
  --conteudo="olá mundo" \
  --algoritmo=SHA256withRSA
```

Para validação:

```bash
java -jar assinador.jar \
  --modo=local \
  --operacao=validar \
  --pin=1234 \
  --conteudo="olá mundo" \
  --assinatura="MIIB...AAAA=="
```

### 5.2. Canais de saída

| Canal    | Conteúdo                                                          |
|----------|-------------------------------------------------------------------|
| `stdout` | JSON de sucesso (mesma estrutura da resposta HTTP 200)            |
| `stderr` | JSON de erro (mesma estrutura da resposta HTTP 400) + diagnósticos |
| `exit`   | `0` em sucesso; `1` em erro                                       |

### 5.3. Parâmetros aceitos

| Flag             | Obrigatório (assinar) | Obrigatório (validar) | Observação                       |
|------------------|-----------------------|-----------------------|----------------------------------|
| `--modo`         | Sim                   | Sim                   | `local` ou `servidor`            |
| `--operacao`     | Sim                   | Sim                   | `assinar` ou `validar`           |
| `--pin`          | Sim                   | Sim                   | mínimo 4 dígitos numéricos       |
| `--conteudo`     | Sim                   | Sim                   | ≤ 10.000 caracteres              |
| `--algoritmo`    | Não                   | —                     | default `SHA256withRSA`          |
| `--assinatura`   | —                     | Sim                   | —                                |
| `--porta`        | —                     | —                     | usado apenas em `--modo=servidor` (default `8080`) |

---

## 6. Contrato JSON canônico

Todos os retornos do `assinador.jar` — seja por HTTP, seja pelo `stdout`/`stderr` no
modo local — seguem **o mesmo esquema canônico** definido a seguir.

### 6.1. Envelope de sucesso

```json
{
  "status": "success",
  "operation": "sign | validate",
  "data": { /* específico da operação */ }
}
```

### 6.2. Envelope de erro

```json
{
  "status": "error",
  "error_code": "ERR_INVALID_PARAM",
  "message": "descrição legível do problema",
  "details": []
}
```

### 6.3. Códigos de erro

| `error_code`          | Causa típica                                                                 |
|-----------------------|------------------------------------------------------------------------------|
| `ERR_INVALID_PARAM`   | Parâmetro ausente, vazio, fora do domínio ou com formato inválido            |

> Códigos adicionais podem ser introduzidos em versões futuras seguindo o prefixo `ERR_`.

---

## 7. Códigos HTTP

| Status | Significado                                                                     |
|--------|---------------------------------------------------------------------------------|
| `200`  | Operação executada com sucesso (corpo segue o envelope de sucesso)              |
| `400`  | Falha de validação de parâmetros (corpo segue o envelope de erro)               |
| `405`  | Método não permitido para o endpoint (ex.: `GET /sign`)                         |

> Erros de transporte (servidor indisponível, timeout, conexão recusada) são tratados
> pelo `assinatura` como falhas do CLI e exibidos com a dica de uso de `--local`.

---

## 8. Gerenciamento de ciclo de vida (servidor)

### 8.1. Inicialização sob demanda

Antes de cada requisição HTTP, o `assinatura` executa:

1. `GET /health` para verificar se há instância viva (`isAlive`);
2. Se não houver, executa `java -jar assinador.jar --modo=servidor --porta=8080`
   em background e aguarda até 20 segundos pelo `/health` responder;
3. Persiste estado em `~/.hubsaude/assinador-server.json`:

   ```json
   { "pid": 12345, "port": 8080 }
   ```

### 8.2. Encerramento

Atualmente o encerramento ocorre via sinal ao PID registrado (`SIGINT` em
Unix, `Kill` em Windows). Veja [`server.Stop`](projetos/assinatura-go/internal/server/manager.go#L34).

---

## 9. Exemplos ponta a ponta

### 9.1. Assinatura via HTTP (default)

```bash
$ assinatura sign --conteudo "olá mundo" --pin 1234
{
  "status": "success",
  "operation": "sign",
  "data": {
    "signature": "MIIB...[base64-simulado-SHA256withRSA]...AAAA==",
    "algorithm": "SHA256withRSA"
  }
}
```

### 9.2. Assinatura local (cold start)

```bash
$ assinatura sign --local --conteudo "olá mundo" --pin 1234
{ "status":"success", "operation":"sign", "data": { ... } }
```

### 9.3. Validação com PIN inválido

```bash
$ assinatura validate --conteudo "x" --assinatura "MIIB..." --pin 12
Error: O parâmetro --pin possui tamanho incorreto (deve ter no mínimo 4 dígitos)
```

### 9.4. Chamada HTTP direta (sem o CLI)

```bash
$ curl -sS -X POST http://localhost:8080/sign \
    -H "Content-Type: application/json" \
    -d '{"conteudo":"texto","pin":"1234"}' | jq
```

---

## 10. Mapeamento código ↔ contrato

| Item do contrato                  | Implementação                                                              |
|-----------------------------------|----------------------------------------------------------------------------|
| Servidor HTTP e endpoints         | [projetos/assinador-java/src/main/java/com/runner/assinador/ServidorHttp.java](projetos/assinador-java/src/main/java/com/runner/assinador/ServidorHttp.java) |
| Validação de parâmetros           | [projetos/assinador-java/src/main/java/com/runner/assinador/FakeSignatureService.java](projetos/assinador-java/src/main/java/com/runner/assinador/FakeSignatureService.java) |
| Parser de flags do JAR            | [projetos/assinador-java/src/main/java/com/runner/assinador/AssinadorMain.java](projetos/assinador-java/src/main/java/com/runner/assinador/AssinadorMain.java) |
| Cliente HTTP (Go)                 | [projetos/assinatura-go/internal/invoker/http.go](projetos/assinatura-go/internal/invoker/http.go) |
| Invocação local (Go)              | [projetos/assinatura-go/internal/invoker/local.go](projetos/assinatura-go/internal/invoker/local.go) |
| Gestão do servidor (Go)           | [projetos/assinatura-go/internal/server/manager.go](projetos/assinatura-go/internal/server/manager.go) |
