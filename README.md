# Person Service API

## Описание
RESTful сервис для управления данными о людях с обогащением (возраст, пол, национальность).

## Функционал
- **POST /api/v1/persons**: Создать персону.
- **GET /api/v1/persons**: Список персон (пагинация, фильтр по имени).
- **GET /api/v1/persons/{id}**: Получить персону по ID.
- **PUT /api/v1/persons/{id}**: Обновить персону.
- **DELETE /api/v1/persons/{id}**: Удалить персону.
- Документация: `/swagger/index.html`.

## Технологии
- Go, Gin, PostgreSQL, Swagger.

## Установка
1. Клонировать: `git clone https://github.com/username/repo.git`
2. Установить: `go mod tidy`
3. Запустить БД: `docker run --name pg -e POSTGRES_PASSWORD=pass -p 5432:5432 -d postgres`
4. Настроить `.env`
