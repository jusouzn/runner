# Backlog Inicial — Sistema Runner

Data base do planejamento: 2026-03-31  
Referência de requisitos: `kyriosdata/runner` (especificação em `especificacao.md`)

## Convenções
- Prioridade: P0 (agora) / P1 (próximo) / P2 (depois)
- Formato das tarefas: BDD (Dado / Quando / Então)
- “CLI assinatura” = aplicativo de linha de comando que invoca `assinador.jar`
- “Assinador” = `assinador.jar` (Java), com simulação e **validação rigorosa** de parâmetros
- “CLI simulador” = aplicativo de linha de comando para obter/iniciar/parar/checar `simulador.jar` (via releases)

---

## US-01 — Invocar assinador.jar via CLI (P0)
**Como** usuário do Sistema Runner  
**Quero** executar comandos de assinatura digital através da linha de comandos  
**Para que** eu possa invocar o `assinador.jar` sem conhecer detalhes de configuração Java, tanto para assinar quanto para validar

### Tarefas (priorizadas)
1. **(P0) Definir contrato de comandos do CLI (help, subcomandos e parâmetros)**
   - Dado que sou um usuário no terminal  
     Quando eu executo `assinatura --help`  
     Então vejo descrição do produto, comandos disponíveis e exemplos de uso
   - Dado que sou um usuário no terminal  
     Quando eu executo `assinatura assinar --help` e `assinatura validar --help`  
     Então vejo os parâmetros obrigatórios/opcionais e mensagens consistentes

2. **(P0) Implementar parsing/validação inicial de argumentos no CLI**
   - Dado que informo um comando inexistente  
     Quando executo o CLI  
     Então recebo erro claro e retorno (exit code) diferente de zero
   - Dado que omito parâmetro obrigatório do CLI  
     Quando executo o CLI  
     Então recebo mensagem indicando o parâmetro ausente e exemplo

3. **(P0) Implementar invocação local (modo CLI) do `assinador.jar`**
   - Dado que o `assinador.jar` está disponível localmente  
     Quando executo `assinatura assinar ...`  
     Então o CLI executa `java -jar assinador.jar ...` e imprime o resultado formatado
   - Dado que o `assinador.jar` retorna erro (stderr/exit code != 0)  
     Quando executo o comando  
     Então o CLI apresenta erro compreensível, preserva detalhes técnicos em modo verbose e encerra com exit code != 0

4. **(P0) Implementar modo servidor (HTTP) do `assinador.jar` e fallback/seleção no CLI**
   - Dado que não informei `--local`  
     Quando executo um comando de assinatura  
     Então o CLI tenta usar o modo servidor via HTTP
   - Dado que não existe instância do servidor rodando na porta padrão  
     Quando executo um comando sem `--local`  
     Então o CLI inicia o `assinador.jar` em modo servidor na porta padrão e usa essa instância
   - Dado que já existe instância do servidor rodando  
     Quando executo um comando sem `--local`  
     Então o CLI reutiliza a instância existente (sem reiniciar)

5. **(P1) Implementar comandos de controle do servidor `assinador.jar`**
   - Dado que o servidor está rodando na porta padrão (ou indicada)  
     Quando executo `assinatura servidor parar`  
     Então o servidor encerra e o CLI confirma a ação
   - Dado que solicito parada programada por inatividade  
     Quando executo `assinatura servidor parar --apos <minutos>`  
     Então o servidor é configurado para autoencerrar após o tempo sem interação

6. **(P1) Padronizar formatação de saída e códigos de retorno**
   - Dado que a operação foi bem-sucedida  
     Quando o CLI imprime o resultado  
     Então a saída é legível e adequada para uso humano e automação (ex.: inclui status e payload)
   - Dado que houve falha de validação/execução  
     Quando o CLI finaliza  
     Então retorna exit code != 0 e mensagem útil

---

## US-02 — Simular assinatura digital com validação de parâmetros (P0)
**Como** usuário do Sistema Runner  
**Quero** que o `assinador.jar` valide rigorosamente os parâmetros de entrada antes de simular a operação  
**Para que** eu receba feedback imediato sobre erros e apenas requisições bem formadas sejam aceitas

### Tarefas (priorizadas)
1. **(P0) Criar projeto Java do `assinador.jar` com estrutura mínima e build reproduzível**
   - Dado que o repositório foi clonado  
     Quando executo o build (ex.: Maven/Gradle)  
     Então gero um `assinador.jar` executável de forma consistente

2. **(P0) Implementar comando “assinar” com validação rigorosa de parâmetros**
   - Dado que todos os parâmetros são válidos  
     Quando executo o comando de assinatura  
     Então recebo resposta simulada pré-construída (determinística)
   - Dado que existe parâmetro inválido (formato/ausente/incompatível)  
     Quando executo o comando  
     Então recebo mensagem de erro clara, apontando o campo inválido e motivo

3. **(P0) Implementar comando “validar” com validação rigorosa e resposta simulada**
   - Dado que os parâmetros são válidos  
     Quando executo validação  
     Então recebo “válido”/“inválido” conforme regra simulada definida
   - Dado que os parâmetros são inválidos  
     Quando executo validação  
     Então recebo erro estruturado e mensagem clara

4. **(P1) Implementar API HTTP do `assinador.jar` (modo servidor)**
   - Dado que inicio o `assinador.jar` em modo servidor  
     Quando envio uma requisição HTTP de assinatura válida  
     Então recebo a mesma resposta simulada que no modo local
   - Dado que envio uma requisição inválida  
     Quando o servidor processa  
     Então retorna status HTTP apropriado (4xx) e corpo com erro explicativo

5. **(P2) Preparar “stubs”/interface para PKCS#11 (sem implementar criptografia real)**
   - Dado que informo parâmetros relacionados a dispositivo criptográfico (quando aplicável)  
     Quando executo o comando  
     Então o sistema valida coerência/forma e segue em modo simulado, sem operação criptográfica real

---

## US-03 — Gerenciar ciclo de vida do Simulador do HubSaúde (P1)
**Como** usuário do Sistema Runner  
**Quero** iniciar, parar e monitorar o `simulador.jar` via CLI  
**Para que** eu possa gerenciar o simulador sem conhecer comandos Java

### Tarefas (priorizadas)
1. **(P1) Implementar “obter versão mais recente” do simulador via GitHub Releases**
   - Dado que não tenho `simulador.jar` localmente  
     Quando executo `simulador obter`  
     Então o CLI baixa a versão mais recente do repositório de releases definido e salva localmente
   - Dado que já tenho a versão mais recente local  
     Quando executo `simulador obter`  
     Então o CLI não baixa novamente e informa que está atualizado

2. **(P1) Verificar disponibilidade de portas antes de iniciar**
   - Dado que as portas necessárias estão livres  
     Quando executo `simulador iniciar`  
     Então o simulador inicia com sucesso
   - Dado que uma porta necessária está em uso  
     Quando executo `simulador iniciar`  
     Então o CLI informa a porta em conflito e não inicia o processo

3. **(P1) Implementar start/stop/status do simulador**
   - Dado que o simulador não está rodando  
     Quando executo `simulador status`  
     Então vejo “não está em execução”
   - Dado que o simulador está rodando  
     Quando executo `simulador status`  
     Então vejo PID/porta(s) e estado
   - Dado que o simulador está rodando  
     Quando executo `simulador parar`  
     Então o processo é encerrado e o CLI confirma

---

## US-04 — Provisionar JDK automaticamente (P1)
**Como** usuário do Sistema Runner  
**Quero** que o sistema baixe e configure o JDK quando não estiver disponível  
**Para que** eu não precise instalar/configurar Java manualmente

### Tarefas (priorizadas)
1. **(P1) Detectar Java/JDK existente e versão mínima exigida**
   - Dado que existe Java instalado compatível  
     Quando executo `assinatura` / `simulador`  
     Então o sistema usa o Java do sistema sem baixar outro
   - Dado que não existe Java (ou versão incompatível)  
     Quando executo `assinatura` / `simulador`  
     Então o sistema detecta ausência/incompatibilidade e sugere/procede com provisionamento

2. **(P1) Implementar download do JDK por plataforma (Windows/Linux/macOS)**
   - Dado que estou no Windows/Linux/macOS e não tenho JDK  
     Quando o sistema baixa o JDK  
     Então baixa o pacote correto para a plataforma e valida integridade (quando possível)

3. **(P2) Cache local e reutilização do JDK provisionado**
   - Dado que já existe JDK provisionado localmente  
     Quando executo comandos novamente  
     Então o sistema reutiliza o JDK baixado e não rebaixa

---

## US-05 — Disponibilizar binários multiplataforma (P2)
**Como** usuário do Sistema Runner  
**Quero** baixar uma versão pré-compilada do CLI para minha plataforma  
**Para que** eu use o sistema sem compilar

### Tarefas (priorizadas)
1. **(P2) Pipeline de build para Windows/Linux/macOS**
   - Dado que há commit na branch principal (ou tag)  
     Quando o pipeline roda  
     Então gera binários para as 3 plataformas suportadas (amd64) e armazena como artefatos

2. **(P2) Publicar release com SemVer**
   - Dado que uma tag SemVer foi criada  
     Quando o pipeline executa  
     Então cria release no GitHub e publica os binários

3. **(P2) Gerar e publicar checksums SHA256**
   - Dado que os binários foram gerados  
     Quando o pipeline finaliza  
     Então publica checksums SHA256 junto da release

4. **(P2) Assinar artefatos com Cosign (OIDC) e anexar `.sig` e `.pem`**
   - Dado que a release está sendo criada  
     Quando o pipeline assina os artefatos  
     Então cada artefato tem `.sig` e `.pem` publicados e a assinatura fica verificável

---

## Itens técnicos transversais (P0 → P2)
1. **(P0) Estrutura base do repositório e padrões**
   - Dado que o repo foi inicializado  
     Quando alguém novo clona  
     Então consegue buildar e rodar com instruções mínimas (README, comandos, dependências)

2. **(P0) Testes (unitários + integração)**
   - Dado que um parâmetro inválido é fornecido  
     Quando executo testes  
     Então a suíte cobre os principais cenários de erro e sucesso (CLI e jar)

3. **(P1) Observabilidade básica**
   - Dado que uma execução falha  
     Quando rodo em modo `--verbose`  
     Então consigo ver comando Java/URL chamada, erros originais e contexto

4. **(P2) Empacotamento e experiência do usuário**
   - Dado que sou usuário final  
     Quando baixo e executo o binário  
     Então consigo usar com mensagens claras, exemplos e documentação de uso
