# HTTP мультиплексор

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

* приложение представляет собой http-сервер с одним хендлером
* хендлер на вход получает POST-запрос со списком url в json-формате
* сервер запрашивает данные по всем этим url и возвращает результат клиенту в json-формате
* если в процессе обработки хотя бы одного из url получена ошибка, обработка всего списка прекращается и клиенту возвращается текстовая ошибка
* для реализации задачи следует использовать Go 1.13 или выше
* использовать можно только компоненты стандартной библиотеки Go
* сервер не принимает запрос если количество url в нем больше 20
* сервер не обслуживает больше чем 100 одновременных входящих http-запросов
* для каждого входящего запроса должно быть не больше 4 одновременных исходящих
* таймаут на запрос одного url - секунда
* обработка запроса может быть отменена клиентом в любой момент, это должно повлечь за собой остановку всех операций связанных с этим запросом
* сервис должен поддерживать 'graceful shutdown'

## Флаги

По умолчанию сервер ведет себя так, как описано выше. Но можно изменить поведение с помощью флагов.

```text
Usage of ./httpMultiplexer:
  -maxIncome uint
        Quantity of incoming http requests working simultaneously (default 100)
  -maxOutgoing int
        Quantity of outgoing http requests per each incoming request (default 4)
  -path string
        HTTP path (default "/api/v1")
  -port string
        HTTP listen port (default "8080")
  -timeoutPerRequest duration
        Timeout for each outgoing request, 1s for example (default 1s)
  -timeoutStop duration
        Timeout for graceful shutdown, 20s for example (default 30s)
  -urlQuantity int
        Maximum quantity of urls in each incoming request (default 20)
```

## Покрытие тестами

Для тестов логики обработки запросов создан небольшой [сервер](test/server/server.go), который имитирует работу внешнего сервиса. Позволяет проверить успешные, неуспешные и слишком долгие запросы. В [test.http](test.http) приведены тесты запросов и ответов на запущенном приложении.

```text 



```text
?       github.com/akrillis/affise-http-multiplexer/cmd/httpMultiplexer [no test files]
?       github.com/akrillis/affise-http-multiplexer/internal/auxiliary  [no test files]
ok      github.com/akrillis/affise-http-multiplexer/internal/check      0.434s  coverage: 33.3% of statements
ok      github.com/akrillis/affise-http-multiplexer/internal/handler    2.742s  coverage: 92.2% of statements
ok      github.com/akrillis/affise-http-multiplexer/internal/limit      0.234s  coverage: 46.2% of statements
ok      github.com/akrillis/affise-http-multiplexer/internal/server     0.913s  coverage: 60.0% of statements
?       github.com/akrillis/affise-http-multiplexer/service     [no test files]
?       github.com/akrillis/affise-http-multiplexer/test/server [no test files]
?       github.com/akrillis/affise-http-multiplexer/types       [no test files]
```
