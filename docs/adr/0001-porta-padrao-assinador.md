# ADR-0001: Porta padrão 8080 do servidor do assinador

- **Status:** Aceito
- **Data:** 2026-06-16
- **Decisores:** Verônica Ribeiro Oliveira Palmeira, Júlia de Souza Nascimento

## Contexto

O CLI `assinatura`, em modo servidor (warm start), inicia o `assinador.jar` como um
servidor HTTP local de longa duração e reaproveita essa instância em chamadas
subsequentes para reduzir a latência do *cold start* da JVM
(ver [design.md](../../design.md), seção 3.2).

Para que o CLI e o `assinador.jar` se comuniquem sem configuração manual por parte do
usuário final, é necessário definir uma **porta padrão** conhecida por ambos os lados.
A porta precisa de um valor previsível, mas também configurável, já que o ambiente do
usuário pode tê-la ocupada por outro processo.

Restrições e forças consideradas:

- O servidor escuta apenas em `localhost` (comunicação entre processos na mesma máquina).
- Não há requisito de TLS para o assinador (diferente do simulador, que usa 8443).
- O valor deve estar na faixa de portas não privilegiadas (≥ 1024) para não exigir
  privilégios de administrador.
- Deve ser um valor de convenção amplamente reconhecido como "HTTP local de aplicação".

## Decisão

A porta padrão do servidor HTTP do `assinador.jar` é a **8080**.

- O CLI Go define `const DefaultPort = 8080`
  ([`internal/server/manager.go`](../../projetos/assinatura-go/internal/server/manager.go)).
- O `assinador.jar` adota `8080` como default quando `--porta` não é informado
  ([`AssinadorMain.java`](../../projetos/assinador-java/src/main/java/com/runner/assinador/AssinadorMain.java)).
- O valor é **configurável** via flag `--porta=<n>`, repassada do CLI ao processo Java,
  permitindo contornar conflitos de porta sem alteração de código.

## Alternativas consideradas

- **Porta efêmera/aleatória (porta 0):** o SO atribuiria uma porta livre e o CLI a leria
  do arquivo de estado. Elimina conflitos, mas torna o comportamento menos previsível
  para depuração manual (`curl localhost:8080/...`) e documentação. Descartada por
  reduzir a transparência didática do projeto.
- **Porta 80 / 443:** privilegiadas (< 1024), exigiriam elevação de privilégios.
  Descartadas.
- **8443:** já é a porta padrão do **simulador** (que usa TLS). Reusá-la causaria
  colisão entre os dois serviços. Descartada.
- **Porta de alto número arbitrária (ex.: 49152+):** sem ganho frente a 8080 e menos
  reconhecível como "HTTP de aplicação". Descartada.

## Consequências

**Positivas**

- Valor de convenção amplamente reconhecido para HTTP local, fácil de lembrar e documentar.
- Comportamento previsível: o usuário consegue inspecionar manualmente o endpoint.
- Não exige privilégios elevados.
- Mantém separação clara de portas entre assinador (8080) e simulador (8443).

**Negativas / mitigações**

- 8080 é uma porta popular e pode estar ocupada por outro serviço de desenvolvimento.
  **Mitigação:** flag `--porta` permite redefinir a porta; o CLI registra a porta em uso
  no arquivo de estado (ver [ADR-0002](0002-descoberta-de-instancia.md)), de modo que
  CLI e servidor permaneçam coerentes mesmo com porta customizada.
- Comunicação em texto claro (sem TLS). Aceitável porque o tráfego nunca deixa o
  `localhost` e o conteúdo é simulado neste estágio do projeto.
