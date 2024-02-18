# Distributed Calculator
## Документация
Документация по проекту доступна по [ссылке](https://aleksandrvishniakov.github.io/distributed-calculator/project-overview.html)
## Начало работы
Для работы приложения понадобится docker

### Запуск проекта
* клонируйте репозиторий и перейдите в корневую папку:
    ```Bash
  $ git clone github.com/AleksandrVishniakov/distributed-calculator
  $ cd distributed-calculator
  ```
* выполните команду ```make``` или выполните команды из Makefile последовательно вручную:
    ```Bash
  $ docker build -t dc-api-gateway:local ./api-gateway
  $ docker build -t dc-daemon:local ./daemon
  $ docker build -t dc-page-parser:local ./page-parser
  $ docker compose up
  ```
* Готово! GUI приложения доступно на [http://localhost:8080](http://localhost:8080)

### Изменения в структуре приложения
Можно изменить порты, количество агентов, макк=симальное количество горутин и т.д. с помощью файла ```docker-compose.yml```:
#### Сервисы docker-compose
#### Api Gateway
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

#### Daemon
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
* `HTTP_PORT` - порт, на которм работает сервер. При изменении необходимо также изменить ```ports```
* `DAEMON_ID` - иденитификатор демона
* `DAEMON_HOST` - адрес агента, по которому к нему можно обратиться
* `ORCHESTRATOR_HOST` - адрес оркестратора (api-gateway)
* `PING_PERIOD_MS: 25000` - период в миллисекудах, через который агент оправляет ping к оркестратору
* `MAX_GOROUTINES` - маскимальное количество горутин, которые могут работать внутри агента


Также можно добавить дополнительных агентов, изменив их названия, порты и идентификаторы

#### Page Parser
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
