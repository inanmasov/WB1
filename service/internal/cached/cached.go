package cached

import (
	"log"
	"time"

	database "example.com/service/service/internal/database"
	"github.com/patrickmn/go-cache"
)

var GlobalCacheManager *cache.Cache

func InitCacheDB() bool {
	if GlobalCacheManager == nil {
		GlobalCacheManager = NewCache(0, 0)

		db, err := database.Initialize()
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		order, err := db.GetAllOrdersDatabase()
		if err != nil {
			log.Fatal(err)
		}

		if len(order) == 0 {
			return false
		} else {
			for i := 0; i < len(order); i++ {
				key := order[i].OrderUID
				GlobalCacheManager.Set(key, order[i], -1)
			}
		}
	}
	return true
}

// NewCache создает новый CacheManager.
func NewCache(defaultExpiration, cleanupInterval time.Duration) *cache.Cache {
	return cache.New(defaultExpiration, cleanupInterval)
}
