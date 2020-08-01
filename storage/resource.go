package storage

import (
	"database/sql"
	"rangkingserver/config"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

//----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
// Integrate with DB
//----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

// GetAllUserEventDataFromDB get daily data from game database `play_event` for store in redis
func GetAllUserEventDataFromDB(ds *DataSource) ([]UserData, error) {
	var userDataList []UserData
	db, err := sql.Open("mysql", ds.DataSourceName)
	if err != nil {
		zap.L().Panic("cannot open connection", zap.String("source", ds.DataSourceName), zap.Error(err))
	}
	defer db.Close()
	rows, err := db.Query("SELECT event_type, uid, sum(value) FROM `play_event`  GROUP by uid,event_type")
	if err != nil {
		return userDataList, err
	}

	defer rows.Close()

	for rows.Next() {
		userData := UserData{}
		err := rows.Scan(&userData.EventType, &userData.UID, &userData.Amount)
		if err != nil {
			return userDataList, err
		}
		userDataList = append(userDataList, userData)
	}

	if err := rows.Err(); err != nil {
		return userDataList, err
	}

	return userDataList, nil
}

// GetAllUserStatisticFromDB get user statistic data from game database `user_dummy` for store in redis
// func GetAllUserStatisticFromDB(ds *DataSource) ([]UserStatistic, error) {
// 	var userDataList []UserStatistic
// 	db, err := sqlx.Open("mysql", ds.DataSourceName)
// 	if err != nil {
// 		zap.L().Panic("cannot open connection", zap.String("source", ds.DataSourceName), zap.Error(err))
// 	}
// 	defer db.Close()
// 	rows, err := db.Queryx("SELECT user.uid as uid, `knock`, `knock_color`, `dark_knock`, `dark_knock_color`, `dai_ngo`, `win_no_1`, user.name as name, user.level as level, profile_pic_id, platform +0 as platform, total_chip_gain  FROM `log_user_dummy` JOIN user ON log_user_dummy.uid = user.uid")
// 	if err != nil {
// 		return userDataList, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		userData := UserStatistic{}

// 		err := rows.StructScan(&userData)
// 		if err != nil {
// 			zap.S().Info("Load StructScan GetAllUserStatisticFromDB ", err)
// 			return userDataList, err
// 		}

// 		userDataList = append(userDataList, userData)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return userDataList, err
// 	}
// 	return userDataList, nil

// }

// // GetUserProfileFromDB get user profile item data from DB
// func GetUserProfileFromDB(ds *DataSource) (map[string]string, error) {

// }

//----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
// Integrate with Redis
//----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

// rankingName = eventtype + GameMode + SubTitle

// IncreaseScoreDataRedisByRankingData ZIncrBy increase value in redis and get ranking type
func IncreaseScoreDataRedisByRankingData(ds *DataSource, rankingName string, score float64, uid string) error {
	_, err := ds.RedisClient.ZIncrBy(rankingName+config.EventRankingKey, score, uid).Result()
	_, err = ds.RedisClient.SAdd(config.EventRankingKey, rankingName+config.EventRankingKey).Result()
	return err
}

// IncreaseScoreDataRedisByRankingKey increase value by ranking key
func IncreaseScoreDataRedisByRankingKey(ds *DataSource, rankingName string, score float64, uid string, rankingKey string) error {
	_, err := ds.RedisClient.ZIncrBy(rankingName+rankingKey, score, uid).Result()
	_, err = ds.RedisClient.SAdd(rankingKey, rankingName+rankingKey).Result()
	return err
}

// IncreaseScoreDataWorldRankingRedis ZIncrBy increase value in redis Hall of fame
func IncreaseScoreDataWorldRankingRedis(ds *DataSource, rankingName string, score float64, uid string) error {
	_, err := ds.RedisClient.ZIncrBy(rankingName, score, uid).Result()
	_, err = ds.RedisClient.SAdd(config.WorldRankingKey, rankingName).Result()

	return err
}

// SetScoreDataRedis value by score
func SetScoreDataRedis(ds *DataSource, rankingName string, score float64, uid string) error {
	_, err := ds.RedisClient.ZAdd(rankingName, redis.Z{
		Score:  score,
		Member: uid,
	}).Result()

	return err
}

// DeleteRedis  delete value in redis via ranking name
func DeleteRedis(ds *DataSource, rankingName string, score float64, uid string) error {
	_, err := ds.RedisClient.Del(rankingName).Result()
	return err
}

// GetScoreRedis get user score via rankingName
func GetScoreRedis(ds *DataSource, rankingName string, uid string) (float64, error) {
	score, err := ds.RedisClient.ZScore(rankingName, uid).Result()
	return score, err
}

// GetUserRank get user rank via rankingName
func GetUserRank(ds *DataSource, rankingName string, uid string) (int64, error) {
	val, err := ds.RedisClient.ZRevRank(rankingName, uid).Result()
	return val + 1, err
}

// GetRedisAllRanking get all data and can get data limit by limit
func GetRedisAllRanking(ds *DataSource, rankingName string, isGetNoLimit bool) ([]redis.Z, error) {

	var vals []redis.Z
	var err error
	if isGetNoLimit {
		vals, err = ds.RedisClient.ZRevRangeByScoreWithScores(rankingName, redis.ZRangeBy{
			Min:    "0",
			Max:    "+inf",
			Offset: 0,
		}).Result()
	} else {
		vals, err = ds.RedisClient.ZRevRangeByScoreWithScores(rankingName, redis.ZRangeBy{
			Min:    "1",
			Max:    "+inf",
			Offset: 0,
			Count:  config.NumLimitRankingData,
		}).Result()
	}

	return vals, err

}

// ClearAllRankingByKey clear type daily ranking
func ClearAllRankingByKey(ds *DataSource, key string) (int64, error) {
	listKey, err := ds.RedisClient.SMembers(key).Result()
	for _, key := range listKey {
		ds.RedisClient.Del(key).Result()
	}
	result, err := ds.RedisClient.Del(key).Result()
	return result, err
}

// ClearAllWorldRankingDataRedis clear all world ranking data store in redis
func ClearAllWorldRankingDataRedis(ds *DataSource) (int64, error) {
	listKey, err := ds.RedisClient.SMembers(config.WorldRankingKey).Result()
	for _, key := range listKey {
		ds.RedisClient.Del(key).Result()
	}
	result, err := ds.RedisClient.Del(config.WorldRankingKey).Result()
	return result, err
}

// GetAllKeyRankingByDuraion get all key by member
func GetAllKeyRankingByDuraion(ds *DataSource, durationKey string) ([]string, error) {
	listKey, err := ds.RedisClient.SMembers(durationKey).Result()
	return listKey, err
}
