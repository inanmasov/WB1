nats-serv:
	cd nats-streaming/ && docker compose up -d

publish:
	cd nats-streaming/ && go run cmd/main.go

server:
	cd service/ && go run cmd/main.go

all: nats-serv publish server

clean:
	cd nats-streaming/ && docker compose down