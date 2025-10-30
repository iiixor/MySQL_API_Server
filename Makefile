# MySQL TUI Editor - Makefile
# Быстрая сборка бинарников для разных платформ

.PHONY: all clean build-editor build-server build-all help
.PHONY: build-macos build-windows build-linux
.PHONY: build-editor-macos-arm build-editor-macos-intel build-editor-windows
.PHONY: test run-server run-editor

# Версия приложения
VERSION ?= 1.0.0
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Пути к бинарникам
BIN_DIR := bin
EDITOR_DIR := editor
SERVER_DIR := server

# Имена бинарников
EDITOR_NAME := mysql-editor
SERVER_NAME := mysql-server

# Build flags
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

# По умолчанию показать help
.DEFAULT_GOAL := help

## help: Показать эту справку
help:
	@echo "MySQL TUI Editor - Makefile команды:"
	@echo ""
	@echo "  make build-all              - Собрать все бинарники (editor + server)"
	@echo "  make build-editor           - Собрать editor для текущей платформы"
	@echo "  make build-server           - Собрать server для текущей платформы"
	@echo ""
	@echo "  make build-macos            - Собрать editor для macOS (ARM64 + Intel)"
	@echo "  make build-windows          - Собрать editor для Windows"
	@echo "  make build-linux            - Собрать editor для Linux"
	@echo ""
	@echo "  make clean                  - Удалить все собранные бинарники"
	@echo "  make test                   - Запустить тесты"
	@echo "  make run-server             - Запустить сервер"
	@echo "  make run-editor             - Запустить редактор"
	@echo ""

## clean: Удалить все собранные бинарники
clean:
	@echo "🧹 Очистка бинарников..."
	@rm -rf $(BIN_DIR)
	@rm -f editor/editor server/server
	@echo "✅ Очистка завершена"

## build-all: Собрать все компоненты
build-all: build-editor build-server
	@echo "✅ Все компоненты собраны"

## build-editor: Собрать editor для текущей платформы
build-editor:
	@echo "🔨 Сборка editor для текущей платформы..."
	@mkdir -p $(BIN_DIR)
	@cd $(EDITOR_DIR) && go build $(LDFLAGS) -o ../$(BIN_DIR)/$(EDITOR_NAME) ./cmd/editor
	@echo "✅ Editor собран: $(BIN_DIR)/$(EDITOR_NAME)"

## build-server: Собрать server
build-server:
	@echo "🔨 Сборка server..."
	@mkdir -p $(BIN_DIR)
	@cd $(SERVER_DIR) && go build $(LDFLAGS) -o ../$(BIN_DIR)/$(SERVER_NAME) ./cmd/server
	@echo "✅ Server собран: $(BIN_DIR)/$(SERVER_NAME)"

## build-macos: Собрать editor для macOS (ARM64 + Intel)
build-macos: build-editor-macos-arm build-editor-macos-intel
	@echo "✅ macOS бинарники собраны"

## build-editor-macos-arm: Собрать editor для macOS Apple Silicon (ARM64)
build-editor-macos-arm:
	@echo "🍎 Сборка editor для macOS ARM64..."
	@mkdir -p $(BIN_DIR)
	@cd $(EDITOR_DIR) && GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) \
		-o ../$(BIN_DIR)/$(EDITOR_NAME)-macos-arm64 ./cmd/editor
	@echo "✅ Готово: $(BIN_DIR)/$(EDITOR_NAME)-macos-arm64"

## build-editor-macos-intel: Собрать editor для macOS Intel (AMD64)
build-editor-macos-intel:
	@echo "🍎 Сборка editor для macOS Intel..."
	@mkdir -p $(BIN_DIR)
	@cd $(EDITOR_DIR) && GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) \
		-o ../$(BIN_DIR)/$(EDITOR_NAME)-macos-intel ./cmd/editor
	@echo "✅ Готово: $(BIN_DIR)/$(EDITOR_NAME)-macos-intel"

## build-windows: Собрать editor для Windows
build-windows: build-editor-windows

## build-editor-windows: Собрать editor для Windows (AMD64)
build-editor-windows:
	@echo "🪟 Сборка editor для Windows..."
	@mkdir -p $(BIN_DIR)
	@cd $(EDITOR_DIR) && GOOS=windows GOARCH=amd64 go build $(LDFLAGS) \
		-o ../$(BIN_DIR)/$(EDITOR_NAME)-windows.exe ./cmd/editor
	@echo "✅ Готово: $(BIN_DIR)/$(EDITOR_NAME)-windows.exe"

## build-linux: Собрать editor для Linux (AMD64)
build-linux:
	@echo "🐧 Сборка editor для Linux..."
	@mkdir -p $(BIN_DIR)
	@cd $(EDITOR_DIR) && GOOS=linux GOARCH=amd64 go build $(LDFLAGS) \
		-o ../$(BIN_DIR)/$(EDITOR_NAME)-linux ./cmd/editor
	@echo "✅ Готово: $(BIN_DIR)/$(EDITOR_NAME)-linux"

## build-release: Собрать релизные бинарники для всех платформ
build-release: clean build-editor-macos-arm build-editor-macos-intel build-editor-windows build-linux build-server
	@echo ""
	@echo "📦 Релизные бинарники собраны:"
	@ls -lh $(BIN_DIR)/
	@echo ""
	@echo "✅ Готово для распространения!"

## test: Запустить тесты
test:
	@echo "🧪 Запуск тестов..."
	@cd server && go test -v ./...
	@echo "✅ Тесты завершены"

## test-coverage: Запустить тесты с покрытием
test-coverage:
	@echo "🧪 Запуск тестов с покрытием..."
	@cd server && go test -cover -coverprofile=coverage.out ./...
	@cd server && go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Отчёт о покрытии: server/coverage.html"

## run-server: Запустить сервер для разработки
run-server:
	@echo "🚀 Запуск server..."
	@cd $(SERVER_DIR) && go run ./cmd/server -config config/config.yml

## run-editor: Запустить редактор для разработки
run-editor:
	@echo "🚀 Запуск editor..."
	@cd $(EDITOR_DIR) && go run ./cmd/editor

## deps: Установить/обновить зависимости
deps:
	@echo "📦 Установка зависимостей..."
	@cd editor && go mod tidy
	@cd server && go mod tidy
	@echo "✅ Зависимости установлены"

## fmt: Форматировать код
fmt:
	@echo "🎨 Форматирование кода..."
	@cd editor && go fmt ./...
	@cd server && go fmt ./...
	@echo "✅ Код отформатирован"

## vet: Проверить код на ошибки
vet:
	@echo "🔍 Проверка кода..."
	@cd editor && go vet ./...
	@cd server && go vet ./...
	@echo "✅ Проверка завершена"

## check: Запустить все проверки (fmt, vet, test)
check: fmt vet test
	@echo "✅ Все проверки пройдены"

## info: Показать информацию о сборке
info:
	@echo "ℹ️  Информация о сборке:"
	@echo "  Версия:        $(VERSION)"
	@echo "  Время сборки:  $(BUILD_TIME)"
	@echo "  Git commit:    $(GIT_COMMIT)"
	@echo "  Go версия:     $$(go version)"
	@echo "  Платформа:     $$(go env GOOS)/$$(go env GOARCH)"
