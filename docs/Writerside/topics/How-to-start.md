# Начало работы
Для работы приложения понадобится docker

## Запуск проекта
* клонируйте репозиторий и перейдите в корневую папку:
    ```
  git clone https://github.com/AleksandrVishniakov/distributed-calculator
  cd distributed-calculator
  ```
* выполните команду ```make``` или выполните команды из Makefile последовательно вручную:
    ```
  docker build -t dc-api-gateway:local ./api-gateway
  docker build -t dc-daemon:local ./daemon
  docker build -t dc-auth:local ./auth
  docker build -t dc-page-parser:local ./page-parser
  docker compose up
  ```
* Если при запуске прокета возникли ошибки инициализации базы данных или у мигратора (goose), то помогает перезапуск (`make` или `docker compose up`). Если возники другие ошибки, пожалуйста, напишите [сюда](https://t.me/landowner7)
* Готово! GUI приложения доступно на [http://localhost:8080](http://localhost:8080)

## О проекте
- интеркативный gui, связанный с сервером по http
- взаимодействие оркестратор-демон по grpc
- сохранение состояния в базу данных postgresql
- возможность восстановления после повторного включения
- proto-файлы можно найти в [репозитории](https://github.com/AleksandrVishniakov/dc-protos)

## Изменения в структуре приложения
Можно изменить порты, количество агентов, максимальное количество горутин и т.д. с помощью файла ```docker-compose.yml```:
### Сервисы docker-compose
### Api Gateway
api-gateway - основной сервис приложения. Сюда приходят запросы из GUI, а также агентов
  ```yaml
  api-gateway:
    image: dc-api-gateway:local
    container_name: dc-api-gateway
    environment:
      HTTP_PORT: 8000
      WORKERS_MONITORING_PERIOD_MS: 30000
      DB_PASSWORD: "admin"
    ports:
      - "8000:8000"
    depends_on:
      goose:
        condition: service_completed_successfully
  ```
Переменные окружения:
* `HTTP_PORT` - порт, на которм работает сервер. При изменении необходимо также изменить ```ports```
* `WORKERS_MONITORING_PERIOD_MS` - период в миллисекундах, через который сервер проверяет, получен ли ping от всех агентов и удаляет неактивные
* `DB_PASSWORD` - пароль для базы данных PostgreSQL

### Daemon
daemon - агент, который выполняет арифметические вычисления
  ```yaml
  daemon1:
    depends_on:
      - api-gateway
    image: dc-daemon:local
    container_name: dc-daemon-1
    environment:
      HTTP_PORT: 8001
      DAEMON_ID: 1
      DAEMON_HOST: "http://daemon1:8001" 
      ORCHESTRATOR_HOST: "http://api-gateway:8000"
      PING_PERIOD_MS: 25000
      MAX_GOROUTINES: 1
    ports:
      - "8001:8001"
  ```
Переменные окружения:
* `HTTP_PORT` - порт, на которм работает сервер. При изменении необходимо также изменить ```ports``` и ```DAEMON_HOST```
* `DAEMON_ID` - иденитификатор демона, уникальный для каждого демона
* `DAEMON_HOST` - адрес агента, по которому к нему можно обратиться
* `ORCHESTRATOR_HOST` - адрес оркестратора (api-gateway)
* `PING_PERIOD_MS: 25000` - период в миллисекудах, через который агент оправляет ping к оркестратору
* `MAX_GOROUTINES` - маскимальное количество горутин, которые могут работать внутри агента


Также можно добавить дополнительных агентов, изменив их названия, порты и идентификаторы

### Page Parser
Сервис отображает графический интерфейс
```yaml
page-parser:
    depends_on:
      - api-gateway
    image: dc-page-parser:local
    container_name: dc-page-parser
    environment:
      HTTP_PORT: "8080"
      ORCHESTRATOR_HOST: "http://localhost:8000"
    ports:
      - "8080:8080"
```
Переменные окружения:
* `HTTP_PORT` - порт, на которм работает сервер. При изменении необходимо также изменить ```ports```
* `ORCHESTRATOR_HOST` - адрес оркестратора (api-gateway)

[По вопросам](https://t.me/landowner7)
