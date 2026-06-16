# Runner

Implementação do **Sistema Runner** — trabalho prático da disciplina de Implementação e Integração de Software.

O Runner facilita o uso, via linha de comando, de aplicações Java do ecossistema da disciplina, sem exigir instalação manual de Java. É composto por três artefatos:

- **`assinatura`** (CLI em Go): cria e valida assinaturas, invoca `assinador.jar` em modo local ou HTTP
- **`assinador.jar`** (Java 21): valida parâmetros rigorosamente e simula operações de assinatura
- **`simulador`** (CLI em Go): gerencia o ciclo de vida do Simulador do HubSaúde

## Equipe

- Verônica Ribeiro Oliveira Palmeira
- Júlia de Souza Nascimento

## Status atual

| Componente | Estado |
|------------|--------|
| `assinatura` CLI (modos local e HTTP) | Implementado |
| `assinador.jar` (validação + simulação + servidor HTTP) | Implementado |
| `simulador` CLI (start/stop/status, download via Releases) | Implementado |
| Provisionamento automático de JDK | Implementado (Linux/macOS/Windows amd64) |
| Build multiplataforma (CI) | Linux validado em CI; Windows pendente |
| Releases assinadas com Cosign + checksums SHA256 | Pipeline pronto, primeira tag pendente |
| PKCS#11 real | Pendente — escopo do marco final |

A versão estável é a `main`. Releases serão publicadas em [GitHub Releases](https://github.com/jusouzn/runner/releases) sob versionamento [SemVer](https://semver.org/lang/pt-BR/).

## Documentação do projeto

| Documento | Descrição |
|-----------|-----------|
| [especificacao.md](especificacao.md) | Ponteiro para a especificação oficial no upstream, fixada por commit |
| [design.md](design.md) | Decisões e detalhes específicos desta implementação |
| [INTEGRACAO.md](INTEGRACAO.md) | Contrato HTTP do `assinador.jar`, payloads e mapeamento de erros |
| [USO_SIMULADOR.md](USO_SIMULADOR.md) | Guia de uso do CLI `simulador` |
| [JDK.md](JDK.md) | Provisionamento automático do JDK |
| [BACKLOG.md](BACKLOG.md) | Backlog por história de usuário |
| [PLANO_ACAO.md](PLANO_ACAO.md) | Plano de ação com marcos e cronograma |

A especificação oficial é mantida no upstream e referenciada por commit fixo (não `main`):

- <https://github.com/kyriosdata/runner/blob/4d7d40fff32b3b50372e7fbe41fe713b2bbddb4c/especificacao.md>

## Pré-requisitos

- **Go 1.22+** (apenas para compilar os CLIs a partir do código-fonte)
- **JDK 21+** (apenas para compilar o `assinador.jar`; usuários finais não precisam — o `assinatura` provisiona um JDK automaticamente quando necessário)
- **Git**

> Para uso final, basta baixar o binário pré-compilado da release correspondente à sua plataforma. Nenhuma instalação adicional é exigida.

## Como gerar os executáveis

### 1. CLI `assinatura` (Go)

```bash
cd projetos/assinatura-go
go mod tidy
go build -o assinatura .          # Linux/macOS
go build -o assinatura.exe .      # Windows
```

Cross-compiling para outra plataforma:

```bash
GOOS=linux   GOARCH=amd64 go build -o assinatura-linux-amd64 .
GOOS=windows GOARCH=amd64 go build -o assinatura-windows-amd64.exe .
GOOS=darwin  GOARCH=amd64 go build -o assinatura-macos-amd64 .
```

### 2. CLI `simulador` (Go)

```bash
cd projetos/simulador-go
go mod tidy
go build -o simulador .           # Linux/macOS
go build -o simulador.exe .       # Windows
```

### 3. `assinador.jar` (Java)

```bash
cd projetos/assinador-java
./mvnw package                    # Linux/macOS
mvnw.cmd package                  # Windows
# Resultado: target/assinador.jar
```

> O Maven Wrapper (`mvnw`) baixa o Maven necessário automaticamente. Apenas o JDK 21 precisa estar instalado para esta etapa.

## Como executar

### Assinatura — modo servidor (padrão, recomendado)

```bash
./assinatura sign --conteudo "texto a assinar" --pin 1234
./assinatura validate --conteudo "texto" --assinatura "MIIB..." --pin 1234
```

O CLI inicia automaticamente o `assinador.jar` em modo servidor HTTP (warm start) e reutiliza a instância em chamadas seguintes.

### Assinatura — modo local (cold start)

```bash
./assinatura sign --local --conteudo "texto" --pin 1234
```

### Controle do servidor

```bash
./assinatura servidor parar
```

### Modo verboso

```bash
./assinatura --verbose sign --conteudo "texto" --pin 1234
```

### Simulador

```bash
./simulador obter      # baixa simulador.jar mais recente via GitHub Releases (cache local)
./simulador iniciar    # inicia o simulador (verifica portas)
./simulador status     # consulta estado/PID
./simulador parar      # encerra o processo
```

Para detalhes completos, ver [USO_SIMULADOR.md](USO_SIMULADOR.md) e [INTEGRACAO.md](INTEGRACAO.md).

## Como executar os testes

### Testes do `assinador.jar` (JUnit)

```bash
cd projetos/assinador-java
./mvnw test
```

### Testes do CLI `assinatura` (Go)

```bash
cd projetos/assinatura-go
go test ./...
```

### Testes do CLI `simulador` (Go)

```bash
cd projetos/simulador-go
go test ./...
```

### Pipeline de CI

Toda alteração em PR ou push para `main` dispara o workflow [.github/workflows/ci.yml](.github/workflows/ci.yml), que executa testes Java, testes Go e build dos binários para Linux, Windows e macOS amd64.

## Como gerar os diagramas C4

Os arquivos-fonte ficam em [diagramas/](diagramas/) (PlantUML) e as imagens SVG em `diagramas/imagens/`.

```bash
# Linux / macOS
chmod +x geraimagens.sh
./geraimagens.sh

# Windows
geraimagens.bat
```

Os scripts baixam o `plantuml.jar` automaticamente quando ausente.

## Como contribuir

Este é um trabalho acadêmico em equipe de duas pessoas; ainda assim, o repositório segue padrões abertos para tornar a contribuição navegável.

1. Crie uma issue descrevendo o problema ou melhoria, referenciando a história de usuário (US-01…US-05) quando aplicável
2. Abra uma branch a partir de `main` com prefixo `feat/`, `fix/`, `chore/`, `docs/` ou `ci/`
3. Use [Conventional Commits](https://www.conventionalcommits.org/pt-br/v1.0.0/) nas mensagens (ex.: `feat(assinatura): suportar timeout configurável`)
4. Abra um PR pequeno e atômico, ligado à issue (`Closes #N`), descrevendo o que foi feito e como validar
5. CI deve passar (lint + testes + build) antes do merge

## Licença

[MIT](LICENSE) © 2026 Verônica Ribeiro Oliveira Palmeira, Júlia de Souza Nascimento.
