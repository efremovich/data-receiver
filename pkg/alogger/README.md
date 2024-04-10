# ALogger пакет для логирования

Пакет ALogger предоставляет набор методов для регистрации событий.

## Использование

Для использования ALogger выполните следующие шаги:

* Создайте конфигурацию логгера, которую можно передать для создания экземпляра ALogger.

```go
cfg := &alogger.Config{
Level:       alogger.LevelDebug,
}
```

- **Level** - задаёт уровень логирования "debug","info","warn","error", по умолчанию - "info".
- **Output** задаёт поток вывода событий, по умолчанию - os.Stdout.
- **TextFormat** включает вывод в текстовом человекочитаемом виде, по умолчанию использует формат JSON.
- **OnlyError** включает логирование событий только при событии типа, по умолчанию логирует при отсутствии событий
  уровня "error".

Так же есть возможность заменить дефолтную конфигурацию для ALogger.

```go
alogger.SetDefaultConfig(&alogger.Config{
Level:       alogger.LevelInfo,
})
```

* Создайте новый экземпляр ALogger с помощью функции NewALogger. Если вместо конфигурации передать nil - будет
  использоваться дефолтная конфигурация.

```go
log := alogger.NewALogger(context.Background(), cfg)
```

* Отложите вызов метода Print, чтобы гарантировать логирование событий в конце вашей функции или программы. Если при
  выводе событий произойдет ошибка - она будет выведена в StdOut и возвращена.

```go
defer log.Flush()
```

* Используйте различные методы, предоставляемые ALogger, для создания событий разных уровней:

```go
log.Debugf("debug message")
log.Infof("info message")
log.Warnf("warn message")
log.Errorf("info message")
```

* Вы также можете установить дополнительные атрибуты события с помощью метода SetAttrs.

```go
log.Errorf("warn message").SetAttrs(map[string]interface{}{
"key1": "value1",
"key2": true,
})
```

* Вы также можете обернуть ошибку с помощью метода Wrap. Поддерживается тип ошибок AError.

```go
log.Errorf("error message").Wrap(errors.New("new Errors"))
Пример
```

### Пример, демонстрирующий использование пакета ALogger:

```go
package main

import (
	"context"
	"errors"

	"alogger"
)

func main() {
	// Создает контекст с trace_id
	ctx := context.WithValue(context.Background(), alogger.TraceIdKey, "uuid123")

	// Устанавливаем дефолтную конфигурацию alogger. Можно установить при запуске сервиса.
	alogger.SetDefaultConfig(&alogger.Config{
		Level: alogger.LevelInfo,
	})

	// Создаем экземпляр логгера с дефолтной конфигурацией
	log := alogger.NewALogger(ctx, nil)

	// Логгируем события по завершению полезной работы
	defer log.Flush()

	// Логгируем ошибку, добавляем атрибуты и стек трейс
	log.Errorf("%s", "my message for error").
		Wrap(errors.New("new Errors")).
		SetAttrs(map[string]interface{}{
			"key1": "value1",
			"key2": true,
		}).Stack()
}
```

### Пример, демонстрирующий использование пакета ALogger без создания экземляра ALogger:

```go
package main

import (
	"context"
	"errors"

	"alogger"
)

func main() {
	// Создает контекст с trace_id
	ctx := context.WithValue(context.Background(), alogger.TraceIdKey, "uuid123")

	// Устанавливаем дефолтную конфигурацию alogger. Можно установить при запуске сервиса.
	alogger.SetDefaultConfig(&alogger.Config{
		Level: alogger.LevelInfo,
	})

	// Логирование без создания экземпляра логгера, логирование произойдет с использованием дефолтной конфигурации и сразу будет отправлена на вывод.
	stackOn := true

	// Добавление атрибутов к событию
	attrs := map[string]interface{}{
		"key1": "value1",
		"key2": true,
	}

	// Некая исходная ошибка
	err := errors.New("MyError")

	// Логирование события без экземпляра ALogger
	alogger.ErrorFromCtx(ctx, "Message", err, attrs, stackOn)
}
```

## Middleware для пакета ALogger

### gRPC

Middleware для проброса `trace_id` из контекста gRPC запроса в контекст обработчика этого запроса.

- `alogger.UnaryTraceIdInterceptor`
- `alogger.StreamTraceIdInterceptor`

#### Пример использования можно посмотреть в `alogger/tests/middleware/grpc_mock_server.go`

---

### fasthttp

Middleware для проброса `trace_id` из заголовков HTTP запроса в контекст обработчика этого запроса,
а так же в заголовки ответа на данный запрос.

- `alogger.TraceIdMiddlewareFastHTTP`

#### Пример использования можно посмотреть в `alogger/tests/middleware/fasthttp_mock_server.go`

---

### pgx

Middleware для логирования sql запросов (сам запрос, аргументы, время выполнения и `trace_id`, если он был передан в
контексте). Логи пишутся с уровнем `debug`

- `alogger.PGXQueryTracer`

Возможны несколько вариантов использования:

- модуль `github.com/jackc/pgx/v5` + драйвер `pgx`
- модуль `github.com/jmoiron/sqlx` + драйвер `pgx`
- модуль `database/sql` + драйвер `pgx`

#### Примеры использования

БД во всех примерах использования состоит из таблицы `foobar`, которая следующий вид:

| id | foo  | bar   | baz |
|----|------|-------|-----|
| 1  | aaaa | true  | 10  |
| 2  | bbbb | false | 20  |

```go
package main

import (
  "context"
  "database/sql"
  "os"

  "github.com/jackc/pgx/v5"
  "github.com/jackc/pgx/v5/stdlib"
  "github.com/jmoiron/sqlx"

  "alogger"
)

type R struct {
  Id  int
  Foo string
  Bar bool
  Baz int
}

const (
  connString   = "postgres://username:password@host:port/dbname"
  testTraceId  = "foo123456789"
  testSQLQuery = "select * from foobar where id = $1;"
)

func main() {
  alogger.SetDefaultConfig(
    &alogger.Config{
      Output:     os.Stdout,
      Level:      alogger.LevelDebug,
      OnlyError:  false,
      TextFormat: true,
    })

  examplePGXWithDriverPGX()

  exampleDbSQLWithDriverPGX()

  exampleSQLXWithDriverPGX()
}

// examplePGXWithDriverPGX - Пример использования с модулем github.com/jackc/pgx/v5 и драйвером pgx
func examplePGXWithDriverPGX() {
  pgxConf, _ := pgx.ParseConfig(connString)
  pgxConf.Tracer = new(alogger.PGXQueryTracer)

  conn, _ := pgx.ConnectConfig(context.Background(), pgxConf)
  defer func() { _ = conn.Close(context.Background()) }()

  ctx := context.WithValue(context.Background(), alogger.TraceIdKey, testTraceId)

  row := conn.QueryRow(ctx, testSQLQuery, 1)
  r := new(R)
  _ = row.Scan(&r.Id, &r.Foo, &r.Bar, &r.Baz)
}

// exampleDbSQLWithDriverPGX - Пример использования с модулем database/sql и драйвером pgx
func exampleDbSQLWithDriverPGX() {
  pgxConf, _ := pgx.ParseConfig(connString)
  pgxConf.Tracer = &alogger.PGXQueryTracer{SkipFirstArg: true}

  conn, _ := sqlx.ConnectContext(context.Background(), "pgx", stdlib.RegisterConnConfig(pgxConf))
  defer func() { _ = conn.Close() }()

  ctx := context.WithValue(context.Background(), alogger.TraceIdKey, testTraceId)

  row := conn.QueryRowContext(ctx, testSQLQuery, 1)

  r := new(R)
  _ = row.Scan(&r.Id, &r.Foo, &r.Bar, &r.Baz)
}

// exampleSQLXWithDriverPGX- Пример использования с модулем github.com/jmoiron/sqlx и драйвером pgx
func exampleSQLXWithDriverPGX() {
  pgxConf, _ := pgx.ParseConfig(connString)
  pgxConf.Tracer = &alogger.PGXQueryTracer{SkipFirstArg: true}

  conn, _ := sql.Open("pgx", stdlib.RegisterConnConfig(pgxConf))
  defer func() { _ = conn.Close() }()

  ctx := context.WithValue(context.Background(), alogger.TraceIdKey, testTraceId)

  row := conn.QueryRowContext(ctx, testSQLQuery, 1)

  r := new(R)
  _ = row.Scan(&r.Id, &r.Foo, &r.Bar, &r.Baz)
}
```

Во всех примерах при выполнении SQL запроса в логах появятся одинаковые записи вида:
```
[01.03.2024 01:06:30.946] debug: SQL START
TraceId: foo123456789
PackageName: alogger
Caller: /home/derv-dice/projects/work/Edi.Operator/alogger/alogger.go:52
Attrs:
	args: [1]
	sql: select * from foobar where id = $1;

[01.03.2024 01:06:30.947] debug: SQL END
TraceId: foo123456789
PackageName: alogger
Caller: /home/derv-dice/projects/work/Edi.Operator/alogger/alogger.go:52
Attrs:
	args: [1]
	time: 1.398956ms
	sql: select * from foobar where id = $1;
```

---







