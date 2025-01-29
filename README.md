# API de Usuários com Documentação Scalar

Esta é uma API simples de gerenciamento de usuários com documentação utilizando Scalar.

## Pré-requisitos

- Go 1.19 ou superior

## Como executar

1. Inicialize o módulo e instale as dependências:
```bash
# Inicializar o módulo
go mod init api-usuarios

# Instalar dependências
go get github.com/go-chi/chi/v5
go get github.com/MarceloPetrucio/go-scalar-api-reference
```

2. Inicie a API:
```bash
go run api.go
```

3. Acesse a documentação da API:
- Abra seu navegador e acesse: http://localhost:8080/docs

## Endpoints disponíveis

- GET /users - Lista todos os usuários
- POST /users - Cria um novo usuário

## Estrutura do projeto

- `api.go` - Código fonte da API
- `docs/swagger.json` - Documentação OpenAPI/Swagger
- `go.mod` - Gerenciamento de dependências 