**Windows (PowerShell):**

```powershell
# ВАЖНО: Копируем файл в контейнер, чтобы избежать проблем с кодировкой
docker cp dumps\data-only-dump.sql backend-db-1:/tmp/restore.sql
docker exec -i backend-db-1 psql -U postgres -d Medvito -f /tmp/restore.sql
```

**Linux/Mac:**

```bash
# Способ 1: Через pipe (работает корректно на Linux/Mac)
cat dumps/data-only-dump.sql | docker exec -i backend-db-1 psql -U postgres -d Medvito

# Способ 2: Копирование в контейнер (универсальный метод)
docker cp dumps/data-only-dump.sql backend-db-1:/tmp/restore.sql
docker exec -i backend-db-1 psql -U postgres -d Medvito -f /tmp/restore.sql
```
