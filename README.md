## backend

run locally:
```
make run
```
requires:
- make
- dotenv

or docker compose:
```
make compose up
```
both requires .env file for example:
```
HTTP_HOST=localhost
HTTP_PORT=4000
HTTP_READ_TIMEOUT=15s
HTTP_WRITE_TIMEOUT=15s

PG_USER=user
PG_PASS=pass
PG_HOST=localhost
PG_PORT=5432
PG_DBNAME=db
PG_SSL=disable
PG_DSN=postgresql://${PG_USER}:${PG_PASS}@${PG_HOST}:${PG_PORT}/${PG_DBNAME}?sslmode=${PG_SSL}

JWT_KEY="verysecretkey"

STEAM_API_KEY=apikey #https://steamcommunity.com/dev/apikey

SWAGGER_URL=/api/swagger/doc.json
```
