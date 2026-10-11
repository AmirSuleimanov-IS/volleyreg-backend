# VolleyReg — бэкенд

Регистрация участников на спортивное соревнование (тема 16): **участники** (услуги) и **заявки команд** на участие в Первенстве ЗГУ по волейболу 2026.

Стек: Go, Gin, html/template, Minio. Каждая лабораторная работа — отдельная ветка (`lab1-volley-design`, `lab2-volley-db`, …).

## Запуск
```
docker compose up -d
go run ./cmd/volleyreg
```
Приложение: http://localhost:8080, консоль Minio: http://localhost:9001 (minioadmin / minioadmin).

Медиафайлы: см. `resources/media/ATTRIBUTION.md`.
