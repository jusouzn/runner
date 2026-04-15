# Runner (Trabalho Prático)

Repositório do **trabalho prático da disciplina**: desenvolvimento do **Sistema Runner** ao longo das aulas, em grupo de **até 2 integrantes**.

Este projeto tem como objetivo **dominar a implementação de software** na prática, incorporando os tópicos apresentados na disciplina (arquitetura, qualidade, testes, automação, releases, segurança de supply chain etc.) e refletindo isso no código e na organização do repositório.

## Objetivo do Sistema Runner

O **Sistema Runner** visa **facilitar o uso (via CLI) de aplicações Java** do ecossistema da disciplina, reduzindo a necessidade de conhecimento técnico sobre instalação/configuração, e padronizando:

- execução local e via servidor (HTTP) para o componente Java (`assinador.jar`);
- gerenciamento do simulador (`simulador.jar`);
- provisionamento automático de JDK;
- releases multiplataforma com integridade/verificação (checksums + assinatura).

## Especificação oficial do trabalho

A especificação/enunciado do trabalho prático está disponível em:

- https://github.com/kyriosdata/runner/tree/main

> Este repositório implementa a solução conforme o enunciado acima e as orientações do docente ao longo da disciplina.

## Equipe

- **Veronica Ribeiro Oliveira Palmeira**
- **Júlia de Souza Nascimento**

(Grupo com 2 integrantes, conforme exigência do trabalho.)

## Entregáveis (visão geral)

Conforme o planejamento do projeto, a entrega final inclui:

1. **CLI `assinatura` (multiplataforma)**
   - cria/valida assinatura;
   - integra com `assinador.jar` em dois modos:
     - **local**: `java -jar ...`
     - **HTTP**: requisições para o serviço em execução.

2. **`assinador.jar` (Java)**
   - validação rigorosa de parâmetros;
   - simulação de assinatura e validação;
   - tratamento de erros com mensagens claras.

3. **CLI `simulador` (multiplataforma)**
   - iniciar/parar/status do simulador;
   - checagem de portas;
   - download dinâmico de `simulador.jar` via GitHub Releases.

4. **Provisionamento automático de JDK**
   - detectar/baixar/configurar automaticamente quando necessário;
   - funcionar em Windows/Linux/macOS.

5. **Releases com qualidade e segurança**
   - binários multiplataforma + checksums SHA256;
   - versionamento SemVer;
   - assinatura dos artefatos com Cosign (Sigstore) no CI/CD.
