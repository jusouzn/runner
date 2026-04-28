// AssinadorMain.java
package com.runner.assinador;

public class AssinadorMain {

    public static void main(String[] args) {
        // Validação básica: verificar se os parâmetros mínimos foram passados
        if (args.length == 0) {
            imprimirErro("ERR_INVALID_PARAM", "Nenhum parametro fornecido. Esperado: --modo=local");
            System.exit(1); // Exit code diferente de zero indica falha
        }

        boolean modoLocal = false;
        boolean temPin = false;

        // Leitura simplificada dos argumentos
        for (String arg : args) {
            if (arg.equals("--modo=local")) modoLocal = true;
            if (arg.startsWith("--pin=")) temPin = true;
        }

        // Cenário de Erro: Se não for modo local ou faltar PIN
        if (!modoLocal || !temPin) {
            imprimirErro("ERR_INVALID_PARAM", "Parametros obrigatorios ausentes (--modo=local ou --pin).");
            System.exit(1);
        }

        // Cenário de Sucesso: Simulação simples de assinatura
        imprimirSucesso();
        System.exit(0); // Sucesso
    }

    // Método para imprimir o contrato de erro em JSON
    private static void imprimirErro(String code, String message) {
        String jsonErro = String.format("""
            {
              "status": "error",
              "error_code": "%s",
              "message": "%s",
              "details": []
            }
            """, code, message);
        System.err.println(jsonErro);
    }

    // Método para imprimir o contrato de sucesso em JSON
    private static void imprimirSucesso() {
        String jsonSucesso = """
            {
              "status": "success",
              "operation": "sign",
              "data": {
                "signature": "MIIB...[string base64 simulada]...",
                "algorithm": "SHA256withRSA"
              }
            }
            """;
        System.out.println(jsonSucesso);
    }
}
