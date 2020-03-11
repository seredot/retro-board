# Retro Board Service

Backend for real time retro board project.

## Running the service
```bash
go build ./... && ./retro-board
```

## Running tests
```
go test -v
```

## Endpoints
### Api health check
```bsh
curl --location --request GET 'http://127.0.0.1:8080/api'
```

### Create a new board
```bsh
curl --location --request POST 'http://127.0.0.1:8080/api/board'
```

### Get a board by ID
```bsh
curl --location --request GET 'http://127.0.0.1:8080/api/board/{{boardId}}'
```

### Add an item to a board
```bsh
curl --location --request POST 'http://127.0.0.1:8080/api/board/{{boardId}}/item' \
--header 'Content-Type: application/json' \
--data-raw '{
    "text": "This is an item",
    "color": "blue"
}'
```

### Update an item in a board
```bsh
curl --location --request PUT 'http://127.0.0.1:8080/api/board/{{boardId}}/item/{{itemId}}' \
--header 'Content-Type: application/json' \
--data-raw '{
    "text": "This is an updated item",
    "color": "green",
    "left": 10,
    "top": 10,
    "width": 100,
    "height": 100
}'
```

### Long poll for changes in a board
```bsh
curl --location --request GET 'http://127.0.0.1:8080/api/board/{{boardId}}/updates/{{version}}'
```
