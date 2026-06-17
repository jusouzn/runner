package com.runner.assinador;

import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.security.KeyStore;
import java.security.Provider;
import java.security.Security;

/**
 * Integração com módulos PKCS#11 reais via {@code SunPKCS11} (incluído na JVM).
 *
 * <p>Permite carregar um provedor a partir de uma biblioteca PKCS#11 nativa
 * (ex.: SoftHSM2, {@code libsofthsm2.so}) e abrir o {@link KeyStore} do token
 * mediante PIN. Os testes de integração usam SoftHSM2 como token real para
 * comprovar resposta a chamadas PKCS#11 (ver ADR-0004).
 */
public final class Pkcs11Support {

    private Pkcs11Support() {
    }

    /**
     * Configura um provedor SunPKCS11 apontando para a biblioteca informada.
     *
     * @param nomeProvedor   nome lógico do provedor (ex.: "SoftHSM")
     * @param bibliotecaPath caminho da biblioteca PKCS#11 nativa (.so/.dll)
     * @param slotListIndex  índice do slot na lista (geralmente 0)
     * @return provedor configurado (ainda não registrado globalmente)
     */
    public static Provider carregarProvedor(String nomeProvedor, String bibliotecaPath, int slotListIndex)
            throws IOException {
        if (nomeProvedor == null || nomeProvedor.isBlank()) {
            throw new IllegalArgumentException("nomeProvedor é obrigatório");
        }
        if (bibliotecaPath == null || bibliotecaPath.isBlank()) {
            throw new IllegalArgumentException("bibliotecaPath é obrigatório");
        }
        Provider base = Security.getProvider("SunPKCS11");
        if (base == null) {
            throw new IllegalStateException("Provedor SunPKCS11 indisponível nesta JVM");
        }
        String config = """
                name = %s
                library = %s
                slotListIndex = %d
                attributes(generate, CKO_PRIVATE_KEY, *) = {
                  CKA_TOKEN = false
                  CKA_SIGN = true
                  CKA_PRIVATE = true
                }
                attributes(generate, CKO_PUBLIC_KEY, *) = {
                  CKA_TOKEN = false
                  CKA_VERIFY = true
                }
                """.formatted(nomeProvedor, bibliotecaPath, slotListIndex);

        Path tmp = Files.createTempFile("pkcs11-", ".cfg");
        tmp.toFile().deleteOnExit();
        Files.writeString(tmp, config, StandardCharsets.UTF_8);
        return base.configure(tmp.toString());
    }

    /**
     * Abre o {@link KeyStore} do token (tipo "PKCS11") fazendo login com o PIN,
     * exercendo {@code C_OpenSession}/{@code C_Login} no módulo.
     */
    public static KeyStore abrirKeyStore(Provider provedor, String pin) throws Exception {
        if (pin == null || pin.isBlank()) {
            throw new IllegalArgumentException(
                    "O parâmetro --pin é obrigatório para abrir o KeyStore PKCS#11");
        }
        KeyStore ks = KeyStore.getInstance("PKCS11", provedor);
        ks.load(null, pin.toCharArray());
        return ks;
    }
}
