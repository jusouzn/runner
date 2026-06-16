# ADR-0003: CLIs implementados em Go

- **Status:** Aceito
- **Data:** 2026-06-16
- **Decisores:** Verônica Ribeiro Oliveira Palmeira, Júlia de Souza Nascimento

## Contexto

O Sistema Runner expõe duas ferramentas de linha de comando — `assinatura` e
`simulador` — cujo papel é **orquestrar** aplicações Java do ecossistema da disciplina
(o `assinador.jar` e o `simulador.jar`) sem exigir que o usuário final instale ou
configure Java manualmente.

Requisitos que pressionam a escolha da linguagem dos CLIs:

- **Distribuição trivial para o usuário final:** idealmente um único binário, sem
  runtime a instalar — afinal, o propósito do Runner é justamente *esconder* a
  complexidade de provisionar Java.
- **Suporte multiplataforma** (Linux, Windows, macOS) a partir de um único ambiente de
  CI (ver [release.yml](../../.github/workflows/release.yml)).
- **Gerência de processos do SO:** iniciar/parar processos filhos, capturar `stdout`,
  controlar PIDs e portas (ver [ADR-0002](0002-descoberta-de-instancia.md)).
- **Cliente HTTP** simples para o modo servidor e **download de arquivos** (JRE,
  `simulador.jar`) para o provisionamento automático.

Observação importante: a escolha de Go aplica-se **aos CLIs**. O componente
`assinador.jar` permanece em **Java 21** por restrição de projeto e por exigir
`SunPKCS11` para a integração PKCS#11 (ver [design.md](../../design.md), seção 2) —
esta decisão não o altera.

## Decisão

Os CLIs `assinatura` e `simulador` são implementados em **Go (1.22+)**.

## Alternativas consideradas

### Java

- **Contra:** exigiria uma JVM para *rodar o próprio orquestrador*, contradizendo o
  objetivo central do Runner (poupar o usuário de instalar Java). A distribuição
  envolveria JAR + runtime ou empacotamento via `jlink`/`jpackage`, aumentando tamanho e
  complexidade. *Startup* mais lento para um utilitário de linha de comando de vida curta.
- **A favor:** reuso de conhecimento da equipe e do ecossistema do assinador.
- **Veredito:** descartada para os CLIs; mantida apenas para o `assinador.jar`, onde é
  requisito.

### Python

- **Contra:** depende de um interpretador instalado e na versão correta no ambiente do
  usuário, ou de empacotadores (PyInstaller) que produzem artefatos grandes e frágeis em
  cross-compiling. Distribuição multiplataforma confiável é trabalhosa.
- **A favor:** prototipagem rápida e legibilidade.
- **Veredito:** descartada — atrito de distribuição incompatível com o objetivo de
  "baixar um binário e usar".

### Go (escolhida)

- **A favor:**
  - **Binário único, estaticamente ligado**, sem runtime externo — o usuário baixa e
    executa.
  - **Cross-compiling nativo** via `GOOS`/`GOARCH`, gerando os três alvos
    (Linux/Windows/macOS) a partir de um único runner Ubuntu no CI.
  - **Biblioteca padrão rica** para o que o Runner precisa: `os/exec` (processos),
    `net/http` (cliente HTTP e downloads), manipulação de arquivos e JSON — minimizando
    dependências de terceiros.
  - *Startup* praticamente instantâneo, adequado a um utilitário de vida curta.
- **Contra:** introduz uma segunda linguagem no repositório (além de Java). Aceitável
  dado o ganho em distribuição.

## Consequências

**Positivas**

- Usuário final obtém um **binário autocontido** por plataforma, sem instalar Go nem Java.
- Pipeline de release simples: cross-compile dos três alvos em um único job
  ([release.yml](../../.github/workflows/release.yml)).
- Baixo acoplamento a dependências externas graças à stdlib.

**Negativas / mitigações**

- **Duas linguagens no projeto** (Go + Java) exigem familiaridade da equipe com ambas.
  **Mitigação:** fronteira clara — Go orquestra, Java assina; a comunicação entre eles é
  um contrato HTTP/JSON estável (ver [INTEGRACAO.md](../../INTEGRACAO.md)).
- Ferramentas de assinatura/criptografia de baixo nível (PKCS#11) permanecem do lado
  Java, então os CLIs não absorvem essa responsabilidade — o que é intencional.
