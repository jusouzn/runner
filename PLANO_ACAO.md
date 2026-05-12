# Plano de Ação (2026-01) — Sistema Runner (Trabalho Prático) 
**Equipe:** Veronica Ribeiro Oliveira Palmeira E Júlia de Souza Nascimento  
**Data de início:** **17/03/2026**  
**Data de entrega final:** **30/06/2026**  
**Checkpoint fixos:** 17/03, 31/03, 28/04, 05/05, 12/05, 19/05, 16/06, 30/06


## 1) Objetivo do projeto (o que vamos entregar)
Construir o **Sistema Runner**, cujo objetivo é **facilitar o uso (via CLI) de aplicações Java** do ecossistema da disciplina, reduzindo a necessidade de conhecimento técnico sobre instalação/configuração de Java.

O trabalho tem relevância prática por estar alinhado ao contexto de interoperabilidade em saúde (SES-GO/UFG) e aos requisitos de integração propostos.

### Entregáveis finais (conforme enunciado)
Ao final da disciplina, o repositório e as releases devem conter:

1. **Aplicação `assinatura` (CLI multiplataforma)**
   - Implementação completa
   - Compatível com Windows, Linux e macOS
   - Integração com `assinador.jar` em **dois modos**:
     - execução direta (CLI → `java -jar ...`)
     - execução via HTTP (CLI → requisição HTTP para o serviço do assinador)

2. **Aplicação `assinador.jar` (Java)**
   - Validação rigorosa de parâmetros (conforme especificações FHIR aplicáveis no enunciado)
   - Simulação de criação e validação de assinatura
   - Tratamento de erros e mensagens claras

3. **Aplicação `simulador` (CLI multiplataforma)**
   - Gerenciar ciclo de vida do **Simulador do HubSaúde** (iniciar/parar/status)
   - Verificar portas antes de iniciar
   - **Obter dinamicamente** o `simulador.jar` via GitHub Releases (baixar somente se não existir versão mais recente local)

4. **Provisionamento automático de JDK**
   - Detectar se JDK exigido existe; baixar/configurar automaticamente caso ausente
   - Funcionar em Windows/Linux/macOS
   - JDK disponível para uso pelo `assinador` e pelo `simulador`

5. **Releases com binários multiplataforma**
   - Windows amd64, Linux amd64, macOS amd64
   - Publicar via **GitHub Releases**
   - Incluir **checksums SHA256**
   - Versionamento **SemVer**
   - (Importante) Publicar os artefatos exigidos (nomes exemplificados no enunciado)

6. **Qualidade e evidências**
   - Testes unitários + integração + cenários de erro
   - Testes de aceitação baseados nos critérios das user stories
   - Documentação de uso + documentação técnica de integração
   - Especificação/escopo/contexto organizados
   - Diagramas C4

7. **Segurança de supply chain: assinatura dos artefatos com Cosign (Sigstore)**
   - Assinatura automática no CI/CD durante a release
   - Para cada artefato publicar:
     - `<artefato>`
     - `<artefato>.sig`
     - `<artefato>.pem`
   - Evidência de verificação via `cosign verify-blob ...`

---

## 2) Requisitos funcionais (síntese operacional)
Este plano se guia pelas user stories do enunciado:

- **US‑01 — Invocar Assinador via CLI**
  - CLI suporta criação/validação de assinatura
  - Suporta modo local (invocação direta) e modo servidor (HTTP)
  - Resultado apresentado de forma legível

- **US‑02 — Assinador valida parâmetros e simula operações**
  - Validação rigorosa dos parâmetros
  - Simulação com respostas pré-construídas
  - Mensagens de erro claras
  - Suporte a interface PKCS#11 (token/smartcard) *no nível exigido pelo enunciado (simulado/integração conforme definido)*

- **US‑03 — Gerenciar Simulador do HubSaúde**
  - CLI inicia/para/status
  - Verifica portas
  - Baixa dinamicamente `simulador.jar` via releases (cache local)

- **US‑04 — Provisionar JDK**
  - Detectar versão e baixar quando necessário (3 plataformas)

- **US‑05 — Releases multiplataforma**
  - Binários + checksums + SemVer + distribuição via releases

---

## 3) Estratégia de execução (como vamos trabalhar)
### 3.1 Organização e fluxo de trabalho
- **Branch principal:** `main` sempre estável.
- Mudanças por **Pull Request** (PR), com revisão da dupla.
- Cada PR deve trazer:
  - descrição do que foi feito,
  - como validar,
  - referência ao requisito (US‑01…US‑05) quando aplicável.

### 3.2 Divisão de responsabilidades (sem silos)
Para evitar “uma pessoa só sabe uma parte”, vamos alternar liderança por marco:

- **Veronica (V):** coordenação do cronograma, integração ponta a ponta, revisão final de documentação/entregas.
- **Júlia (J):** foco em qualidade, testes, validação de critérios de aceitação, documentação técnica e revisão cruzada.

> Regra: em toda entrega relevante, **uma implementa/escreve** e a outra **revisa e valida executando**.

### 3.3 Definição de “pronto”
Uma funcionalidade só é considerada pronta se:
- existe evidência no repositório (código/documento/scripts),
- há instrução para executar/testar,
- e há validação (teste automatizado ou checklist manual reprodutível).

---

## 4) Cronograma com datas fixas e entregas objetivas

### Marco 1 — **17/03/2026 (hoje) | Kickoff e alinhamento**
**Objetivo:** iniciar com controle de escopo, organização e plano de execução.

**Entregas**
- `PLANO_DE_ACAO.md` (este arquivo) no repositório.

**Responsáveis**
- V: publicar este plano.
- J: publicar este plano.

---

### Marco 2 — **31/03/2026 | Design e arquitetura congelados (primeira versão)**
**Objetivo:** “travar” a arquitetura e interfaces antes de implementar pesado.

**Entregas**
  - Documento curto `ESCOPO_E_ACEITE.md` (ou seção no README) contendo:
  - o que entra e o que não entra (copiado/sintetizado do enunciado),
  - critérios de aceite por user story (US‑01…US‑05),
  - lista de riscos iniciais (JDK, releases, cosign, multiplataforma).
- `design.md` atualizado com:
  - arquitetura (componentes `assinatura`, `assinador.jar`, `simulador`),
  - fluxos dos dois modos (local e HTTP),
  - contrato de erro (formato/estrutura de mensagens) e logs.
- **Diagramas C4 (níveis necessários)** em `diagramas/` + instrução de geração (via `geraimagens.sh/.bat` se aplicável).
- Backlog inicial em `BACKLOG.md` com tarefas por US (priorizadas).

**Responsáveis**
- J: revisar o enunciado e transformar critérios em checklist inicial. liderar diagramas C4 e revisão de consistência (requisito ↔ design).
- V:  + estrutura inicial de entregáveis. consolidar decisões técnicas e backlog priorizado.

---

### Marco 3 — **28/04/2026 | Implementação base + “walking skeleton”**
**Objetivo:** ter um esqueleto executável ponta a ponta, mesmo que mínimo.

**Entregas (mínimo viável, mas real)**
- `assinatura` CLI executa comandos e:
  - chama `assinador.jar` **localmente** (modo direto) com parâmetros mínimos.
- `assinador.jar`:
  - valida parâmetros básicos (já com mensagens de erro claras),
  - responde a uma simulação simples de assinatura/validação.
- Primeiros testes:
  - unitários para validação de parâmetros,
  - integração mínima CLI → jar.

**Responsáveis**
- V: integração CLI → jar (modo local) e estrutura de comandos.
- J: suite de testes inicial + matriz de cenários de erro.

---

### Marco 4 — **05/05/2026 | Modo servidor (HTTP) e contrato de integração**
**Objetivo:** implementar o segundo modo de invocação (warm start) e documentar o contrato.

**Entregas**
- `assinador.jar` operando como **servidor HTTP** (start/stop básico) + endpoints necessários.
- `assinatura` suportando:
  - modo local (direto)
  - modo HTTP (servidor)
- Documento `INTEGRACAO.md`:
  - como iniciar o servidor,
  - endpoints e payloads,
  - exemplos de request/response,
  - mapeamento de erros.

**Responsáveis**
- J: validar contrato HTTP e escrever `INTEGRACAO.md`.
- V: implementar e estabilizar modo HTTP.

---

### Marco 5 — **12/05/2026 | Simulador CLI + download dinâmico via Releases**
**Objetivo:** cumprir US‑03 (gerenciamento do simulador e download inteligente).

**Entregas**
- CLI `simulador`:
  - start / stop / status
  - checagem de portas antes de iniciar
- Implementação do **download do `simulador.jar`** via GitHub Releases:
  - baixa “latest” quando necessário,
  - não baixa se versão mais recente já está local.
- Documentação:
  - `USO_SIMULADOR.md` (ou seção no README) com exemplos.

**Responsáveis**
- V: lógica de download/cache + integração com releases.
- J: checagem de portas + validação com cenários (porta ocupada, sem rede, jar ausente).

---

### Marco 6 — **19/05/2026 | Provisionamento automático de JDK (US‑04)**
**Objetivo:** remover dependência do usuário instalar/configurar Java manualmente.

**Entregas**
- Detecção de JDK exigido + download/configuração automática:
  - Windows, Linux, macOS (amd64)
- Documentação `JDK.md`:
  - onde instala (diretório local do app),
  - como limpa cache,
  - como diagnosticar problemas.

**Responsáveis**
- J: especificar e testar matriz de ambientes (por OS) + documentação.
- V: implementação do provisionamento + integração com `assinador` e `simulador`.

---

### Marco 7 — **16/06/2026 | Pipeline de Releases + Cosign + checksums (US‑05 e segurança)**
**Objetivo:** fechar distribuição profissional e rastreável (releases prontas).

**Entregas**
- CI/CD produzindo artefatos multiplataforma:
  - Windows amd64
  - Linux amd64
  - macOS amd64
- Release com:
  - binários
  - `SHA256SUMS` (ou equivalente) com checksums
  - assinatura Cosign para cada artefato:
    - `<artefato>`
    - `<artefato>.sig`
    - `<artefato>.pem`
- Documento `RELEASE.md`:
  - como gerar release,
  - como verificar com `cosign verify-blob`,
  - convenção de versionamento (SemVer) adotada no projeto.

**Responsáveis**
- V: pipeline e empacotamento dos binários + checksums.
- J: Cosign e verificação ponta a ponta + documentação de verificação.

---

### Entrega final — **30/06/2026 | Finalização e validação completa**
**Objetivo:** consolidar tudo com evidência de qualidade, documentação e reprodutibilidade.

**Entregas finais**
- Repositório organizado, com:
  - documentação de uso (manual do usuário do `assinatura`),
  - documentação técnica de integração,
  - guias de instalação,
  - diagramas C4 atualizados,
  - rastreabilidade requisitos → implementação/testes.
- Testes e evidências:
  - unitários e integração (e/ou aceitação automatizada quando aplicável),
  - checklist final de critérios de aceitação por US (executado e marcado).
- Release final estável (SemVer) com artefatos e assinaturas (Cosign).

**Responsáveis**
- V + J: revisão final conjunta (auditável) + “freeze” do repositório.

---

## 5) Checklists (para garantir nota e evitar surpresas)

### 5.1 Checklist por User Story (aceite)
- [ ] US‑01: CLI cria/valida assinatura e suporta modo local e HTTP
- [ ] US‑02: `assinador.jar` valida parâmetros rigorosamente + simula operações + erros claros
- [ ] US‑03: CLI gerencia simulador + verifica portas + baixa jar via Releases (com cache)
- [ ] US‑04: JDK automático (detectar/baixar/configurar) nas 3 plataformas
- [ ] US‑05: releases com binários 3 plataformas + SHA256 + SemVer

### 5.2 Checklist de segurança (Cosign / supply chain)
- [ ] Pipeline assina automaticamente cada artefato na release
- [ ] Para cada artefato existe `.sig` e `.pem`
- [ ] Existe documentação de verificação com `cosign verify-blob`
- [ ] Evidência de verificação bem-sucedida (log/output ou instrução reproduzível)

### 5.3 Checklist de documentação mínima
- [ ] README com “Como instalar / Como executar / Exemplos”
- [ ] Manual do usuário do CLI (`assinatura`)
- [ ] Documentação técnica da integração (local e HTTP)
- [ ] Guia do provisionamento de JDK
- [ ] Guia de releases (SemVer + checksums + cosign)
