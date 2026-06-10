# Runner 

Repositório do **trabalho prático da disciplina de Implementação e Integração de Software**: desenvolvimento do **Sistema Runner** ao longo das aulas, em grupo de **até 2 integrantes**.

## Equipe

- Veronica Ribeiro Oliveira Palmeira
- Júlia de Souza Nascimento

## Documentação do projeto

| Documento | Descrição |
|-----------|-----------|
| [especificacao.md](especificacao.md) | Ponteiro para a especificação oficial no upstream, fixada por commit |
| [design.md](design.md) | Decisões e detalhes específicos desta implementação |
| [BACKLOG.md](BACKLOG.md) | Backlog por história de usuário com tarefas priorizadas |
| [PLANO_ACAO.md](PLANO_ACAO.md) | Plano de ação da equipe com marcos e cronograma |

## Objetivo do Sistema Runner

Facilitar o uso (via CLI) de aplicações Java do ecossistema da disciplina, sem exigir instalação ou configuração manual de Java.

**Componentes:**

- **`assinatura`** (CLI em Go): cria e valida assinaturas, invoca `assinador.jar` em modo local ou HTTP
- **`assinador.jar`** (Java 21): valida parâmetros rigorosamente e simula operações de assinatura
- **`simulador`** (CLI em Go): gerencia o ciclo de vida do Simulador do HubSaúde

## Especificação oficial

A especificação original do trabalho é mantida no upstream e não é duplicada neste repositório. A referência abaixo é fixada por commit para garantir rastreabilidade e evitar deriva (não usar `main`):

- <https://github.com/kyriosdata/runner/blob/4d7d40fff32b3b50372e7fbe41fe713b2bbddb4c/especificacao.md>

## Como gerar os diagramas C4

```bash
# Linux / macOS
./geraimagens.sh

# Windows
geraimagens.bat
```

Os scripts baixam o `plantuml.jar` automaticamente se necessário. As imagens SVG são salvas em `diagramas/imagens/`.
