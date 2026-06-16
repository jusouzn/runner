# Architecture Decision Records (ADRs)

Este diretório registra as **decisões arquiteturais** do Sistema Runner no formato
[ADR](https://adr.github.io/). Cada ADR é um documento imutável que captura o
contexto, a decisão tomada e suas consequências em um determinado momento do projeto.

## Por que ADRs?

- **Rastreabilidade:** vincula cada decisão técnica ao seu racional e ao momento em que foi tomada.
- **Onboarding:** novos integrantes entendem *por que* o sistema é como é, não apenas *como* ele é.
- **Auditoria de configuração:** decisões ficam versionadas junto ao código, sob o mesmo controle de versão.

## Convenções

- Arquivos nomeados como `NNNN-titulo-em-kebab-case.md`, com `NNNN` sequencial e zero-padded.
- Um ADR **não é editado** após aceito. Se uma decisão muda, cria-se um novo ADR que
  **supera** (`Supersedes`) o anterior, e o antigo passa a `Status: Substituído por ADR-XXXX`.
- Cada ADR segue o template abaixo.

## Índice

| ADR | Título | Status |
|-----|--------|--------|
| [0001](0001-porta-padrao-assinador.md) | Porta padrão 8080 do servidor do assinador | Aceito |
| [0002](0002-descoberta-de-instancia.md) | Descoberta de instância via arquivo de estado em `~/.hubsaude/` | Aceito |
| [0003](0003-cli-em-go.md) | CLIs implementados em Go | Aceito |

## Template

```markdown
# ADR-NNNN: Título

- **Status:** Proposto | Aceito | Substituído por ADR-XXXX | Descontinuado
- **Data:** AAAA-MM-DD
- **Decisores:** ...

## Contexto

Qual problema ou força motriz levou a esta decisão?

## Decisão

O que foi decidido (de forma afirmativa e clara).

## Alternativas consideradas

Que opções foram avaliadas e por que foram descartadas.

## Consequências

Resultados positivos e negativos decorrentes da decisão.
```
