# Escopo, Aceite e Riscos - Sistema Runner (Marco 2)

Este documento estabelece os limites do projeto, os critérios de aceitação atrelados às histórias de usuário e a análise inicial de riscos técnicos para o Sistema Runner, servindo como base para o congelamento da arquitetura da primeira iteração.

## 1. Escopo do Sistema

A definição clara do escopo garante o foco na entrega de valor estrutural e de integração.

### 1.1. O que ESTÁ no Escopo
* Desenvolvimento da aplicação `assinatura` (CLI multiplataforma).
* Desenvolvimento da aplicação `assinador.jar` (Java).
* Integração estruturada entre as duas aplicações (CLI e Java).
* Validação rigorosa de parâmetros pelo `assinador.jar`.
* Simulação de criação e validação de assinatura (`assinador.jar`).
* Tratamento unificado de erros dos parâmetros e exceções (`assinador.jar`).
* Desenvolvimento da aplicação `simulador` (CLI multiplataforma).
* Criação de testes automatizados (unidade e integração).
* Documentação técnica e de uso.

### 1.2. O que NÃO ESTÁ no Escopo
* Implementação real de assinatura digital criptográfica.
* Implementação real de validação de assinatura digital criptográfica.
* Integração com autoridades certificadoras (ACs).
* Armazenamento persistente de assinaturas geradas.
* Interface gráfica (GUI - Graphical User Interface).
* Fluxos de autenticação de usuários.
* Geração de certificados digitais.

## 2. Critérios de Aceite por User Story (US)

Conforme a especificação, as entregas serão validadas contra as seguintes histórias de usuário:

**US-01: Invocar assinador.jar via CLI**
* O CLI deve aceitar comandos para criação e validação de assinatura.
* O CLI deve invocar o `assinador.jar` com os parâmetros fornecidos.
* O CLI deve permitir a invocação direta do `assinador.jar` (modo local/CLI).
* O CLI deve permitir a invocação do `assinador.jar` via HTTP (modo servidor).
* O CLI deve exibir o resultado da operação de forma legível.
* O CLI deve iniciar o `assinador.jar` no modo servidor usando a porta padrão quando não orientado de forma diferente.
* O CLI deve detectar se instância do `assinador.jar` já está em execução no modo servidor e usar essa instância, se não orientado de forma diferente.
* O CLI deve fazer uso do `assinador.jar` no modo servidor quando não orientado para usar o modo local.
* O CLI deve interromper a execução do `assinador.jar` na porta padrão ou outra indicada.
* O CLI deve permitir a requisição de interrupção programada do `assinador.jar` após o número de minutos fornecidos sem interação.

**US-02: Simular assinatura digital com validação de parâmetros**
* O `assinador.jar` deve validar todos os parâmetros conforme especificações.
* O `assinador.jar` deve simular a criação de assinatura retornando resposta pré-construída quando parâmetros válidos.
* O `assinador.jar` deve simular validação de assinatura retornando resultado pré-determinado.
* O `assinador.jar` deve suportar interação com dispositivo criptográfico via interface PKCS#11.
* O `assinador.jar` deve retornar mensagens de erro claras quando parâmetros forem inválidos.

**US-03: Gerenciar Ciclo de Vida do Simulador do HubSaúde**
* O CLI deve permitir iniciar o Simulador.
* O CLI deve verificar se as portas necessárias para o Simulador estão disponíveis antes de iniciar.
* O CLI deve permitir parar o Simulador.
* O CLI deve exibir o status atual do Simulador (ou informar que não está em execução).
* O Simulador (`simulador.jar`) não faz parte do escopo de desenvolvimento deste sistema.
* O Simulador (`simulador.jar`) deve ser obtido dinamicamente pelo CLI, baixando a versão mais recente disponível no repositório da disciplina (GitHub Releases).
* O CLI não deve baixar o Simulador (`simulador.jar`) se a versão mais recente já estiver disponível localmente.

**US-04: Provisionar JDK Automaticamente**
* O sistema deve detectar se o JDK está presente na máquina (na versão exigida).
* O sistema deve baixar o JDK compatível quando ausente.
* O sistema deve disponibilizar o JDK baixado para uso próprio ou seja, pelo Assinador e Simulador.
* O download deve funcionar nas três plataformas (Windows, Linux, macOS).

**US-05: Disponibilizar binários multiplataforma**
* Disponibilizar binário para Windows (amd64).
* Disponibilizar binário para Linux (amd64).
* Disponibilizar binário para macOS (amd64).
* Distribuir via GitHub Releases.
* Incluir checksums SHA256 para verificação de integridade.
* Utilizar versionamento semântico (SemVer).

## 3. Definition of Done (DoD) - Definição de Pronto

Para que uma funcionalidade seja considerada "pronta" nesta e nas próximas iterações, os seguintes critérios devem ser satisfeitos:
1. Código compilando sem avisos (warnings) críticos.
2. Testes de unidade e de integração implementados e passando no pipeline de CI/CD.
3. Tratamento explícito de exceções padronizado conforme o Contrato de Erros documentado no Design.
4. Código validado por pelo menos um revisor (Pull Request review).
5. Binários gerados, com checksum SHA256 e assinados via **Cosign/Sigstore**.

## 4. Análise e Matriz de Riscos Iniciais

| Risco Técnico | Impacto | Estratégia de Mitigação |
|---|---|---|
| **Provisionamento Automático do JDK** | Alto (Bloqueia o uso do CLI) | Criar abstrações no CLI (em Go) para baixar o `.tar.gz`/`.zip` do Adoptium/Eclipse Temurin, descompactando-o num diretório isolado (`~/.hubsaude/jdk`). |
| **Gerenciamento de Processos (Start/Stop/PID)** | Alto (Processos zumbis ou portas presas) | O CLI deve manter um arquivo local de estado (ex: `~/.hubsaude/runner.db` ou `.json`). O CLI lerá esse banco para checar o status e o PID ativo. |
| **Assinatura de Artefatos com Cosign/Sigstore** | Médio (Pode travar a entrega contínua) | Implementar a assinatura como um step isolado e automatizado no pipeline do GitHub Actions (CI/CD) usando a identidade OIDC. |
| **Cross-compiling e Multiplataforma** | Baixo (Suporte nativo da linguagem) | A escolha da linguagem Go oferece *cross-compiling* nativo, mitigando a complexidade de compilar para Windows, Linux e macOS. |