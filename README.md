# VolleyReg — бэкенд

Регистрация участников на спортивное соревнование (тема 16): **участники** (услуги) и **заявки команд** на участие в Первенстве ЗГУ по волейболу 2026.

- Фронтенд (React SPA, PWA, Tauri): https://github.com/AmirSuleimanov-IS/volleyreg-frontend
- Демонстрация фронтенда в GitHub Pages: https://amirsuleimanov-is.github.io/volleyreg-frontend/
- Курсовая работа «Веб-сервис регистрации участников на спортивное соревнование», дисциплина «Разработка Web-приложений», студент Сулейманов А. Ш., группа ИС-23, ЗГУ им. Н.М. Федоровского, 2026.

Стек: Go, Gin, html/template, GORM + PostgreSQL, Minio, Redis, JWT, Swagger (swaggo). REST API — `/api/...`, Swagger — http://localhost:8080/swagger/index.html, Postman-коллекции — `postman/`.

Для фронтенда на GitHub Pages и приложения Tauri включён CORS; адрес Minio в ответах задаётся `MINIO_PUBLIC_URL` (например, IP компьютера в локальной сети).

Тестовые пользователи (пароль `123456`): `sever_rep`, `polar_rep` — представители команд, `organizer` — организатор соревнования. Каждая лабораторная работа — отдельная ветка (`lab1-volley-design`, `lab2-volley-db`, …).

## Запуск
```
docker compose up -d
go run ./cmd/migrate
docker exec -i volleyreg-db psql -U volleyreg -d volleyreg < resources/db/seed.sql
go run ./cmd/volleyreg
```
Приложение: http://localhost:8080, консоль Minio: http://localhost:9001 (minioadmin / minioadmin), Adminer: http://localhost:8081 (сервер db, volleyreg / volleyreg).

Медиафайлы: см. `resources/media/ATTRIBUTION.md`.
