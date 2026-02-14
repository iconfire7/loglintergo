# loglintergo

`loglintergo` — кастомный линтер для `golangci-lint`, который проверяет текст лог-сообщений (`slog`, `zap`) по набору правил.

## Что проверяет линтер

Сейчас в проекте есть 4 правила:

- `LOG001` — сообщение не должно начинаться с заглавной буквы.
- `LOG002` — сообщение должно быть на английском (латиница).
- `LOG003` — сообщение не должно содержать emoji и non-ASCII символы.
- `LOG004` — сообщение не должно содержать чувствительные слова (`password`, `token`, `bearer` и т.д.).

## Структура проекта

```text
.
├── internal/
│   ├── analyzer/loglinter/   # анализатор go/analysis
│   ├── config/               # структура конфигурации линтера
│   └── rules/                # реализация правил + тесты
├── plugin/                   # точка входа плагина для golangci-lint
├── testdata/                 # пример исходников для локальной проверки
├── .custom-gcl.yml           # конфиг сборки custom golangci-lint бинаря
├── .golangci.yml             # конфиг запуска линтера + его настройки
├── Makefile                  # удобные команды для разработки
└── golangci-lint/            # локальный бинарь golangci-lint (если используется)
```

## Требования

- Go 1.23+
- `golangci-lint` с поддержкой custom module plugin system (`golangci-lint custom`)

## Конфигурация линтера

Линтер настраивается в `.golangci.yml` в блоке:

```yaml
linters:
  settings:
    custom:
      loglintergo:
        type: module
        description: "Custom log linter"
        settings:
          rules:
            lowercase: true
            english: true
            emoji_or_spesial: true
            sensitive: true
```

## Запуск: вариант 1 (через Makefile)

1. Подтянуть зависимости:

```bash
make tidy
```

2. Собрать custom golangci-lint с подключенным плагином:

```bash
make custom-linter
```

3. Запустить проверку:

```bash
make lint-golangci
```

Дополнительно:

```bash
make test   # go test ./...
make clean  # чистка артефактов
```

## Запуск: вариант 2 (вручную)

1. Установить зависимости:

```bash
go mod tidy
```

2. Собрать custom бинарь golangci-lint:

```bash
./golangci-lint/golangci-lint custom
```

3. Запустить линтер через собранный binary (`./custom-gcl`):

```bash
./custom-gcl run --config .golangci.yml ./testdata/src/...
```

## Быстрый локальный цикл разработки

- Меняете правила в `internal/rules`.
- Запускаете тесты:

```bash
go test ./...
```

- Прогоняете линтер на `testdata`:

```bash
make lint-golangci
```

## Полезные замечания

- Линтер анализирует только статический текст первого аргумента лог-вызова.
- Поддерживаются вызовы `slog` и `zap` (включая `SugaredLogger`).
- Если конфигурация не передана, используются значения по умолчанию (все правила включены).
