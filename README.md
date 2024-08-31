# Product Service in Go

## Running the server
1. Clone this repository (optional install task https://taskfile.dev)

2. Mount the repository & run this command to install dependencies
```
go mod tidy
```

3. Migrate data
```
migrate -path db/migrations -database "postgres://${SHOPEEFUN_POSTGRES_USER}:${SHOPEEFUN_POSTGRES_PASSWORD}@${SHOPEEFUN_POSTGRES_HOST}:${SHOPEEFUN_POSTGRES_PORT}/${SHOPEEFUN_POSTGRES_DB}?sslmode=${SHOPEEFUN_POSTGRES_SSL_MODE}" up
```

4. Run server
```
go run ./cmd/bin/main.go
```

5. Server will be running on `localhost:4000`

## Some example from API

1. POST categories
```
curl --location 'http://localhost:4000/products/category' \
--header 'X-USER-ID: 84095313-f3dc-4529-b869-24bb5c77c1a4' \
--header 'Content-Type: application/json' \
--data '{
    "name": "elektronik"
}'
```

2. GET All categories
```
curl --location 'http://localhost:4000/products/categories?paginate=5&page=1' \
--header 'X-USER-ID: 84095313-f3dc-4529-b869-24bb5c77c1a4'
```

3. POST product
```
{
 curl --location 'http://localhost:4000/products' \
--header 'X-USER-ID: 84095313-f3dc-4529-b869-24bb5c77c1a4' \
--header 'Content-Type: application/json' \
--data '{
    "category_id": "c97081c5-6ed3-4649-b7ff-6113ecc09a4e",
    "shop_id": "9aa56858-974a-4c5a-9aa8-0fcce63430f7",
    "name": "Gamis mahal",
    "price": 30000,
    "stock": 20,
    "description": "test", 
    "image_url": "https://fastly.picsum.photos/id/607/200/300.jpg?hmac=ZEvzqI62NudR3rgqTkRZzFnlEeOt9z-b_i8VdLoTgoI", 
    "brand" : "Adidas"
}'
}
```

4. GET All products
```
curl --location 'http://localhost:4000/products?paginate=5&page=1&name=Gamis%20mahal%202' \
--header 'X-USER-ID: 84095313-f3dc-4529-b869-24bb5c77c1a4'
```


## ERD
This ERD describes how this dbserver works.
![Server-ERD](?raw=true)
 