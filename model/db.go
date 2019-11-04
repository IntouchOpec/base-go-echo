package model

import (
	"fmt"
	"time"

	. "github.com/IntouchOpec/base-go-echo/conf"
	"github.com/IntouchOpec/base-go-echo/middleware/cache"
	"github.com/IntouchOpec/base-go-echo/model/orm"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/gommon/log"
)

type Model struct {
	ID uint
}

var db gorm.DB

// "host=localhost port=5432 user=postgres-dev dbname=dev password=password sslmode=disable"
var dbCacheStore cache.CacheStore

// DB connect data base pastgras.
func DB() *gorm.DB {
	log.Debugf("Model NewDB")

	newDb, err := newDB()

	if err != nil {
		panic(err)
	}

	newDb.DB().SetMaxIdleConns(10)
	newDb.DB().SetMaxOpenConns(100)

	newDb.SetLogger(orm.Logger{})
	newDb.LogMode(true)
	// defer newDb.Close()

	return newDb
}

func newDB() (*gorm.DB, error) {

	sqlConnection := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		Conf.DB.Host, Conf.DB.Port, "admin", "todos", Conf.DB.Pwd)
	db, err := gorm.Open("postgres", sqlConnection)
	if err != nil {
		return nil, err
	}

	return db, nil

}

type TestModel struct {
	Namw string
}

// Initialize auto migration.
func Initialize() {
	newDb := DB()
	newDb.AutoMigrate(&CustomerTpye{})
	newDb.AutoMigrate(&Customer{})
	newDb.AutoMigrate(&Promotion{})
	newDb.AutoMigrate(&Account{})
	newDb.AutoMigrate(&User{})
	newDb.AutoMigrate(&ChatAnswer{})
	newDb.AutoMigrate(&ChatChannel{})
	newDb.AutoMigrate(&ChatRequest{})
	newDb.AutoMigrate(&EventLog{})
	newDb.AutoMigrate(&Service{})
	newDb.AutoMigrate(&ServiceSlot{})
	newDb.AutoMigrate(&Booking{})
	newDb.AutoMigrate(&Account{})
	newDb.AutoMigrate(&LoginRespose{})
	newDb.AutoMigrate(&ActionLog{})
	newDb.AutoMigrate(&Setting{})
	newDb.AutoMigrate(&Place{})
}

// CacheStore use cache MEMCACHED or REDIS.
func CacheStore() cache.CacheStore {
	if dbCacheStore == nil {
		switch Conf.CacheStore {
		case MEMCACHED:
			dbCacheStore = cache.NewMemcachedStore([]string{Conf.Memcached.Server}, time.Hour)
		case REDIS:
			dbCacheStore = cache.NewRedisCache(Conf.Redis.Server, Conf.Redis.Pwd, time.Hour)
		default:
			dbCacheStore = cache.NewInMemoryStore(time.Hour)
		}
	}

	return dbCacheStore
}

// Cache config expire time
func Cache(db *gorm.DB) *orm.CacheDB {
	return orm.NewCacheDB(db, CacheStore(), orm.CacheConf{
		Expire: time.Second * 10,
	})
}
