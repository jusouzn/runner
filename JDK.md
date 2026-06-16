# Provisionamento Automático de JRE (US‑04)

Este documento descreve como o Sistema Runner detecta, baixa e configura
automaticamente um Java Runtime Environment (JRE) 21, de modo que o usuário
**não precise instalar ou configurar Java manualmente** para usar as aplicações
`assinatura` e `simulador`.

> Requisito atendido: **US‑04 — Provisionar JDK** (detectar versão e baixar
> quando necessário, nas 3 plataformas: Windows, Linux e macOS — amd64).

---

## 1. Visão geral

As CLIs precisam de um Java compatível (**versão mínima: 21**) para executar o
`assinador.jar` e o `simulador.jar`. Em vez de exigir uma instalação manual, o
Runner resolve o Java na seguinte ordem:

1. **JRE gerenciado** — `~/.hubsaude/jdk/bin/java`, provisionado pelo próprio Runner.
2. **Java do sistema** — primeiro `java` encontrado no `PATH`, desde que a versão
   seja `>= 21`.
3. **Provisionamento sob demanda** — se nada compatível for encontrado, o Runner orienta a executar `assinatura setup`, que baixa o **Eclipse Temurin JRE 21** para a plataforma atual.

Se um `java` for encontrado no `PATH` mas com versão incompatível (< 21), o Runner
**não** o utiliza e orienta o usuário a executar `assinatura setup`.

---

## 2. Onde o JRE é instalado

| Item | Caminho |
| --- | --- |
| Diretório do JRE gerenciado | `~/.hubsaude/jdk/` |
| Binário Java (Linux/macOS) | `~/.hubsaude/jdk/bin/java` |
| Binário Java (Windows) | `%USERPROFILE%\.hubsaude\jdk\bin\java.exe` |

O diretório `~/.hubsaude/` é a raiz de estado do Runner (também guarda
`assinador.jar` e o estado do servidor). A instalação é **local ao usuário** e
não requer privilégios de administrador.

---

## 3. Como provisionar

```bash
# Detecta o Java disponível; baixa e instala o JRE 21 apenas se necessário
assinatura setup
```

Saídas possíveis:

```text
# Já havia Java 21+ disponível (sistema ou gerenciado)
Java já disponível (versão 21): /usr/lib/jvm/temurin-21/bin/java

# Provisionamento realizado
Java não encontrado ou incompatível. Iniciando provisionamento automático...
Baixando JRE 21 (Eclipse Temurin)...
Extraindo JRE...
JRE instalado com sucesso (versão 21): /home/usuario/.hubsaude/jdk/bin/java
```

O provisionamento é acionado explicitamente pelo comando `assinatura setup`.
Quando nenhum Java compatível é encontrado, comandos como `sign` e `validate` falham com uma mensagem orientando a executar o setup.

---

## 4. Fonte e plataformas suportadas

Os binários são obtidos da **API oficial do Adoptium** (Eclipse Temurin), sempre
a última build GA do Java 21:

| Plataforma | Arch | Formato | Endpoint Adoptium |
| --- | --- | --- | --- |
| Linux | amd64 (`x64`) | `.tar.gz` | `/v3/binary/latest/21/ga/linux/x64/jre/hotspot/normal/eclipse` |
| macOS | amd64 (`x64`) | `.tar.gz` | `/v3/binary/latest/21/ga/mac/x64/jre/hotspot/normal/eclipse` |
| Windows | amd64 (`x64`) | `.zip` | `/v3/binary/latest/21/ga/windows/x64/jre/hotspot/normal/eclipse` |

Durante a extração, o diretório de topo do pacote (ex.: `jdk-21.0.3+9/`) é
removido, de modo que o binário fique sempre em `~/.hubsaude/jdk/bin/java`
independentemente da versão exata baixada.

---

## 5. Cache (não baixar se já existe versão válida)

Antes de baixar, o Runner verifica se `~/.hubsaude/jdk/bin/java` já existe e
reporta versão `>= 21`. Em caso afirmativo, **o download é ignorado**:

```text
JRE 21 já provisionado em ~/.hubsaude/jdk — nenhum download necessário.
```

Isso garante execuções subsequentes rápidas e sem rede.

---

## 6. Limpar / reprovisionar

Para forçar um novo download (ex.: instalação corrompida), basta remover o
diretório gerenciado e rodar o setup novamente:

```bash
# Linux / macOS
rm -rf ~/.hubsaude/jdk
assinatura setup

# Windows (PowerShell)
Remove-Item -Recurse -Force "$env:USERPROFILE\.hubsaude\jdk"
assinatura setup
```

---

## 7. Diagnóstico de problemas

| Sintoma | Causa provável | Ação |
| --- | --- | --- |
| `java não encontrado — ... execute assinatura setup` | Nenhum Java no sistema | Rode `assinatura setup` |
| `java encontrado em <path>, mas versão incompatível (mínimo: 21)` | Java do sistema é < 21 | Rode `assinatura setup` (instala JRE gerenciado) |
| `download falhou com status HTTP <n>` | Sem rede ou Adoptium indisponível | Verifique conectividade e tente novamente |
| `JRE instalado mas não encontrado` | Extração incompleta | `rm -rf ~/.hubsaude/jdk` e rode `assinatura setup` |

Para inspecionar manualmente a versão instalada:

```bash
~/.hubsaude/jdk/bin/java -version
```

---

## 8. Segurança

- O download usa os endpoints HTTPS oficiais do Adoptium.
- A extração possui proteção contra **path traversal / zip‑slip**: qualquer
  entrada do pacote que tente escapar de `~/.hubsaude/jdk` é descartada.

---

## 9. Rastreabilidade

| Item | Localização |
| --- | --- |
| Detecção do Java | [`internal/jdk/finder.go`](projetos/assinatura-go/internal/jdk/finder.go) |
| Download e extração | [`internal/jdk/provisioner.go`](projetos/assinatura-go/internal/jdk/provisioner.go) |
| Comando `setup` | [`cmd/setup.go`](projetos/assinatura-go/cmd/setup.go) |
| Testes unitários | [`internal/jdk/finder_test.go`](projetos/assinatura-go/internal/jdk/finder_test.go), [`internal/jdk/provisioner_test.go`](projetos/assinatura-go/internal/jdk/provisioner_test.go) |

> O `simulador` reutiliza a mesma estratégia (pacote `internal/jdk/` em
> `projetos/simulador-go/`), compartilhando o JRE em `~/.hubsaude/jdk`.
