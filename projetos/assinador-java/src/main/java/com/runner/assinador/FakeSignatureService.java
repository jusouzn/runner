package com.runner.assinador;

public class FakeSignatureService implements SignatureService {

    private static final String ASSINATURA_SIMULADA =
            "MIIB...[base64-simulado-SHA256withRSA]...AAAA==";

    @Override
    public String sign(String conteudo, String algoritmo, String pin) {
        validarPin(pin);
        validarConteudo(conteudo);
        String alg = (algoritmo != null && !algoritmo.isBlank()) ? algoritmo : "SHA256withRSA";
        return """
                {
                  "status": "success",
                  "operation": "sign",
                  "data": {
                    "signature": "%s",
                    "algorithm": "%s"
                  }
                }""".formatted(ASSINATURA_SIMULADA, alg);
    }

    @Override
    public String validate(String conteudo, String assinatura, String pin) {
        validarPin(pin);
        validarConteudo(conteudo);
        if (assinatura == null || assinatura.isBlank()) {
            throw new IllegalArgumentException("O parâmetro --assinatura é obrigatório");
        }
        return """
                {
                  "status": "success",
                  "operation": "validate",
                  "data": {
                    "valid": true,
                    "algorithm": "SHA256withRSA"
                  }
                }""";
    }

    private void validarPin(String pin) {
        if (pin == null || pin.length() < 4) {
            throw new IllegalArgumentException(
                    "O parâmetro --pin possui tamanho incorreto (deve ter no mínimo 4 dígitos)");
        }
    }

    private void validarConteudo(String conteudo) {
        if (conteudo == null || conteudo.isBlank()) {
            throw new IllegalArgumentException("O parâmetro --conteudo é obrigatório");
        }
    }
}
