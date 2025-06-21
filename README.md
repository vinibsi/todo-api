# todo-api

## Estrutura de diretórios do projeto 
```text
todo-api/
├── cmd/
│   └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── dto/
│   │   └── todo_dto.go
│   ├── entity/
│   │   └── todo.go
│   ├── controller/
│   │   └── todo_controller.go
│   ├── service/
│   │   └── todo_service.go
│   └── repository/
│       └── todo_repository.go
├── pkg/
│   └── database/
│       └── connection.go
├── test/
│   ├── integration/
│   │   ├── todo_integration_test.go
│   │   └── test_helpers.go
│   └── unit/
│       ├── controller/
│       │   └── todo_controller_test.go
│       ├── service/
│       │   └── todo_service_test.go
│       └── repository/
│           └── todo_repository_test.go
├── mocks/
│   ├── mock_todo_repository.go
│   └── mock_todo_service.go
└── go.mod
```

## Comandos para executar
```shell
$ go mod init todo-api 
$ go mod tidy 
# Criar arquivo .env com as variáveis de ambiente
$ go run cmd/main.go 
```

## Endpoints disponíveis
```text
GET    /v1/todos              - Lista todas as tarefas (com paginação)
GET    /v1/todos/:id          - Busca tarefa por ID
POST   /v1/todos              - Cria nova tarefa
PUT    /v1/todos/:id          - Atualiza tarefa
DELETE /v1/todos/:id          - Deleta tarefa
PATCH  /v1/todos/:id/complete - Marca tarefa como concluída
```