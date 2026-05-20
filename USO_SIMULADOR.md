# Uso do CLI `simulador`

Este documento descreve como usar o CLI **`simulador`** para gerenciar o ciclo de vida
do **Simulador do HubSaúde** (`simulador.jar`), conforme a US-03 da
[especificação](especificacao.md).

O CLI cuida automaticamente do **download do JAR**, do **provisionamento do JRE 21**
e da **inicialização/encerramento** do processo Java, eliminando a necessidade de
configurar o ambiente manualmente.

---

## 1. Instalação

Baixe o binário pré-compilado para sua plataforma a partir do GitHub Releases:

| Plataforma   | Binário                                  |
|--------------|------------------------------------------|
| Windows x64  | `simulador-<versão>-windows-amd64.exe`   |
| Linux x64    | `simulador-<versão>-linux-amd64`         |
| macOS x64    | `simulador-<versão>-macos-amd64`         |

> Em Linux/macOS lembre-se de conceder permissão de execução: `chmod +x simulador-*`.

### Compilação a partir do fonte

```bash
cd projetos/simulador-go
go build -o simulador .
```

---

## 2. Visão geral dos comandos

```text
simulador [comando] [flags]

Comandos disponíveis:
  obter      Baixa o simulador.jar e o JRE (se necessário)
  iniciar    Inicia o simulador.jar
  parar      Para o simulador.jar em execução
  status     Exibe o status atual do simulador
  version    Exibe a versão do CLI
  help       Ajuda sobre qualquer comando
```

Use `simulador <comando> --help` para detalhes de cada comando.

---

## 3. Primeiro uso (quickstart)

```bash
# 1. Baixa o simulador.jar (versão mais recente) e o JRE 21, se necessário
simulador obter

# 2. Inicia o simulador na porta padrão (8443)
simulador iniciar

# 3. Consulta status
simulador status

# 4. Encerra o simulador
simulador parar
```

Todos os artefatos baixados ficam em `~/.hubsaude/`:

```
~/.hubsaude/
├── jdk/                       # JRE 21 provisionado automaticamente
├── simulador.jar              # JAR mais recente
├── simulador-version.txt      # Versão instalada
└── simulador-server.json      # PID e porta do processo em execução
```

---

## 4. Comandos em detalhe

### 4.1. `simulador obter`

Verifica a versão remota disponível (consultando o `release.json` do repositório
oficial) e baixa o `simulador.jar` apenas se a versão local estiver desatualizada.
Também provisiona o JRE 21 (Eclipse Temurin / Adoptium) caso não esteja presente.

```bash
simulador obter
```

**Saída típica (atualização):**
```text
Verificando versão mais recente do simulador...
Atualizando de v1.1.0 → v1.2.0
Baixando simulador.jar v1.2.0...
JRE já disponível.
```

**Saída típica (já atualizado):**
```text
Verificando versão mais recente do simulador...
simulador.jar já está na versão mais recente (1.2.0).
JRE já disponível.
```

> Requer acesso à internet. Se a rede estiver indisponível e o JAR já existir
> localmente, é possível pular este passo e ir direto para `simulador iniciar`.

---

### 4.2. `simulador iniciar`

Inicia o `simulador.jar` em **background** na porta padrão **8443**.

```bash
simulador iniciar
simulador iniciar --porta 9443     # porta customizada
```

**Verificações automáticas:**
1. Detecta se já existe instância rodando na porta — se sim, não inicia nova;
2. Verifica se a porta está livre (TCP `localhost:<porta>`);
3. Localiza o JRE 21 em `~/.hubsaude/jdk` (ou no `JAVA_HOME`);
4. Executa `java -jar ~/.hubsaude/simulador.jar`;
5. Aguarda até **30 segundos** o processo aceitar conexões TCP.

**Flags:**

| Flag       | Default | Descrição                          |
|------------|---------|------------------------------------|
| `--porta`  | `8443`  | Porta onde o simulador irá escutar |

**Erros comuns:**

| Mensagem                                                             | Causa / Solução                                          |
|----------------------------------------------------------------------|----------------------------------------------------------|
| `porta 8443 está em uso por outro processo`                          | Outro serviço ocupa a porta — libere ou use `--porta`    |
| `simulador.jar não encontrado — execute `simulador obter``           | Execute `simulador obter` primeiro                       |
| `Dica: execute `simulador obter` para provisionar o JRE`             | JRE 21 ausente — execute `simulador obter`               |
| `simulador iniciado (PID …) mas não respondeu em 30s`                | Falha de boot do JAR — consulte stderr do processo       |

---

### 4.3. `simulador status`

Consulta o endpoint `/api/info` do simulador e exibe o JSON formatado.

```bash
simulador status
simulador status --porta 9443
```

**Saída — simulador em execução:**
```text
Versão instalada: 1.2.0
Simulador em execução na porta 8443:
{
  "app": "Simulador HubSaúde",
  "version": "1.2.0",
  "uptime": "00:12:34"
}
```

**Saída — simulador parado:**
```text
Versão instalada: 1.2.0
Simulador não está em execução (porta 8443).
```

> O CLI aceita certificados auto-assinados (TLS) do simulador. A consulta tenta
> `https://` primeiro e cai para `http://` se necessário.

---

### 4.4. `simulador parar`

Encerra o simulador. A estratégia é:

1. **Graceful shutdown** via `POST /shutdown` (HTTPS e HTTP);
2. Caso o endpoint não responda, **encerra pelo PID** salvo em
   `~/.hubsaude/simulador-server.json`.

```bash
simulador parar
simulador parar --porta 9443
```

**Saída típica:**
```text
Encerrando simulador na porta 8443...
Simulador encerrado.
```

Se nada estiver rodando:
```text
Simulador não está em execução na porta 8443.
```

---

### 4.5. `simulador version`

Exibe a versão do CLI (injetada via `ldflags` em builds de release).

```bash
$ simulador version
simulador v0.1.0
```

---

## 5. Estado e diretórios

| Caminho                                  | Conteúdo                                              |
|------------------------------------------|-------------------------------------------------------|
| `~/.hubsaude/jdk/`                       | JRE 21 baixado do Adoptium                            |
| `~/.hubsaude/simulador.jar`              | JAR do simulador (versão mais recente)                |
| `~/.hubsaude/simulador-version.txt`      | Texto plano com a versão local                        |
| `~/.hubsaude/simulador-server.json`      | `{ "pid": <int>, "port": <int> }` enquanto rodando    |

Para **resetar completamente** o ambiente:

```bash
# Linux/macOS
rm -rf ~/.hubsaude

# Windows (PowerShell)
Remove-Item -Recurse -Force $HOME\.hubsaude
```

---

## 6. Endpoints HTTP do simulador

> Endpoints expostos pelo `simulador.jar` e consumidos pelo CLI. Provedor: aplicação
> Web do Simulador do HubSaúde.

| Método | Caminho      | Uso pelo CLI                |
|--------|--------------|-----------------------------|
| `GET`  | `/api/info`  | `simulador status`          |
| `POST` | `/shutdown`  | `simulador parar` (graceful)|

A porta padrão é **8443** (HTTPS com certificado auto-assinado).

---

## 7. Receitas

### 7.1. Subir o simulador em CI

```bash
simulador obter
simulador iniciar --porta 8443
# ... rodar testes que dependem do simulador ...
simulador parar --porta 8443
```

### 7.2. Diagnosticar “não inicia”

```bash
# Veja se a porta está livre
simulador status --porta 8443

# Tente uma porta alternativa
simulador iniciar --porta 18443
```

### 7.3. Forçar atualização do JAR

```bash
rm ~/.hubsaude/simulador.jar ~/.hubsaude/simulador-version.txt
simulador obter
```

---

## 8. Solução de problemas

| Sintoma                                                        | Causa provável                              | Ação                                                          |
|----------------------------------------------------------------|---------------------------------------------|----------------------------------------------------------------|
| `não foi possível verificar a versão remota`                   | Sem acesso à internet ou GitHub indisponível| Verifique a rede; o JAR local (se existir) ainda funciona     |
| `release.json não contém URL do JRE para esta plataforma`      | Plataforma não suportada (arch != x64)      | Instale o JRE 21 manualmente e exponha em `JAVA_HOME`         |
| `simulador iniciado (PID N) mas não respondeu em 30s`          | Falha na inicialização do JAR               | Inspecione `stderr`; tente outra porta                        |
| `simulador não está respondendo na porta 8443`                 | Processo travado / encerrado abruptamente   | `simulador parar` (fallback por PID) e `simulador iniciar`    |

---

## 9. Mapeamento código ↔ comandos

| Comando            | Implementação                                                                    |
|--------------------|----------------------------------------------------------------------------------|
| `simulador obter`  | [projetos/simulador-go/cmd/obter.go](projetos/simulador-go/cmd/obter.go)         |
| `simulador iniciar`| [projetos/simulador-go/cmd/iniciar.go](projetos/simulador-go/cmd/iniciar.go)    |
| `simulador parar`  | [projetos/simulador-go/cmd/parar.go](projetos/simulador-go/cmd/parar.go)        |
| `simulador status` | [projetos/simulador-go/cmd/status.go](projetos/simulador-go/cmd/status.go)      |
| `simulador version`| [projetos/simulador-go/cmd/version.go](projetos/simulador-go/cmd/version.go)    |
| Ciclo de processo  | [projetos/simulador-go/internal/process/manager.go](projetos/simulador-go/internal/process/manager.go) |
| Download do JAR    | [projetos/simulador-go/internal/release/manager.go](projetos/simulador-go/internal/release/manager.go) |
| Provisionamento JRE| [projetos/simulador-go/internal/jdk/provisioner.go](projetos/simulador-go/internal/jdk/provisioner.go) |
