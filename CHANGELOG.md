# Changelog

Todas as mudanças relevantes deste projeto são documentadas aqui.

O formato segue [Keep a Changelog](https://keepachangelog.com/pt-BR/1.1.0/)
e o projeto adota [Semantic Versioning](https://semver.org/lang/pt-BR/).
A versão corrente é declarada em [`release.json`](release.json), fonte única
verificada pelo workflow de release contra a tag.

## [Não lançado]

## [0.1.0] — 2026-06-24

Primeira versão marcada. Consolida os três componentes (CLIs `assinatura` e
`simulador` em Go e o `assinador.jar`), invocação local e via HTTP, provisionamento
de JDK e o simulador do HubSaúde. Artefatos multiplataforma (Linux/Windows/macOS
amd64) com `SHA256SUMS` e assinatura Cosign keyless, publicados via GitHub Releases.

### Adicionado
- CLI `assinatura` (Go): `assinar`/`validar` em modo local (`java -jar`) e via
  HTTP, controle do servidor (`servidor parar [--apos N]`) e `setup` de JRE (US-01/US-04).
- `assinador.jar` (Java): validação rigorosa de parâmetros, simulação de
  assinatura/validação e API HTTP com erros estruturados (US-02).
- Integração PKCS#11 real via SunPKCS11 + SoftHSM2, validada em job dedicado na CI (#15).
- CLI `simulador` (Go): `obter`/`iniciar`/`parar`/`status` com checagem de porta e
  download dinâmico do `simulador.jar` via GitHub Releases (US-03).
- Pipeline de release multiplataforma com `SHA256SUMS` e assinatura Cosign keyless,
  e documentação de verificação (`cosign verify-blob`) em [RELEASE.md](RELEASE.md) (US-05).
- CI multiplataforma: testes Go e Java em matriz Windows + Linux (#13).
- Lint Go (`golangci-lint`) e relatório de cobertura (Go e JaCoCo) na CI (#13).
- Testes de contrato CLI ↔ `assinador.jar` (subprocess e HTTP reais) e cenários
  negativos: JAR/JVM ausentes, conexão recusada, payload inválido, porta ocupada (#14).
- Teste do timer de auto-shutdown por inatividade (#14).
- `version` passa a reportar tag + commit curto (rastreabilidade) (#16).
- Manifesto `release.json` gerado por release e verificado contra a tag (#16).

### Alterado
- `assinatura servidor`: falha rápido e com mensagem clara quando a porta está
  ocupada, em vez de aguardar o timeout completo de readiness (#14).
