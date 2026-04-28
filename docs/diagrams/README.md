# Design e Arquitetura - Sistema Runner

Este projeto utiliza o [Modelo C4](https://c4model.com/) para documentação da arquitetura. Os diagramas são escritos em código utilizando [PlantUML](https://plantuml.com/) e ficam armazenados na pasta `diagramas/`.

## Níveis de Arquitetura

1. **Contexto (`contexto.puml`)**: Visão de alto nível que mostra o Sistema Runner e suas interações com o usuário, dispositivos criptográficos e a aplicação HubSaúde.
2. **Contêineres (`conteineres.puml`)**: Visão detalhada da fronteira do sistema, evidenciando as aplicações CLI (`assinatura` e `simulador` construídas em Go) e o núcleo em Java (`assinador.jar`).

## Como gerar as imagens (.svg)

Os scripts de geração automatizam a conversão dos arquivos `.puml` em imagens `.svg`. O script fará o download automático do `plantuml.jar` (caso não esteja disponível localmente) e processará os arquivos. As imagens serão salvas automaticamente dentro de um subdiretório chamado `imagens/` de onde o arquivo `.puml` estiver localizado.

### No Linux / macOS
Utilize o script bash disponibilizado na raiz do projeto:

```bash
# Dar permissão de execução (necessário apenas na primeira vez)
chmod +x geraimagens.sh

# Executar a geração padrão (processa apenas os arquivos modificados recentemente)
./geraimagens.sh
