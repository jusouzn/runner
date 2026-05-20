package com.runner.assinador;

public class AssinadorMain {

    private static final SignatureService service = new FakeSignatureService();

    public static void main(String[] args) throws Exception {
        if (args.length == 0) {
            erro("ERR_INVALID_PARAM",
                    "Nenhum parâmetro fornecido. Use --modo=local|servidor --operacao=assinar|validar --pin=... --conteudo=...");
            System.exit(1);
        }

        String modo       = null;
        String pin        = null;
        String operacao   = null;
        String conteudo   = null;
        String assinatura = null;
        String algoritmo  = null;
        int    porta      = 8080;

        for (String arg : args) {
            if      (arg.startsWith("--modo="))        modo       = arg.substring(7);
            else if (arg.startsWith("--pin="))         pin        = arg.substring(6);
            else if (arg.startsWith("--operacao="))    operacao   = arg.substring(11);
            else if (arg.startsWith("--conteudo="))    conteudo   = arg.substring(11);
            else if (arg.startsWith("--assinatura="))  assinatura = arg.substring(13);
            else if (arg.startsWith("--algoritmo="))   algoritmo  = arg.substring(12);
            else if (arg.startsWith("--porta=")) {
                try {
                    porta = Integer.parseInt(arg.substring(8));
                } catch (NumberFormatException e) {
                    erro("ERR_INVALID_PARAM", "Valor inválido para --porta: " + arg.substring(8));
                    System.exit(1);
                }
            }
        }

        if (modo == null) {
            erro("ERR_INVALID_PARAM", "Parâmetro --modo é obrigatório (local ou servidor)");
            System.exit(1);
        }

        switch (modo) {
            case "servidor" -> new ServidorHttp(porta, service).iniciar();
            case "local"    -> runLocal(operacao, conteudo, assinatura, algoritmo, pin);
            default -> {
                erro("ERR_INVALID_PARAM",
                        "Valor inválido para --modo: '" + modo + "'. Use 'local' ou 'servidor'");
                System.exit(1);
            }
        }
    }

    private static void runLocal(String operacao, String conteudo,
                                  String assinatura, String algoritmo, String pin) {
        try {
            String resultado = switch (operacao != null ? operacao : "") {
                case "assinar" -> service.sign(conteudo, algoritmo, pin);
                case "validar" -> service.validate(conteudo, assinatura, pin);
                default -> {
                    erro("ERR_INVALID_PARAM", "Parâmetro --operacao deve ser 'assinar' ou 'validar'");
                    System.exit(1);
                    yield "";
                }
            };
            System.out.println(resultado);
        } catch (IllegalArgumentException e) {
            erro("ERR_INVALID_PARAM", e.getMessage());
            System.exit(1);
        }
    }

    private static void erro(String code, String message) {
        System.err.printf("""
                {
                  "status": "error",
                  "error_code": "%s",
                  "message": "%s",
                  "details": []
                }%n""", code, message);
    }
}
