// AssinadorMain.java
package com.runner.assinador;

import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpServer;

import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.io.PrintStream;
import java.net.InetSocketAddress;
import java.nio.charset.StandardCharsets;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.Executors;

/**
 * Entrada do {@code assinador.jar}. Suporta dois modos de operacao:
 * <ul>
 *   <li>{@code --modo=local}  : executa a operacao e imprime o resultado em STDOUT.</li>
 *   <li>{@code --modo=server} : sobe um servidor HTTP local (warm start).</li>
 * </ul>
 * O contrato de mensagens (sucesso/erro em JSON) esta descrito em
 * {@code Design.md} e {@code docs/INTEGRACAO.md}.
 */
public class AssinadorMain {

    private static final String VERSION = "0.2.0";
    private static final int PORTA_PADRAO = 8080;

    public static void main(String[] args) throws Exception {
        Map<String, String> p = parseArgs(args);

        if (p.containsKey("--version") || p.containsKey("-v")) {
            System.out.println(VERSION);
            return;
        }
        if (p.containsKey("--help") || p.containsKey("-h")) {
            imprimirAjuda();
            return;
        }

        String modo = p.getOrDefault("--modo", "local");
        switch (modo) {
            case "local":
                executarLocal(p);
                break;
            case "server":
                int porta = parsePorta(p.getOrDefault("--porta", String.valueOf(PORTA_PADRAO)));
                iniciarServidor(porta);
                break;
            default:
                imprimirErro(System.err, "ERR_INVALID_PARAM", "Modo invalido: " + modo);
                System.exit(1);
        }
    }

    // --------- Modo local ---------

    private static void executarLocal(Map<String, String> p) {
        String operacao = p.getOrDefault("--operacao", "assinar");
        String pin = p.get("--pin");
        if (pin == null || pin.length() < 4) {
            imprimirErro(System.err, "ERR_INVALID_PARAM",
                    "Parametro --pin ausente ou invalido (minimo 4 digitos).");
            System.exit(1);
        }
        switch (operacao) {
            case "assinar":
                System.out.println(sucessoAssinarJson());
                break;
            case "validar":
                System.out.println(sucessoValidarJson());
                break;
            default:
                imprimirErro(System.err, "ERR_INVALID_PARAM", "Operacao invalida: " + operacao);
                System.exit(1);
        }
    }

    // --------- Modo servidor (HTTP) ---------

    private static void iniciarServidor(int porta) throws IOException {
        HttpServer server = HttpServer.create(new InetSocketAddress("localhost", porta), 0);
        server.setExecutor(Executors.newFixedThreadPool(4));

        server.createContext("/health", ex -> responder(ex, 200,
                "{\"status\":\"success\",\"data\":{\"versao\":\"" + VERSION + "\"}}"));

        server.createContext("/sign",     ex -> handleOperacao(ex, "sign"));
        server.createContext("/assinar",  ex -> handleOperacao(ex, "sign"));
        server.createContext("/validate", ex -> handleOperacao(ex, "validate"));
        server.createContext("/validar",  ex -> handleOperacao(ex, "validate"));

        server.createContext("/shutdown", ex -> {
            responder(ex, 200, "{\"status\":\"success\",\"message\":\"servidor encerrando\"}");
            new Thread(() -> {
                try { Thread.sleep(200); } catch (InterruptedException ignored) {}
                server.stop(0);
                System.exit(0);
            }, "shutdown-hook").start();
        });

        System.out.println("[assinador] Servidor HTTP iniciado em http://localhost:" + porta);
        System.out.println("[assinador] Versao " + VERSION + " | endpoints: /health /sign /validate /shutdown");
        server.start();
    }

    private static void handleOperacao(HttpExchange ex, String op) throws IOException {
        try {
            if (!"POST".equalsIgnoreCase(ex.getRequestMethod())) {
                responderErro(ex, 405, "ERR_METHOD_NOT_ALLOWED", "Metodo nao permitido. Use POST.");
                return;
            }
            String body = lerBody(ex);
            String pin = extrairCampoTexto(body, "pin");
            if (pin == null || pin.length() < 4) {
                responderErro(ex, 400, "ERR_INVALID_PARAM",
                        "Campo 'pin' ausente ou invalido (minimo 4 digitos).");
                return;
            }
            if ("sign".equals(op)) {
                responder(ex, 200, sucessoAssinarJson());
            } else {
                responder(ex, 200, sucessoValidarJson());
            }
        } catch (RuntimeException e) {
            responderErro(ex, 500, "ERR_INTERNO",
                    e.getMessage() == null ? "erro interno" : e.getMessage());
        }
    }

    // --------- Utilidades ---------

    private static Map<String, String> parseArgs(String[] args) {
        Map<String, String> m = new HashMap<>();
        for (String a : args) {
            int idx = a.indexOf('=');
            if (idx > 0) {
                m.put(a.substring(0, idx), a.substring(idx + 1));
            } else {
                m.put(a, "true");
            }
        }
        return m;
    }

    private static int parsePorta(String s) {
        try {
            int p = Integer.parseInt(s);
            if (p < 1 || p > 65535) throw new NumberFormatException();
            return p;
        } catch (NumberFormatException e) {
            imprimirErro(System.err, "ERR_INVALID_PARAM", "Porta invalida: " + s);
            System.exit(1);
            return -1;
        }
    }

    private static String lerBody(HttpExchange ex) throws IOException {
        try (InputStream is = ex.getRequestBody()) {
            return new String(is.readAllBytes(), StandardCharsets.UTF_8);
        }
    }

    /**
     * Extrator JSON ingenuo no formato {@code "campo":"valor"}. Suficiente para
     * o walking skeleton; sera substituido por parser real em iteracoes futuras.
     */
    private static String extrairCampoTexto(String json, String campo) {
        if (json == null) return null;
        String chave = "\"" + campo + "\"";
        int i = json.indexOf(chave);
        if (i < 0) return null;
        int colon = json.indexOf(':', i + chave.length());
        if (colon < 0) return null;
        int aspa1 = json.indexOf('"', colon + 1);
        if (aspa1 < 0) return null;
        int aspa2 = json.indexOf('"', aspa1 + 1);
        if (aspa2 < 0) return null;
        return json.substring(aspa1 + 1, aspa2);
    }

    private static void responder(HttpExchange ex, int status, String body) throws IOException {
        byte[] bytes = body.getBytes(StandardCharsets.UTF_8);
        ex.getResponseHeaders().add("Content-Type", "application/json; charset=utf-8");
        ex.sendResponseHeaders(status, bytes.length);
        try (OutputStream os = ex.getResponseBody()) {
            os.write(bytes);
        }
    }

    private static void responderErro(HttpExchange ex, int status, String code, String msg) throws IOException {
        responder(ex, status, erroJson(code, msg));
    }

    private static String sucessoAssinarJson() {
        return "{\n" +
               "  \"status\": \"success\",\n" +
               "  \"operation\": \"sign\",\n" +
               "  \"data\": {\n" +
               "    \"signature\": \"MIIB...[base64 simulada]...\",\n" +
               "    \"algorithm\": \"SHA256withRSA\"\n" +
               "  }\n" +
               "}";
    }

    private static String sucessoValidarJson() {
        return "{\n" +
               "  \"status\": \"success\",\n" +
               "  \"operation\": \"validate\",\n" +
               "  \"data\": {\n" +
               "    \"valid\": true,\n" +
               "    \"algorithm\": \"SHA256withRSA\"\n" +
               "  }\n" +
               "}";
    }

    private static String erroJson(String code, String msg) {
        String safe = msg == null ? "" : msg.replace("\\", "\\\\").replace("\"", "\\\"");
        return "{\n" +
               "  \"status\": \"error\",\n" +
               "  \"error_code\": \"" + code + "\",\n" +
               "  \"message\": \"" + safe + "\",\n" +
               "  \"details\": []\n" +
               "}";
    }

    private static void imprimirErro(PrintStream out, String code, String msg) {
        out.println(erroJson(code, msg));
    }

    private static void imprimirAjuda() {
        System.out.println("assinador.jar - validacao e simulacao de assinatura");
        System.out.println();
        System.out.println("Uso:");
        System.out.println("  java -jar assinador.jar --modo=local --operacao=assinar|validar --pin=NNNN");
        System.out.println("  java -jar assinador.jar --modo=server [--porta=8080]");
        System.out.println();
        System.out.println("Endpoints (modo server):");
        System.out.println("  GET  /health");
        System.out.println("  POST /sign      (alias /assinar)   body: {\"pin\":\"NNNN\"}");
        System.out.println("  POST /validate  (alias /validar)   body: {\"pin\":\"NNNN\"}");
        System.out.println("  POST /shutdown");
    }
}
