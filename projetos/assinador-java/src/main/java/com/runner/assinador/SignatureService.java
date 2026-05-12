package com.runner.assinador;

public interface SignatureService {

    /**
     * Simula a criação de uma assinatura digital.
     *
     * @param conteudo  conteúdo a ser assinado (Base64 ou texto)
     * @param algoritmo algoritmo de assinatura (ex.: "SHA256withRSA")
     * @param pin       PIN do dispositivo criptográfico
     * @return JSON com o resultado da operação
     * @throws IllegalArgumentException se algum parâmetro for inválido
     */
    String sign(String conteudo, String algoritmo, String pin);

    /**
     * Simula a validação de uma assinatura digital.
     *
     * @param conteudo    conteúdo original (Base64 ou texto)
     * @param assinatura  assinatura a ser validada (Base64)
     * @param pin         PIN do dispositivo criptográfico
     * @return JSON com o resultado da validação
     * @throws IllegalArgumentException se algum parâmetro for inválido
     */
    String validate(String conteudo, String assinatura, String pin);
}
