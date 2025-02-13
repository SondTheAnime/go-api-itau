# Documentação da API

Esta pasta contém a documentação da API gerada com o Scalar.

## Estrutura

- `scalar.yaml`: Arquivo de configuração OpenAPI/Swagger
- `dist/`: Pasta contendo a documentação gerada
- `README.md`: Este arquivo

## Como usar

### Gerar documentação

Para gerar a documentação estática:

```bash
make docs
```

A documentação será gerada na pasta `dist/`.

### Servidor de desenvolvimento

Para iniciar o servidor de desenvolvimento com hot-reload:

```bash
make docs-serve
```

O servidor estará disponível em `http://localhost:8088`.

## Atualizando a documentação

1. Edite o arquivo `scalar.yaml` com as alterações desejadas
2. Execute `make docs` para gerar a nova versão
3. Commit as alterações

## Referências

- [Scalar Documentation](https://github.com/MarceloPetrucio/go-scalar-api-reference)
- [OpenAPI Specification](https://swagger.io/specification/) 