package model

import (
	"fmt"
	"time"

	. "github.com/IntouchOpec/base-go-echo/conf"
	"github.com/hb-go/echo-web/middleware/cache"
	"github.com/hb-go/gorm"

	// "github.com/hb-go/gorm"

	"github.com/IntouchOpec/base-go-echo/model/orm"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/gommon/log"
)

// "host=localhost port=5432 user=postgres-dev dbname=dev password=password sslmode=disable"
var db *gorm.DB
var dbCacheStore cache.CacheStore

// DB connect data base pastgras.
func DB() *gorm.DB {
	if db == nil {
		log.Debugf("Model NewDB")

		newDb, err := newDB()

		if err != nil {
			panic(err)
		}

		newDb.DB().SetMaxIdleConns(10)
		newDb.DB().SetMaxOpenConns(100)

		newDb.SetLogger(orm.Logger{})
		newDb.LogMode(true)
		db = newDb
	}

	return db
}

func newDB() (*gorm.DB, error) {

	sqlConnection := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		Conf.DB.Host, Conf.DB.Port, Conf.DB.UserName, Conf.DB.Name, Conf.DB.Pwd)
	db, err := gorm.Open("postgres", sqlConnection)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Initialize auto migration.
func Initialize() {
	newDb, err := newDB()

	if err != nil {
		panic(err)
	}

	newDb.AutoMigrate(&Account{})
	newDb.AutoMigrate(&User{})
	newDb.AutoMigrate(&Customer{})
	newDb.AutoMigrate(&ChatAnswer{})
	newDb.AutoMigrate(&ChatChannel{})
	newDb.AutoMigrate(&ChatRequest{})
	newDb.AutoMigrate(&EventLog{})
	newDb.AutoMigrate(&Product{})
	newDb.AutoMigrate(&SubProduct{})
	newDb.AutoMigrate(&Booking{})
	newDb.AutoMigrate(&Promotion{})
	newDb.AutoMigrate(&Setting{})
	newDb.AutoMigrate(&Account{})
	newDb.AutoMigrate(&LoginRespose{})
	newDb.AutoMigrate(&TemplateSocial{})
	newDb.AutoMigrate(&TemplateSocialDetail{})
	newDb.AutoMigrate(&KeyTemplate{})
	newDb.AutoMigrate(&ActionLog{})
	// Foreign Key Account Table.
	newDb.Model(&User{}).AddForeignKey("account_id", "accounts(id)", "CASCADE", "RESTRICT")

	// Booking Foreign.
	newDb.Model(&Booking{}).AddForeignKey("customer_id", "customers(id)", "CASCADE", "RESTRICT")
	newDb.Model(&Booking{}).AddForeignKey("sub_product_id", "sub_products(id)", "CASCADE", "RESTRICT")
	newDb.Model(&Booking{}).AddForeignKey("chat_channel_id", "chat_channels(id)", "CASCADE", "RESTRICT")

	// ChatRequest Foreign.
	newDb.Model(&ChatRequest{}).AddForeignKey("chat_answer_id", "chat_answers(id)", "CASCADE", "RESTRICT")

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
