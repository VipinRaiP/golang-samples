Ingestion server

$ go run ingestion.go

curl --location --request GET 'http://localhost:8080/ingest' \
--header 'Content-Type: application/json' \
--data '{
    "UserId": "user7",
    "Payload": "message7"
}'

Delivery server

$ go run delivery-server.go 

Delivery client

$ go run delivery-client.go clientId


TODOS

1. Add ack from client side

2. Locks for event list

3. Handling multiple connections of same client

