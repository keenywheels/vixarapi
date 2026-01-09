# Vixar backend

Основной репозиторий бэкэнда [vixar.tech](https://vixar.tech). Репозиторий состоит из двух частей:
- processor: процессор, который обрабатывает задачи из очереди (kafka)
- vixarapi: веб-сервер, которые предаставляет API для фронта

## Локальное тестирование
Для того, чтобы локально запустить проект достаточно прописать:
```bash
make docker-build && make docker-up
```

Чтобы было проще тестировать, можно добавить тестового пользователя, для этого добавляем юзера в postgres:
```sql
INSERT INTO users (username, email, vkid) VALUES ('test_user', 'test_user@mail.ru', 1234)
```

Затем добавляем в redis сессию:
```redis
SET "test_session" '{"id": "<USER_ID>", "username": "test_user", "email": "test_user@mail.ru", "vkid": 1234, "tguser": null}'
```
где `<USER_ID>` это id пользователя, которого мы только что добавили в postgres (можно узнать через `SELECT id FROM users WHERE username='test_user'`).

теперь можно использовать "test_session" в куке "session_id" для запросов к приватным ручкам.