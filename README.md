# Ссылки
[mindmap структура записи продажи](https://www.mindmeister.com/app/map/3281518214?t=WEI6jjzFF4)


# Тестовый эдпоинт
curl --location 'localhost:8080/receiver/cms/v1'

## Импорт для миграции
export DB_HOST=31.131.254.247 DB_PORT=5432 DB_USER=erp_db_usr DB_PASSWORD=E4SeIFlKFSCRBQIz DB_NAME=erp_db

goose -dir="data_receiver_db/" postgres "$DSN_QUEUE" up
