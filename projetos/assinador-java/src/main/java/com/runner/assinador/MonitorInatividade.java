package com.runner.assinador;

import java.util.function.LongSupplier;

/**
 * Monitora a inatividade do servidor e dispara um encerramento quando a janela
 * configurada é ultrapassada sem novas requisições.
 *
 * <p>A lógica de decisão ({@link #inativo(long)}) é separada do laço de espera
 * para permitir teste rápido do timer sem depender do tempo de parede.
 */
final class MonitorInatividade {

    private final long janelaMs;
    private final long intervaloMs;
    private final LongSupplier relogio;
    private final LongSupplier ultimaAtividade;
    private final Runnable aoExpirar;

    /**
     * @param janelaMs        janela de inatividade tolerada, em milissegundos
     * @param intervaloMs     intervalo entre verificações, em milissegundos
     * @param relogio         fonte do "agora" (ex.: {@code System::currentTimeMillis})
     * @param ultimaAtividade fonte do instante da última atividade registrada
     * @param aoExpirar       ação executada quando a janela é ultrapassada
     */
    MonitorInatividade(long janelaMs, long intervaloMs, LongSupplier relogio,
                       LongSupplier ultimaAtividade, Runnable aoExpirar) {
        this.janelaMs = janelaMs;
        this.intervaloMs = intervaloMs;
        this.relogio = relogio;
        this.ultimaAtividade = ultimaAtividade;
        this.aoExpirar = aoExpirar;
    }

    /** Decisão pura: o servidor está inativo há tempo suficiente para encerrar? */
    boolean inativo(long agora) {
        return agora - ultimaAtividade.getAsLong() >= janelaMs;
    }

    /** Inicia o monitor em uma thread virtual; encerra ao expirar ou ser interrompido. */
    void iniciar() {
        Thread.ofVirtual().start(this::laco);
    }

    private void laco() {
        while (!Thread.currentThread().isInterrupted()) {
            try {
                Thread.sleep(intervaloMs);
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                return;
            }
            if (inativo(relogio.getAsLong())) {
                aoExpirar.run();
                return;
            }
        }
    }
}
