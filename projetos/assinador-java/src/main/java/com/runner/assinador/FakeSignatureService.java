package com.runner.assinador;

import java.util.Set;

public class FakeSignatureService implements SignatureService {

    private static final String ASSINATURA_SIMULADA =
            "MIIB...[base64-simulado-SHA256withRSA]...AAAA==";

    private static final String ALGORITMO_DEFAULT = "SHA256withRSA";

    private static final Set<String> ALGORITMOS_ACEITOS =
            Set.of("SHA256withRSA", "SHA512withRSA", "SHA1withRSA");

    private static final int CONTEUDO_TAMANHO_MAXIMO = 10_000;

    @Override
    public String sign(String conteudo, String algoritmo, String pin) {
        validarPin(pin);
        validarConteudo(conteudo);
        String alg;
        if (algoritmo == null || algoritmo.isBlank()) {
            alg = ALGORITMO_DEFAULT;
        } else {
            validarAlgoritmo(algoritmo);
            alg = algoritmo;
        }
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
        if (!pin.matches("\\d+")) {
            throw new IllegalArgumentException(
                    "O parâmetro --pin deve conter apenas dígitos numéricos");
        }
    }

    private void validarConteudo(String conteudo) {
        if (conteudo == null || conteudo.isBlank()) {
            throw new IllegalArgumentException("O parâmetro --conteudo é obrigatório");
        }
        if (conteudo.length() > CONTEUDO_TAMANHO_MAXIMO) {
            throw new IllegalArgumentException(
                    "O parâmetro --conteudo excede o tamanho máximo permitido ("
                            + CONTEUDO_TAMANHO_MAXIMO + " caracteres)");
        }
    }

    private void validarAlgoritmo(String algoritmo) {
        if (!ALGORITMOS_ACEITOS.contains(algoritmo)) {
            throw new IllegalArgumentException(
                    "Algoritmo não suportado: " + algoritmo
                            + ". Algoritmos aceitos: SHA256withRSA, SHA512withRSA, SHA1withRSA");
        }
    }
}
