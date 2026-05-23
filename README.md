# Rinha de Backend 2023 Q3 - Go + Postgres

Implementacao da Rinha de Backend 2023 Q3 usando Go, Postgres, Nginx e Docker Compose.

## Stack

- Go
- Postgres
- Nginx
- Docker Compose

## Rodar

```bash
docker compose up --build
```

A API fica disponivel em:

```txt
http://localhost:9999
```

Para resetar o banco e subir do zero:

```bash
docker compose down -v
docker compose up --build
```

## Endpoints

- `POST /pessoas`
- `GET /pessoas/:id`
- `GET /pessoas?t=termo`
- `GET /contagem-pessoas`

## Exemplos

Criar pessoa:

```bash
curl -i -X POST http://localhost:9999/pessoas \
  -H "Content-Type: application/json" \
  -d '{
    "apelido": "ana",
    "nome": "Ana Barbosa",
    "nascimento": "1985-09-23",
    "stack": ["Node", "Postgres"]
  }'
```

Buscar por ID:

```bash
curl -i http://localhost:9999/pessoas/<id>
```

Buscar por termo:

```bash
curl -i "http://localhost:9999/pessoas?t=node"
```

Contar pessoas:

```bash
curl -i http://localhost:9999/contagem-pessoas
```

Buscar sem `t` deve retornar `400 Bad Request`:

```bash
curl -i http://localhost:9999/pessoas
```
