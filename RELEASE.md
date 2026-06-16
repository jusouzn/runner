# Processo de Release

Este documento descreve como o Sistema Runner é versionado, empacotado, assinado e
verificado. O processo é **automatizado** pelo workflow
[.github/workflows/release.yml](.github/workflows/release.yml), disparado pela criação de
uma tag de versão.

## Sumário

1. [Política de versionamento (SemVer)](#1-política-de-versionamento-semver)
2. [Criação de uma release (tags)](#2-criação-de-uma-release-tags)
3. [Artefatos publicados](#3-artefatos-publicados)
4. [Verificação com `cosign verify-blob`](#4-verificação-com-cosign-verify-blob)

---

## 1. Política de versionamento (SemVer)

O projeto adota [Semantic Versioning 2.0.0](https://semver.org/lang/pt-BR/). Toda versão
tem o formato `MAJOR.MINOR.PATCH`, e a **tag git** recebe o prefixo `v`
(ex.: `v1.4.2`). Observação: em [release.yml](.github/workflows/release.yml) o filtro
`on.push.tags` é um **glob** (não regex) e, do jeito que está hoje, só cobre versões estáveis
com **um dígito por parte** (ex.: `v1.2.3`). Para suportar `v10.0.0`, ajuste o padrão no workflow.

Incremente:

| Parte | Quando | Exemplo |
|-------|--------|---------|
| **MAJOR** | Mudança incompatível na CLI ou no contrato HTTP/JSON (flags removidas/renomeadas, formato de resposta alterado de forma quebrante) | `1.x.x → 2.0.0` |
| **MINOR** | Funcionalidade nova mantendo retrocompatibilidade (novo comando, nova flag opcional) | `1.4.x → 1.5.0` |
| **PATCH** | Correção de bug retrocompatível, ajustes internos, documentação que afete o binário | `1.4.1 → 1.4.2` |

Regras adicionais:

- **`0.y.z` (desenvolvimento inicial):** a API é considerada instável; mudanças
  quebrantes podem ocorrer em incrementos de `MINOR`. A primeira release estável pública
  é a `1.0.0`.
- **Pré-lançamentos:** use sufixos SemVer, ex.: `v1.0.0-rc.1`, `v1.0.0-beta.1`.
  > Nota: o gatilho atual do workflow cobre apenas versões estáveis
  > (`vX.Y.Z`). Para publicar pré-lançamentos automaticamente, amplie o padrão de tags
  > em [release.yml](.github/workflows/release.yml).
- A versão é injetada no binário em tempo de build via `-ldflags` (variável
  `cmd.Version`), permitindo `assinatura --version` / `simulador --version` reportarem a
  versão correta.

## 2. Criação de uma release (tags)

Pré-requisitos: a `main` deve estar verde (CI passando) e conter tudo o que entra na release.

```bash
# 1. Garanta que está na main atualizada
git checkout main
git pull origin main

# 2. Crie uma tag anotada seguindo SemVer (prefixo "v")
git tag -a v1.0.0 -m "Release v1.0.0"

# 3. Publique a tag — isso dispara o workflow de release
git push origin v1.0.0
```

Ao receber a tag, o workflow [release.yml](.github/workflows/release.yml):

1. Extrai a versão da tag (remove o prefixo `v`).
2. Compila os CLIs `assinatura` e `simulador` para Linux, Windows e macOS (amd64).
3. Compila o `assinador.jar`.
4. Gera o arquivo de checksums `SHA256SUMS`.
5. **Assina** cada artefato com Cosign no modo *keyless* (OIDC do GitHub Actions),
   produzindo `.sig` (assinatura) e `.pem` (certificado) por artefato.
6. Cria a **GitHub Release** `vX.Y.Z` com notas geradas automaticamente e faz upload de
   todos os artefatos (`dist/*`).

> **Atenção:** tags são imutáveis por convenção. Se algo estiver errado, **não**
> reescreva a tag — crie uma nova versão PATCH (ex.: `v1.0.1`).

## 3. Artefatos publicados

Cada release contém, para cada binário/JAR:

| Arquivo | Descrição |
|---------|-----------|
| `assinatura-<versão>-<os>-amd64[.exe]` | CLI `assinatura` por plataforma |
| `simulador-<versão>-<os>-amd64[.exe]` | CLI `simulador` por plataforma |
| `assinador-<versão>.jar` | Servidor/validador Java |
| `<artefato>.sig` | Assinatura Cosign (keyless) do artefato |
| `<artefato>.pem` | Certificado efêmero usado na assinatura |
| `SHA256SUMS` | Checksums SHA256 de todos os artefatos (também assinado) |

## 4. Verificação com `cosign verify-blob`

As assinaturas são **keyless** (Sigstore + OIDC do GitHub Actions): não há chave pública
fixa. A confiança é estabelecida verificando que o artefato foi assinado pelo **workflow
de release deste repositório**, comprovado pelo certificado efêmero emitido pela Fulcio.

Pré-requisito: [Cosign](https://docs.sigstore.dev/cosign/installation/) instalado.

### 4.1. Verificar um binário

Baixe da release o artefato, seu `.sig` e seu `.pem`. Em seguida:

```bash
ARTIFACT="assinatura-1.0.0-linux-amd64"

cosign verify-blob \
  --certificate      "${ARTIFACT}.pem" \
  --signature        "${ARTIFACT}.sig" \
  --certificate-identity-regexp "^https://github.com/jusouzn/runner/.github/workflows/release.yml@refs/tags/v.*$" \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
  "${ARTIFACT}"
```

Saída esperada: `Verified OK`.

O que cada flag garante:

- `--certificate` / `--signature`: o certificado e a assinatura a conferir.
- `--certificate-identity-regexp`: a identidade (workflow) que assinou — fixa a origem
  no `release.yml` deste repositório, sob uma tag `v*`. Para travar em uma versão exata,
  troque por `--certificate-identity` com a URL completa terminando em
  `@refs/tags/v1.0.0`.
- `--certificate-oidc-issuer`: o emissor OIDC esperado (GitHub Actions).

### 4.2. Verificar checksums (e o `SHA256SUMS` assinado)

Primeiro verifique a assinatura do próprio manifesto de checksums:

```bash
cosign verify-blob \
  --certificate      "SHA256SUMS.pem" \
  --signature        "SHA256SUMS.sig" \
  --certificate-identity-regexp "^https://github.com/jusouzn/runner/.github/workflows/release.yml@refs/tags/v.*$" \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
  "SHA256SUMS"
```

Depois confira a integridade dos artefatos baixados contra o manifesto:

```bash
sha256sum --check SHA256SUMS --ignore-missing
```

Cada arquivo presente deve reportar `OK`.

### 4.3. Resultado

- ✅ **`Verified OK`** + checksums `OK`: o artefato é autêntico (assinado pelo nosso
  workflow) e íntegro (não foi adulterado).
- ❌ Qualquer falha: **não utilize o artefato** e reporte uma issue.
