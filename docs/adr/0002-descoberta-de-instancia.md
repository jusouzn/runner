# ADR-0002: Descoberta de instância via arquivo de estado em `~/.hubsaude/`

- **Status:** Aceito
- **Data:** 2026-06-16
- **Decisores:** Verônica Ribeiro Oliveira Palmeira, Júlia de Souza Nascimento

## Contexto

Tanto o CLI `assinatura` (modo servidor) quanto o CLI `simulador` gerenciam processos
Java de longa duração (warm start). Cada invocação do CLI é um **processo novo e
efêmero** que precisa descobrir se já existe uma instância do serviço em execução para:

- reaproveitá-la em vez de iniciar outra (evitando *cold starts* e conflitos de porta);
- saber em qual porta ela está escutando;
- conseguir encerrá-la depois (`servidor parar` / `simulador parar`), o que exige
  conhecer o **PID** do processo.

Como processos efêmeros do CLI não compartilham memória entre si, o estado precisa ser
**persistido fora do processo**, em um local previsível, por usuário, e independente do
diretório de trabalho corrente.

O projeto já concentra todos os seus artefatos gerenciados (JRE provisionado,
`assinador.jar`, `simulador.jar`, versões em cache) no diretório `~/.hubsaude/`
(ver [JDK.md](../../JDK.md) e [design.md](../../design.md), seção 5).

## Decisão

A descoberta de instância é feita por um **arquivo de estado JSON** gravado no diretório
`~/.hubsaude/` do usuário (resolvido a partir do *home directory* do SO):

- `assinatura`: `~/.hubsaude/assinador-server.json`
  ([`internal/server/manager.go`](../../projetos/assinatura-go/internal/server/manager.go)).
- `simulador`: `~/.hubsaude/simulador-server.json`
  ([`internal/process/manager.go`](../../projetos/simulador-go/internal/process/manager.go)).

O conteúdo é um objeto JSON com o PID e a porta da instância:

```json
{ "pid": 12345, "port": 8080 }
```

Fluxo de uso:

1. Ao iniciar, o CLI **salva** o estado (PID + porta) após subir o servidor.
2. Em invocações seguintes, o CLI **lê** o estado e/ou sonda a porta para decidir entre
   reaproveitar a instância ou iniciar uma nova.
3. Ao encerrar (`parar`), o CLI lê o PID do estado, encerra o processo e **remove** o
   arquivo de estado.

## Alternativas consideradas

- **Sondar apenas a porta (sem arquivo):** detecta "algo escutando", mas não revela o
  PID — inviabilizando o encerramento controlado e a distinção entre nossa instância e
  um processo alheio na mesma porta. Descartada como solução única (a sondagem de porta
  é usada de forma **complementar** ao arquivo).
- **PID file clássico (`/var/run` ou `/tmp`):** `/var/run` exige privilégios; `/tmp` é
  volátil, compartilhado entre usuários e sujeito a colisões. Descartada.
- **Banco de dados / serviço externo:** sobrecarga desproporcional para um trabalho
  acadêmico; adicionaria dependências. Descartada.
- **Variáveis de ambiente:** não sobrevivem entre invocações independentes do CLI.
  Descartada.

## Consequências

**Positivas**

- Local **único e previsível** por usuário, coerente com os demais artefatos em
  `~/.hubsaude/`, sem poluir o diretório de trabalho.
- Formato JSON simples, legível e sem dependências externas.
- Permite encerramento controlado por PID e coerência CLI↔servidor quanto à porta
  (inclusive quando a porta é customizada — ver [ADR-0001](0001-porta-padrao-assinador.md)).
- Isolamento por usuário: cada `$HOME` tem seu próprio estado.

**Negativas / mitigações**

- **Estado obsoleto (stale):** se o processo morrer abruptamente, o arquivo pode apontar
  para um PID/porta inválidos. **Mitigação:** o CLI valida a instância sondando a porta
  antes de confiar no estado, e recria/remove o arquivo conforme necessário.
- **Concorrência:** invocações simultâneas poderiam disputar o arquivo. Aceitável no
  uso interativo previsto (um operador por vez por `$HOME`).
- **Suposição de `$HOME` gravável:** ambientes sem home definido falhariam. Tratado como
  erro explícito na resolução do caminho.
