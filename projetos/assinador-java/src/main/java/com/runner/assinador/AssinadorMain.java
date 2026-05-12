package com.runner.assinador;

public class AssinadorMain {

    private static final SignatureService service = new FakeSignatureService();

    public static void main(String[] args) {
        if (args.length == 0) {
            erro("ERR_INVALID_PARAM", "Nenhum parâmetro fornecido. Use --modo=local --operacao=assinar|validar --pin=... --conteudo=...");
            System.exit(1);
        }

        String modo = null;
        String pin = null;
        String operacao = null;
        String conteudo = null;
        String assinatura = null;
        String algoritmo = null;

        for (String arg : args) {
            if (arg.startsWith("--modo="))        modo      = arg.substring(7);
            else if (arg.startsWith("--pin="))    pin       = arg.substring(6);
            else if (arg.startsWith("--operacao=")) operacao = arg.substring(11);
            else if (arg.startsWith("--conteudo=")) conteudo = arg.substring(11);
            else if (arg.startsWith("--assinatura=")) assinatura = arg.substring(13);
            else if (arg.startsWith("--algoritmo="))  algoritmo  = arg.substring(12);
        }

        if (!"local".equals(modo)) {
            erro("ERR_INVALID_PARAM", "Parâmetro --modo=local é obrigatório");
            System.exit(1);
        }

        try {
            String resultado = switch (operacao != null ? operacao : "") {
                case "assinar"  -> service.sign(conteudo, algoritmo, pin);
                case "validar"  -> service.validate(conteudo, assinatura, pin);
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
