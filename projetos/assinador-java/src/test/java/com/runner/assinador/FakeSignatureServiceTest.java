package com.runner.assinador;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class FakeSignatureServiceTest {

    private final SignatureService service = new FakeSignatureService();

    @Test
    void signRetornaJsonDeSucesso() {
        String resultado = service.sign("conteudo-teste", "SHA256withRSA", "1234");
        assertTrue(resultado.contains("\"status\": \"success\""));
        assertTrue(resultado.contains("\"operation\": \"sign\""));
        assertTrue(resultado.contains("\"signature\""));
    }

    @Test
    void signUsaAlgoritmoInformado() {
        String resultado = service.sign("conteudo", "SHA512withRSA", "1234");
        assertTrue(resultado.contains("SHA512withRSA"));
    }

    @Test
    void signUsaAlgoritmoDefaultQuandoNulo() {
        String resultado = service.sign("conteudo", null, "1234");
        assertTrue(resultado.contains("SHA256withRSA"));
    }

    @Test
    void validateRetornaJsonDeSucesso() {
        String resultado = service.validate("conteudo-teste", "assinatura-fake", "1234");
        assertTrue(resultado.contains("\"status\": \"success\""));
        assertTrue(resultado.contains("\"valid\": true"));
    }

    @Test
    void signComPinCurtoLancaExcecao() {
        var ex = assertThrows(IllegalArgumentException.class,
                () -> service.sign("conteudo", "SHA256withRSA", "123"));
        assertTrue(ex.getMessage().contains("--pin"));
    }

    @Test
    void signComPinNuloLancaExcecao() {
        assertThrows(IllegalArgumentException.class,
                () -> service.sign("conteudo", "SHA256withRSA", null));
    }

    @Test
    void signComConteudoNuloLancaExcecao() {
        assertThrows(IllegalArgumentException.class,
                () -> service.sign(null, "SHA256withRSA", "1234"));
    }

    @Test
    void validateComAssinaturaAusenteELancaExcecao() {
        assertThrows(IllegalArgumentException.class,
                () -> service.validate("conteudo", null, "1234"));
    }
}
