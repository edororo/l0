Order Service (Go + Kafka + PostgreSQL)

Сервис для обработки заказов из Kafka, сохранения в PostgreSQL, кэширования в памяти и получения через HTTP API.

Возможности:

1. Получение заказов из Kafka
2. Сохранение заказов в PostgreSQL (транзакционно)
3. In-memory кеш с TTL и лимитом
4. REST API для получения заказа по ID
5. Простой web-интерфейс для просмотра заказов
6. SQL-миграции базы данных

Архитектура

Проект разделён на слои:

api — HTTP обработчики и фронт
service — бизнес-логика
db — работа с PostgreSQL
kafka — consumer сообщений
cache — кеш в памяти 
models — структуры данных

Запуск проекта и его работа:

1. Запустить инфраструктуру

docker compose up -d

Запустятся:

 PostgreSQL
  Kafka
   Zookeeper

2. Применить миграции

migrate -path ./migrations -database "postgres://order_user:password@localhost:5432/orders?sslmode=disable" up

Откат одной миграции:

migrate -path ./migrations -database "postgres://order_user:password@localhost:5432/orders?sslmode=disable" down 1

Сброс всех миграций:

migrate -path ./migrations -database "postgres://order_user:password@localhost:5432/orders?sslmode=disable" drop -f

3. Запустить сервис

go run cmd/main.go

Выведет: 
Подключение к PostgreSQL успешно
Кэш заказов создан
Kafka consumer запущен
HTTP сервер запущен на :8081

API:

Получить заказ:

GET http://localhost:8081/order/<order_uid>

Ответ:

json
{
  "order_uid": "test123",
  "items": [...]
}

4. Тест Kafka

Отправить заказ вручную:


docker exec -it kafka kafka-console-producer \
--broker-list localhost:9092 \
--topic orders

Вставь JSON заказа и нажми Enter.

Сервис автоматически:

1. Получит сообщение
2. Сохранит заказ в БД
3. Положит его в кеш

5. Web интерфейс:

Открой:

http://localhost:8081

Введи order_uid и нажми кнопку.

6. Кеш

 TTL удаляет старые записи автоматически
  Ограничение по количеству элементов
   Защита от переполнения памяти
    При повторном запросе заказ берётся из кеша

7. Требования

1. Go 1.22+
2. Docker
3. Kafka
4. PostgreSQL