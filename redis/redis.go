package redis

import (
	"encoding/json"
	"labxsubs/model"
	"log"
	"os"
	"sync"

	"github.com/go-redis/redis/v7"
)

var client RedisClient
var mutex sync.Mutex

var CreatedSupplierChannel map[string]chan *model.Supplier

type RedisClient struct {
	*redis.Client
}

func New() RedisClient {
	CreatedSupplierChannel = map[string]chan *model.Supplier{}

	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	pass := os.Getenv("REDIS_PASS")
	addr := host + ":" + port
	c := RedisClient{
		redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: pass,
			DB:       0,
		}),
	}

	startSubscribingRedis()

	return c
}

func SetClient(rc RedisClient) {
	client = rc
}

func GetClient() RedisClient {
	return client
}

func startSubscribingRedis() {
	go func() {
		pubsub := client.Subscribe("suppliers")
		defer pubsub.Close()

		for {
			_, err := pubsub.Receive()
			if err != nil {
				log.Println(err)
			}

			controlCh := pubsub.Channel()

			for msg := range controlCh {
				m := model.Supplier{}
				if err := json.Unmarshal([]byte(msg.Payload), &m); err != nil {
					log.Println(err)
					continue
				}

				mutex.Lock()
				for _, ch := range CreatedSupplierChannel {
					ch <- &m
				}
				mutex.Unlock()
			}
		}
	}()
}
