# Vira Auth Middleware

Пакет `vira-middleware` предоставляет middleware для аутентификации и логирования HTTP-запросов с использованием JWT.

## Особенности

- Проверка JWT токенов в заголовке Authorization
- Извлечение user_id из токена и сохранение в контексте
- Интеграция с логгером для структурированного логирования
- Поддержка контекстного логирования (request_id, user_id и др.)
- Совместимость с `vira-config`, `vira-jwt` и `vira-logger`

## Установка

```bash
go get github.com/skrolikov/vira-middleware
```

## Использование

### Инициализация middleware

```go
import (
    "github.com/skrolikov/vira-middleware"
    "github.com/skrolikov/vira-config"
    "github.com/skrolikov/vira-logger"
)

func main() {
    cfg := config.Load()
    logger := logger.New(logger.Config{
        Level: logger.INFO,
        // другие параметры логгера
    })

    // Создаем middleware
    authMiddleware := middleware.Auth(cfg, logger)
    loggerMiddleware := middleware.ContextLogger(logger)

    // Настраиваем маршруты
    http.Handle("/secure", loggerMiddleware(authMiddleware(secureHandler)))
}
```

### Middleware Auth

Проверяет JWT токен и сохраняет user_id в контексте запроса.

**Требования:**
- Токен должен быть в заголовке `Authorization` с префиксом `Bearer `
- Токен должен содержать поле `user_id`

**Поведение при ошибках:**
- Возвращает 401 Unauthorized при отсутствии/невалидности токена
- Логирует все ошибки аутентификации

### Middleware ContextLogger

Добавляет логгер в контекст запроса с уже привязанными полями из контекста.

### Получение user_id в обработчике

```go
func secureHandler(w http.ResponseWriter, r *http.Request) {
    userID := middleware.GetUserID(r)
    // ...
}
```

### Использование логгера в обработчике

```go
func secureHandler(w http.ResponseWriter, r *http.Request) {
    logger := r.Context().Value(middleware.loggerKey).(*log.Logger)
    logger.Info("Обработка запроса")
    // ...
}
```

## Пример логирования

При успешной аутентификации логгер запишет:

```
[INFO] 2023-10-01T15:04:05Z auth.go:56 Успешная аутентификация | path=/secure method=GET user_id=123
```

При ошибках:

```
[WARN] 2023-10-01T15:04:05Z auth.go:42 Auth middleware: отсутствует токен | path=/secure method=GET
[WARN] 2023-10-01T15:04:05Z auth.go:49 Auth middleware: неверный токен | path=/secure method=GET error="token is expired"
```

## Лучшие практики

1. Всегда используйте `ContextLogger` перед `Auth` middleware
2. Для защищенных маршрутов применяйте оба middleware
3. Используйте полученный из контекста логгер для логирования в обработчиках
4. Добавляйте дополнительные поля в логгер через `WithFields` при необходимости
5. Настройте корректный уровень логирования в зависимости от окружения

## Интеграция

Пакет работает с:
- `vira-config` - для получения JWT секрета
- `vira-jwt` - для парсинга токенов
- `vira-logger` - для структурированного логирования