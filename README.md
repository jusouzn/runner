# Runner 

Repositório do **trabalho prático da disciplina de Implementação e Integração de Software**: desenvolvimento do **Sistema Runner** ao longo das aulas, em grupo de **até 2 integrantes**.

## Equipe

- Veronica Ribeiro Oliveira Palmeira
- Júlia de Souza Nascimento

## Documentação do projeto

| Documento | Descrição |
|-----------|-----------|
| [especificacao.md](especificacao.md) | Especificação completa: escopo, requisitos (US-01 a US-05), entregáveis |
| [design.md](design.md) | Arquitetura C4, fluxos, contrato de comunicação, riscos |
| [BACKLOG.md](BACKLOG.md) | Backlog por história de usuário com tarefas priorizadas |
| [PLANO_ACAO.md](PLANO_ACAO.md) | Plano de ação da equipe com marcos e cronograma |

## Objetivo do Sistema Runner

Facilitar o uso (via CLI) de aplicações Java do ecossistema da disciplina, sem exigir instalação ou configuração manual de Java.

**Componentes:**

- **`assinatura`** (CLI em Go): cria e valida assinaturas, invoca `assinador.jar` em modo local ou HTTP
- **`assinador.jar`** (Java 21): valida parâmetros rigorosamente e simula operações de assinatura
- **`simulador`** (CLI em Go): gerencia o ciclo de vida do Simulador do HubSaúde

## Especificação oficial

A especificação original do trabalho está em:

- <https://github.com/kyriosdata/runner/tree/main>

## Como gerar os diagramas C4

```bash
# Linux / macOS
./geraimagens.sh

# Windows
geraimagens.bat
```

Os scripts baixam o `plantuml.jar` automaticamente se necessário. As imagens SVG são salvas em `diagramas/imagens/`.
