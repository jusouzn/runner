package com.runner.assinador;

import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpServer;

import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.nio.charset.StandardCharsets;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.Executors;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class ServidorHttp {

    private final int porta;
    private final SignatureService service;

    public ServidorHttp(int porta, SignatureService service) {
        this.porta = porta;
        this.service = service;
    }

    public void iniciar() throws IOException, InterruptedException {
        HttpServer server = HttpServer.create(new InetSocketAddress(porta), 0);
        server.setExecutor(Executors.newVirtualThreadPerTaskExecutor());
        server.createContext("/health",   this::handleHealth);
        server.createContext("/sign",     this::handleSign);
        server.createContext("/validate", this::handleValidate);
        server.start();
        System.err.printf("{\"status\":\"running\",\"port\":%d}%n", porta);
        new CountDownLatch(1).await();
    }

    private void handleHealth(HttpExchange ex) throws IOException {
        if (!"GET".equalsIgnoreCase(ex.getRequestMethod())) {
            respond(ex, 405, "{\"error\":\"Method Not Allowed\"}");
            return;
        }
        respond(ex, 200, "{\"status\":\"ok\"}");
    }

    private void handleSign(HttpExchange ex) throws IOException {
        if (!"POST".equalsIgnoreCase(ex.getRequestMethod())) {
            respond(ex, 405, "{\"error\":\"Method Not Allowed\"}");
            return;
        }
        String body = readBody(ex);
        String conteudo  = extractField(body, "conteudo");
        String pin       = extractField(body, "pin");
        String algoritmo = extractField(body, "algoritmo");
        try {
            respond(ex, 200, service.sign(conteudo, algoritmo, pin));
        } catch (IllegalArgumentException e) {
            respond(ex, 400, errorJson("ERR_INVALID_PARAM", e.getMessage()));
        }
    }

    private void handleValidate(HttpExchange ex) throws IOException {
        if (!"POST".equalsIgnoreCase(ex.getRequestMethod())) {
            respond(ex, 405, "{\"error\":\"Method Not Allowed\"}");
            return;
        }
        String body      = readBody(ex);
        String conteudo  = extractField(body, "conteudo");
        String pin       = extractField(body, "pin");
        String assinatura = extractField(body, "assinatura");
        try {
            respond(ex, 200, service.validate(conteudo, assinatura, pin));
        } catch (IllegalArgumentException e) {
            respond(ex, 400, errorJson("ERR_INVALID_PARAM", e.getMessage()));
        }
    }

    private static String readBody(HttpExchange ex) throws IOException {
        try (InputStream is = ex.getRequestBody()) {
            return new String(is.readAllBytes(), StandardCharsets.UTF_8);
        }
    }

    private static void respond(HttpExchange ex, int status, String body) throws IOException {
        byte[] bytes = body.getBytes(StandardCharsets.UTF_8);
        ex.getResponseHeaders().set("Content-Type", "application/json; charset=utf-8");
        ex.sendResponseHeaders(status, bytes.length);
        try (OutputStream os = ex.getResponseBody()) {
            os.write(bytes);
        }
    }

    // Regex-based parser para JSON plano — evita dependências externas.
    private static String extractField(String json, String field) {
        Pattern p = Pattern.compile(
                "\"" + Pattern.quote(field) + "\"\\s*:\\s*\"((?:[^\"\\\\]|\\\\.)*)\"");
        Matcher m = p.matcher(json);
        return m.find() ? m.group(1) : null;
    }

    private static String errorJson(String code, String message) {
        String safe = message == null ? "" : message.replace("\"", "\\\"");
        return """
                {
                  "status": "error",
                  "error_code": "%s",
                  "message": "%s",
                  "details": []
                }""".formatted(code, safe);
    }
}
