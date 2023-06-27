## Подключение агента Go

### Инициализация модуля

1. Добавьте модуль pobp-trace-go в ваше приложение:
```go env -w GOPRIVATE=git.proto.group```
```go get git.proto.group/protoobp/pobp-trace-go/pobptrace/tracer```
2. Импортируйте модуль в коде вашего приложения:
```go   
import(
  "git.proto.group/protoobp/pobp-trace-go/pobptrace/tracer"
  "git.proto.group/protoobp/pobp-trace-go/pobptrace/opentracer"
)
```
### Конфигурация модуля

 Укажите следующие пременные окружения:
  * ```POBP_AGENT_HOST="proto-backend"``` - Адрес Proto Backend сервера
  * ```POBP_TRACE_AGENT_PORT="9080"``` - Порт Proto Backend сервера
  * ```POBP_SERVICE="my_service_name"``` - Имя сервиса
