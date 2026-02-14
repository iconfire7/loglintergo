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
├── testdata/                 # примеры исходников для локальной проверки
│   └── src/errs/main.go      # демонстрационный файл с намеренно добавленными ошибками
├── .custom-gcl.yml           # конфиг сборки custom golangci-lint бинаря
├── .golangci.yml             # конфиг запуска линтера + его настройки
├── Makefile                  # удобные команды для разработки
└── golangci-lint/            # локальный бинарь golangci-lint (если используется)
```

## Требования

- Go 1.23+
- `golangci-lint` с поддержкой custom module plugin system (`golangci-lint custom`)
- утилита `golangci-lint` с поддержкой custom module plugin system (`golangci-lint custom`)
- `make` (опционально, если используете команды из `Makefile`)

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
            emoji_or_special: true
            sensitive: true
          sensitive_patterns:
            - '(?i)\b(token|secret|api[_-]?key)\b\s*[:=]'
            - '(?i)\bauthorization\b\s*:\s*bearer\b'
```
- Значения в `rules` можно менять прямо в конфиге (`true/false`), чтобы включать или выключать отдельные проверки.
- В `sensitive_patterns` можно добавлять свои регулярные выражения для поиска чувствительных данных в логах.

## Быстрый старт (через Makefile)

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

## Ручной запуск

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

4. Запустить автоисправление поддерживаемых линтеров:

```bash
./custom-gcl run --fix --config .golangci.yml ./testdata/src/...
```

## Проверка на демонстрационных ошибках

В `testdata/src/errs/main.go` специально оставлены ошибки в лог-сообщениях, чтобы можно было проверить работу линтера на реальных примерах.

- Проверить ошибки:

```bash
./custom-gcl run --config .golangci.yml ./testdata/src/...
```

- Попробовать исправить автоматически:

```bash
./custom-gcl run --fix --config .golangci.yml ./testdata/src/...
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

- Поддерживаются вызовы `slog` и `zap` (включая `SugaredLogger`).
- Если конфигурация не передана, используются значения по умолчанию (все правила включены).
