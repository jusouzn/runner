# ADR-0004: Integração PKCS#11 via SunPKCS11 e SoftHSM2 nos testes

- **Status:** Aceito
- **Data:** 2026-06-17
- **Decisores:** Verônica Ribeiro Oliveira Palmeira, Júlia de Souza Nascimento

## Contexto

O `assinador.jar` deve, em produção, assinar usando uma chave residente em um
dispositivo criptográfico acessível por **PKCS#11** (smartcard/HSM), sem que o
material da chave deixe o dispositivo. Para os requisitos da disciplina é
necessário **comprovar, por testes de integração, a resposta a chamadas PKCS#11
reais** — não basta um dublê em memória.

Forças:

- A JVM já inclui o provedor **`SunPKCS11`**, que faz a ponte para qualquer
  módulo PKCS#11 nativo, dispensando dependências de terceiros.
- Precisamos de um **token real** no ambiente de CI, reprodutível e gratuito,
  sem hardware físico.
- A portabilidade dos **CLIs** (Windows + Linux) é coberta à parte; um token
  PKCS#11 em CI é naturalmente dependente de SO.

## Decisão

1. A integração PKCS#11 usa o provedor **`SunPKCS11`** da própria JVM,
   configurado em tempo de execução a partir do caminho da biblioteca nativa
   (`Pkcs11Support.carregarProvedor`).
2. Os **testes de integração** usam **SoftHSM2** como módulo PKCS#11 real,
   provisionado no job de CI (`pkcs11-integration`, Ubuntu). O teste gera um par
   de chaves no token, assina (`C_Sign`) e valida a assinatura.
3. Os testes são **condicionais** à variável `PKCS11_MODULE`
   (`@EnabledIfEnvironmentVariable`): ausentes ela, são ignorados — assim o build
   local e os runners Windows não quebram.
4. O backend padrão do `assinador.jar` permanece o `FakeSignatureService`; o
   `Pkcs11SignatureService` compartilha o mesmo contrato JSON e pode substituí-lo
   quando um token estiver disponível.

## Alternativas consideradas

- **Mock/dublê de PKCS#11 em memória:** não comprova chamadas reais ao módulo —
  reprovado pelo requisito (E5).
- **Outro módulo (OpenSC/pkcs11-proxy, Nitrokey HSM):** SoftHSM2 é o mais simples
  de instalar e inicializar em CI (`apt-get` + `softhsm2-util`), sem hardware.
- **Biblioteca de terceiros (BouncyCastle, IAIK):** adiciona dependência onde a
  stdlib (SunPKCS11) já resolve — contraria "dependências mínimas".

## Consequências

**Positivas**

- Comprovação real (em CI) de assinatura/validação via PKCS#11.
- Sem novas dependências de runtime: apenas a JVM e a lib nativa do token.
- Contrato JSON único entre os backends fake e PKCS#11.

**Negativas / mitigações**

- O teste de integração roda **apenas no Linux** em CI (SoftHSM2 é trivial lá).
  *Mitigação:* a portabilidade dos binários é validada nos demais jobs; a
  integração PKCS#11 é independente de SO no código (SunPKCS11) e poderia ser
  exercida no Windows com um módulo equivalente no futuro.
- Chaves geradas são **de sessão** (`CKA_TOKEN = false`) para não poluir o token
  entre execuções; suficiente para comprovar o fluxo.
