package storage

import (
	"fmt"
	"gamerangkingserver/config"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

// DataSources contain DB connection and redis client
var DataSources *DataSource

// UserData is require from clients when adding user ranking
type UserData struct {
	UID             string `json:"uid"`
	Name            string `json:"name"`
	EventType       string `json:"even_type"`
	Amount          string `json:"amount"`
	RankingDuration string `json:"ranking_duration"`
}

// DataSource struct contain DB connection and RedisClient
type DataSource struct {
	DataSourceName string
	RedisClient    *redis.Client
}

// Close close all connection
func (ds *DataSource) Close() {
	// ds.DB.Close()
	ds.RedisClient.Close()
}

// NewDataSource for initial program
func NewDataSource() *DataSource {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		config.GameDBUser,
		config.GameDBPassword,
		config.GameDBHost,
		config.GameDBPort,
		config.GameDBName,
	)

	zap.L().Info("use db: ", zap.String("data source", dataSourceName))

	// db, err := sql.Open("mysql", dataSourceName)
	// if err != nil {
	// 	zap.L().Panic("cannot open connection", zap.String("source", dataSourceName), zap.Error(err))
	// }
	// if err = db.Ping(); err != nil {
	// 	zap.L().Panic("cannot ping connection", zap.String("source", dataSourceName), zap.Error(err))
	// }

	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.RedisHost + ":" + config.RedisPort,
		Password: config.RedisPassword,
		DB:       0,
	})

	pong, err := redisClient.Ping().Result()
	zap.L().Info("redis status: ", zap.String("pong", pong), zap.Error(err))
	if err != nil {
		zap.L().Fatal("status: ", zap.Error(err))
	}
	return &DataSource{
		DataSourceName: dataSourceName,
		RedisClient:    redisClient,
	}
}
