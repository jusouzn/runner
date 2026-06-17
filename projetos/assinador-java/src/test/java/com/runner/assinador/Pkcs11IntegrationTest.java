package com.runner.assinador;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.condition.EnabledIfEnvironmentVariable;

import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.Provider;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Testes de integração que exercem chamadas PKCS#11 reais contra um token
 * SoftHSM2 (ver ADR-0004). Habilitados apenas quando a variável de ambiente
 * {@code PKCS11_MODULE} aponta para a biblioteca nativa — caso contrário são
 * ignorados (build local e runners sem SoftHSM2 não quebram).
 */
@EnabledIfEnvironmentVariable(named = "PKCS11_MODULE", matches = ".+")
class Pkcs11IntegrationTest {

    private static final Pattern ASSINATURA =
            Pattern.compile("\"signature\":\\s*\"([^\"]+)\"");

    private static String env(String nome, String padrao) {
        String v = System.getenv(nome);
        return (v == null || v.isBlank()) ? padrao : v;
    }

    @Test
    void assinaEValidaUsandoTokenPkcs11Real() throws Exception {
        String modulo = System.getenv("PKCS11_MODULE");
        String pin = env("PKCS11_PIN", "1234");
        int slotIndex = Integer.parseInt(env("PKCS11_SLOT_INDEX", "0"));

        // Configura o provedor e faz login no token (C_OpenSession / C_Login).
        Provider provedor = Pkcs11Support.carregarProvedor("SoftHSM-Test", modulo, slotIndex);
        assertNotNull(Pkcs11Support.abrirKeyStore(provedor, pin),
                "deveria abrir o KeyStore do token após login");

        // Gera o par de chaves DENTRO do token (C_GenerateKeyPair).
        KeyPairGenerator kpg = KeyPairGenerator.getInstance("RSA", provedor);
        kpg.initialize(2048);
        KeyPair par = kpg.generateKeyPair();

        var servico = new Pkcs11SignatureService(provedor, par.getPrivate(), par.getPublic());

        // Assina via token (C_Sign) e confere o JSON de sucesso.
        String jsonSign = servico.sign("documento-de-teste", "SHA256withRSA", pin);
        assertTrue(jsonSign.contains("\"status\": \"success\""), jsonSign);

        Matcher m = ASSINATURA.matcher(jsonSign);
        assertTrue(m.find(), "JSON deveria conter a assinatura: " + jsonSign);
        String assinaturaB64 = m.group(1);
        assertFalse(assinaturaB64.isBlank(), "assinatura não deveria ser vazia");

        // Validação correta → valid: true.
        String jsonOk = servico.validate("documento-de-teste", assinaturaB64, pin);
        assertTrue(jsonOk.contains("\"valid\": true"), jsonOk);

        // Conteúdo adulterado → valid: false.
        String jsonAdulterado = servico.validate("documento-ADULTERADO", assinaturaB64, pin);
        assertTrue(jsonAdulterado.contains("\"valid\": false"), jsonAdulterado);
    }
}
