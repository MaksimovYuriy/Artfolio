# Artfolio
Проект витрины работ для художницы Анны Аветисян

## Интеграционные тесты backend

Repository-тесты используют отдельный PostgreSQL на порту `55432`:

```bash
docker compose -f compose.integration.yml up -d --wait
cd backend
go test -tags=integration ./internal/repo/artwork ./internal/repo/social_link ./internal/controller/restapi/v1
cd ..
docker compose -f compose.integration.yml down
```

Для внешней тестовой БД можно передать `ARTFOLIO_TEST_DATABASE_URL`.
