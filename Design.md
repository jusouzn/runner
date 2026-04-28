# Sistema Runner - Design Detalhado (Marco 2)

Este documento registra a arquitetura do Sistema Runner utilizando o modelo C4, expandido com definições de fluxo de processos, contratos de erro e estratégia de logging estabelecidos para a primeira iteração.

## 1. Diagrama de Contexto (Nível 1)

**Atores e Sistemas Externos:**

| Elemento | Tipo | Descrição |
|----------|------|-----------|
| Usuário | Ator | Pessoa que interage com o sistema via linha de comandos. |
| Dispositivo de Assinatura Digital | Sistema Externo | Hardware criptográfico (token USB, smart card) que armazena certificados e executa operações de assinatura. |
| Simulador do HubSaúde | Sistema Externo | Aplicação Web gerida pelo CLI e que responde a requisições de terceiros. |

## 2. Diagrama de Contêineres (Nível 2) e Decisões Tecnológicas

O CLI será implementado em **Go** pela sua capacidade de cross-compiling nativo, enquanto a lógica de validação de criptografia operará via **Java** e a biblioteca **SunPKCS11**. O estado da aplicação (PID, portas em uso, versão do JDK) será armazenado localmente na pasta `.hubsaude` do usuário.

**Comunicação entre contêineres:**

| Origem | Destino | Protocolo | Descrição |
|--------|---------|-----------|-----------|
| Usuário | assinador  | CLI | Comandos de assinatura (criar, validar) digitados no terminal. |
| Usuário | simulador | CLI | Comandos de gerenciamento do simulador. |
| assinador | assinador.jar | CLI / HTTP | Invocação direta ou requisição HTTP (conforme modo de execução). |
| assinador.jar | Dispositivo Criptográfico | PKCS#11 | Interface padrão para comunicação com tokens e smart cards. |
| simulador | Simulador do HubSaúde | HTTP | Invoca e monitora o ciclo de vida do simulador. |

## 3. Dinâmica de Fluxos (Modo Local vs HTTP)

A aplicação `assinatura` comunica-se com o `assinador.jar` de duas formas lógicas. 

### 3.1. Invocação Direta (Modo Local / Cold Start)
Ideal para scripts isolados.
1. O usuário executa `assinatura criar --local ...`.
2. O CLI CLI valida superficialmente os parâmetros.
3. O CLI executa o processo bloqueante: `~/.hubsaude/jdk/bin/java -jar assinador.jar --modo=local [parametros]`.
4. O Java inicia (cold start), processa a assinatura simulada, imprime o resultado via `STDOUT` (preferencialmente em JSON) e encerra o processo.
5. O CLI captura o `STDOUT`, formata e apresenta ao usuário.

### 3.2. Invocação via Servidor (Modo HTTP / Warm Start)
Ideal para múltiplas requisições sequenciais.
1. O usuário executa `assinatura criar ...`.
2. O CLI verifica localmente se o `assinador.jar` já está rodando (via arquivo de PID/Status em `.hubsaude`). 
   * *Se não estiver rodando:* O CLI provisiona o Java e executa a aplicação em background (`java -jar assinador.jar --modo=server --porta=8080`) e guarda o PID.
3. O CLI dispara uma requisição `POST http://localhost:8080/sign` enviando os parâmetros em JSON.
4. O Java (Spring Boot ou HttpServer nativo) valida os parâmetros e retorna o payload com status `200 OK` (sucesso) ou `400 Bad Request` (erro).
5. O CLI recebe o JSON, formata e apresenta a saída.

## 4. Contrato de Comunicação e Tratamento de Erros

Para garantir que os erros sejam "propagados de forma estruturada" e "apresentados ao usuário de forma clara", todo retorno do `assinador.jar` (seja via stdout ou resposta HTTP) deverá seguir o seguinte esquema estruturado (JSON):

### Formato de Sucesso
```json
{
  "status": "success",
  "operation": "sign",
  "data": {
    "signature": "MIIB...[string base64 simulada]...",
    "algorithm": "SHA256withRSA"
  }
}

```

### Formato de Erro
```json
{
  "status": "error",
  "error_code": "ERR_INVALID_PARAM",
  "message": "A chave privada fornecida é inválida ou o dispositivo não foi encontrado.",
  "details": ["O parâmetro --pin possui tamanho incorreto (deve ter no mínimo 4 dígitos)"]
}
