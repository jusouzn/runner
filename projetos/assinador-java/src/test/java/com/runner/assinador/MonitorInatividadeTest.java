package com.runner.assinador;

import org.junit.jupiter.api.Test;

import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicLong;

import static org.junit.jupiter.api.Assertions.*;

class MonitorInatividadeTest {

    @Test
    void inativoQuandoJanelaUltrapassada() {
        var ultima = new AtomicLong(1_000);
        var monitor = new MonitorInatividade(500, 10, () -> 0, ultima::get, () -> {});
        assertTrue(monitor.inativo(1_600), "1600 - 1000 = 600 >= 500 → inativo");
    }

    @Test
    void ativoQuandoDentroDaJanela() {
        var ultima = new AtomicLong(1_000);
        var monitor = new MonitorInatividade(500, 10, () -> 0, ultima::get, () -> {});
        assertFalse(monitor.inativo(1_300), "1300 - 1000 = 300 < 500 → ainda ativo");
    }

    @Test
    void timerDisparaEncerramentoAposJanelaSemAtividade() throws InterruptedException {
        var encerrou = new CountDownLatch(1);
        // última atividade fixa no passado: a janela já está vencida.
        long ultimaAtividade = System.currentTimeMillis() - 10_000;
        var monitor = new MonitorInatividade(
                50, 10, System::currentTimeMillis, () -> ultimaAtividade, encerrou::countDown);

        monitor.iniciar();

        assertTrue(encerrou.await(2, TimeUnit.SECONDS),
                "o monitor deveria ter disparado o encerramento por inatividade");
    }

    @Test
    void timerNaoDisparaEnquantoHaAtividade() throws InterruptedException {
        var encerrou = new CountDownLatch(1);
        var ultimaAtividade = new AtomicLong(System.currentTimeMillis());
        var monitor = new MonitorInatividade(
                200, 10, System::currentTimeMillis, ultimaAtividade::get, encerrou::countDown);

        monitor.iniciar();

        // Mantém atividade recente por um período maior que a janela.
        for (int i = 0; i < 8; i++) {
            ultimaAtividade.set(System.currentTimeMillis());
            Thread.sleep(40);
        }

        assertEquals(1, encerrou.getCount(),
                "não deveria encerrar enquanto há atividade dentro da janela");
    }
}
