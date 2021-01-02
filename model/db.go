package model

import (
	"database/sql"
	"fmt"
	"time"

	. "github.com/IntouchOpec/base-go-echo/conf"
	"github.com/IntouchOpec/base-go-echo/middleware/cache"
	"github.com/IntouchOpec/base-go-echo/model/orm"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Model struct {
	ID uint
}

var db gorm.DB

// "host=localhost port=5432 user=postgres-dev dbname=dev password=password sslmode=disable"
var dbCacheStore cache.CacheStore

// DB connect data base pastgras.
func DB() *gorm.DB {
	// log.Debugf("Model NewDB")
	newDb, err := newDB()

	if err != nil {
		panic(err)
	}
	newDb.DB().SetMaxOpenConns(20) // Sane default
	newDb.DB().SetMaxIdleConns(0)
	newDb.DB().SetConnMaxLifetime(time.Nanosecond)

	newDb.SetLogger(orm.Logger{})
	newDb.LogMode(Conf.State.IsDevelop)
	// defer newDb.Close()

	return newDb
}

func newDB() (*gorm.DB, error) {
	sqlConnection := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		// "db", "5432", "admin", "todos", "password")
		Conf.DB.Host, Conf.DB.Port, "admin", "line_api", Conf.DB.Pwd)
	// fmt.Printf(sqlConnection)
	db, err := gorm.Open("postgres", sqlConnection)
	if err != nil {
		return nil, err
	}

	return db, nil

}

// Initialize auto migration.
func Initialize() {
	newDb := DB()
	newDb.AutoMigrate(&File{})
	newDb.AutoMigrate(&CustomerType{})
	newDb.AutoMigrate(&ActionLog{})
	newDb.AutoMigrate(&EventLog{})
	newDb.AutoMigrate(&Customer{})
	newDb.AutoMigrate(&Promotion{})
	newDb.AutoMigrate(&Voucher{})
	newDb.AutoMigrate(&Account{})
	newDb.AutoMigrate(&User{})
	newDb.AutoMigrate(&ChatAnswer{})
	newDb.AutoMigrate(&ChatChannel{})
	newDb.AutoMigrate(&Service{})
	newDb.AutoMigrate(&Employee{})
	newDb.AutoMigrate(&EmployeeService{})
	newDb.AutoMigrate(&Booking{})
	newDb.AutoMigrate(&TimeSlot{})
	newDb.AutoMigrate(&Transaction{})
	newDb.AutoMigrate(&LoginRespose{})
	newDb.AutoMigrate(&Setting{})
	newDb.AutoMigrate(&Place{})
	newDb.AutoMigrate(&VoucherCustomer{})
	newDb.AutoMigrate(&BroadcastMessage{})
	newDb.AutoMigrate(&MasterEmployee{})
	newDb.AutoMigrate(&Package{})
	newDb.AutoMigrate(&ServiceItem{})
	newDb.AutoMigrate(&Content{})
	newDb.AutoMigrate(&Report{})
	newDb.AutoMigrate(&OmiseLog{})
	newDb.AutoMigrate(&BookingServiceItem{})
	newDb.AutoMigrate(&BookingPackage{})
	newDb.AutoMigrate(&Payment{})
	newDb.AutoMigrate(&PromotionDetail{})
	newDb.AutoMigrate(&MasterBooking{})
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

func CacheSql(db *sql.DB) *orm.CacheSqlDB {
	return orm.NewCacheDBSql(db, CacheStore(), orm.CacheConf{
		Expire: time.Second * 10,
	})
}
