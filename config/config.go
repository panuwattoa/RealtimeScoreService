package config

import (
	"gamerangkingserver/utils"
)

var (
	ServerType     = utils.GetEnv("SERVER_TYPE", "Development")
	GameDBHost     = utils.GetEnv("DB_HOST", "localhost")
	GameDBPort     = utils.GetEnv("DB_PORT", "3306")
	GameDBName     = utils.GetEnv("DB_NAME", "test")
	GameDBUser     = utils.GetEnv("DB_USERNAME", "test")
	GameDBPassword = utils.GetEnv("DB_PASSWORD", "12345")
	RedisHost      = utils.GetEnv("REDIS_HOST", "127.0.0.1")
	RedisPort      = utils.GetEnv("REDIS_PORT", "6379")
	RedisPassword  = utils.GetEnv("REDIS_PASSWORD", "12345")

	NumLimitRankingData int64  = 100
	WorldRankingKey     string = "WorldRanking"
	EventRankingKey     string = "ScoreKey"
)
