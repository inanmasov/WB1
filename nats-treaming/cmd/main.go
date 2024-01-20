package main

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/nats-io/stan.go"
)

func main() {

	// Подключение к серверу NATS Streaming
	sc, err := stan.Connect("test-cluster", "publisher-client", stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		log.Fatalf("Ошибка подключения к серверу NATS Streaming: %v", err)
	}
	defer sc.Close()

	// Название канала, в который вы хотите отправить сообщение
	path := "G:\\Стажировка\\nats-treaming\\models\\"
	channel := "myrad"

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range files {
		msg, err := readJSON(path + v.Name())
		if err != nil {
			return
		}

		if err = sc.Publish(channel, msg); err != nil {
			log.Fatalf("Ошибка отправки сообщения: %v", err)
		} else {
			log.Printf("Сообщение успешно опубликовано в канале %s: %s", channel, msg)
			// Для того, чтобы NATS Streaming обработало сообщение
			time.Sleep(1 * time.Second)
		}
	}
}

func readJSON(file string) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("Ошибка открытия файла: %v", err)
	}
	defer f.Close()

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
