.PHONY: test test-unit test-integration test-coverage run build clean

# Executar todos os testes
test:
	go test -v ./...

# Executar apnenas tetes unitários
test-unit:
	go test -v ./test/unit/...

# Executar apnenas tetes de integração
test-integration:
	go test -v ./test/integration/...

# Executar testes com cobertura
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo ˜="Relatório de cobertura gerado em coverage.html"

# Executar a aplicação
run:
	go run cmd/main.go

# Fazer o build da aplicação
build:
	go build -o bin/todo-api cmd/main.go

clean:
	rm -f bin/todo-api coverage.out coverage.html

# Instarlar dependências
deps:
	go mod download
	go mod tidy

# Pensar em colocar um linter aqui...
# Executar linter
lint:
	golangci-lint run

# Executar testes em modo watch
test-watch:
	find . -name "*.go" | entr -r go test -v ./...
