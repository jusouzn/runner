package com.runner.assinador;

import java.security.PrivateKey;
import java.security.Provider;
import java.security.PublicKey;
import java.security.Signature;
import java.util.Base64;

/**
 * {@link SignatureService} respaldado por um token PKCS#11 real.
 *
 * <p>A operação de assinatura usa a chave privada residente no token (via o
 * {@link Provider} SunPKCS11), de modo que o material da chave nunca sai do
 * dispositivo — as chamadas {@code C_Sign} são executadas pelo módulo PKCS#11.
 * A validação confere a assinatura com a chave pública correspondente.
 *
 * <p>Mantém o mesmo contrato JSON do {@link FakeSignatureService}, permitindo
 * que o {@code assinador.jar} use indistintamente um ou outro backend.
 */
public class Pkcs11SignatureService implements SignatureService {

    private static final String ALGORITMO_DEFAULT = "SHA256withRSA";

    private final Provider provedor;
    private final PrivateKey chavePrivada;
    private final PublicKey chavePublica;

    public Pkcs11SignatureService(Provider provedor, PrivateKey chavePrivada, PublicKey chavePublica) {
        this.provedor = provedor;
        this.chavePrivada = chavePrivada;
        this.chavePublica = chavePublica;
    }

    @Override
    public String sign(String conteudo, String algoritmo, String pin) {
        if (pin == null || pin.length() < 4) {
            throw new IllegalArgumentException(
                    "O parâmetro --pin possui tamanho incorreto (deve ter no mínimo 4 dígitos)");
        }
        if (!pin.matches("\\d+")) {
            throw new IllegalArgumentException(
                    "O parâmetro --pin deve conter apenas dígitos numéricos");
        }
        if (conteudo == null || conteudo.isBlank()) {
            throw new IllegalArgumentException("O parâmetro --conteudo é obrigatório");
        }
        String alg = (algoritmo == null || algoritmo.isBlank()) ? ALGORITMO_DEFAULT : algoritmo;
        try {
            Signature s = Signature.getInstance(alg, provedor);
            s.initSign(chavePrivada);
            s.update(conteudo.getBytes(java.nio.charset.StandardCharsets.UTF_8));
            String assinatura = Base64.getEncoder().encodeToString(s.sign());
            return """
                    {
                      "status": "success",
                      "operation": "sign",
                      "data": {
                        "signature": "%s",
                        "algorithm": "%s"
                      }
                    }""".formatted(assinatura, alg);
        } catch (java.security.NoSuchAlgorithmException e) {
            throw new IllegalArgumentException("Algoritmo não suportado: " + alg, e);
        } catch (java.security.GeneralSecurityException e) {
            throw new IllegalStateException("Falha ao assinar via PKCS#11: " + e.getMessage(), e);
        }
    }

    @Override
    public String validate(String conteudo, String assinatura, String pin) {
        if (conteudo == null || conteudo.isBlank()) {
            throw new IllegalArgumentException("O parâmetro --conteudo é obrigatório");
        }
        if (assinatura == null || assinatura.isBlank()) {
            throw new IllegalArgumentException("O parâmetro --assinatura é obrigatório");
        }
        boolean valida;
        try {
            Signature v = Signature.getInstance(ALGORITMO_DEFAULT);
            v.initVerify(chavePublica);
            v.update(conteudo.getBytes(java.nio.charset.StandardCharsets.UTF_8));
            valida = v.verify(Base64.getDecoder().decode(assinatura));
        } catch (java.security.GeneralSecurityException | IllegalArgumentException e) {
            valida = false;
        }
        return """
                {
                  "status": "success",
                  "operation": "validate",
                  "data": {
                    "valid": %b,
                    "algorithm": "%s"
                  }
                }""".formatted(valida, ALGORITMO_DEFAULT);
    }
}
