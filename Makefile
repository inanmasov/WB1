nats-serv:
	cd nats-streaming/ && docker compose up -d
	sleep 5

all: nats-serv