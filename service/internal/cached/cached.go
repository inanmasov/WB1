package cached

import (
	"log"
	"sync"
	"time"

	database "example.com/service/service/internal/database"
	"github.com/patrickmn/go-cache"
)

// CacheManager представляет собой менеджер кэша.
type CacheManager struct {
	Cache *cache.Cache
	Mu    sync.Mutex
}

var GlobalCacheManager *CacheManager

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
				GlobalCacheManager.Cache.Set(key, order[i], -1)
			}
		}
	}
	return true
}

// NewCache создает новый CacheManager.
func NewCache(defaultExpiration, cleanupInterval time.Duration) *CacheManager {
	return &CacheManager{
		Cache: cache.New(defaultExpiration, cleanupInterval),
	}
}

// Set устанавливает значение в кэше с заданным ключом и сроком действия.
func (cm *CacheManager) Set(key string, value interface{}, expiration time.Duration) {
	cm.Mu.Lock()
	defer cm.Mu.Unlock()
	cm.Cache.Set(key, value, expiration)
}

// Get получает значение из кэша по ключу.
func (cm *CacheManager) Get(key string) (interface{}, bool) {
	cm.Mu.Lock()
	defer cm.Mu.Unlock()
	return cm.Cache.Get(key)
}

// Delete удаляет значение из кэша по ключу.
func (cm *CacheManager) Delete(key string) {
	cm.Mu.Lock()
	defer cm.Mu.Unlock()
	cm.Cache.Delete(key)
}

// Flush очищает весь кэш.
func (cm *CacheManager) Flush() {
	cm.Mu.Lock()
	defer cm.Mu.Unlock()
	cm.Cache.Flush()
}
