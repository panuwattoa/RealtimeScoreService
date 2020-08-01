package ranking

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"rangkingserver/config"
	"rangkingserver/storage"
	"rangkingserver/utils"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

// rankingName
// EventType ex. 1 =  PlayCount
// in case name of ranking is 1

var eventCh chan event

type event interface{}

type sendRequestSaveRankingEvent struct {
	responseCh chan<- httpResponse
	info       storage.UserData
}

type sendRequestSaveWorldRankingEvent struct {
	responseCh chan<- httpResponse
	info       storage.UserData
}

type getRankingByEvent struct {
	responseCh      chan<- httpResponse
	info            storage.UserData
	isServerRequest string
}

type initRankingSystemDataEvent struct{}

type clearRankingByEvent struct {
	responseCh chan<- httpResponse
	rankingKey string
}

// eventLoop execute user event queue
func eventLoop() {
	for {
		for event := range eventCh {
			switch ev := event.(type) {
			case initRankingSystemDataEvent:
				handleLoadUserEventData()
			case getRankingByEvent:
				handleGetRankingByEventType(ev.info, ev.responseCh, ev.isServerRequest)
			case clearRankingByEvent:
				handleClearRankingByKey(ev.rankingKey, ev.responseCh)
			}

		}
	}
}

// InitHandler initial eventLoop
func InitHandler() {
	eventCh = make(chan event)
	go eventLoop()
}

// handleProcessRankingByEvent save user statistic via game type
func handleProcessRankingByEvent(info storage.UserData, responseCh chan<- httpResponse) {
	rankingName := info.EventType
	if err := storage.IncreaseScoreDataRedisByRankingData(storage.DataSources, rankingName, utils.ToFloat64(info.Amount), info.UID); err != nil {
		responseCh <- httpResponse{
			statusCode: http.StatusInternalServerError,
			err:        err,
		}
		return
	}

	responseCh <- httpResponse{
		statusCode: http.StatusOK,
		err:        nil,
	}
}

// handleGetRankingByEventType for get score by event name
func handleGetRankingByEventType(info storage.UserData, responseCh chan<- httpResponse, isServerRequest string) {
	rankingName := info.EventType
	var vals []redis.Z
	var err error
	rankingName += info.RankingDuration

	if isServerRequest == "1" {
		if vals, err = storage.GetRedisAllRanking(storage.DataSources, rankingName, true); err != nil {
			responseCh <- httpResponse{
				statusCode: http.StatusInternalServerError,
				err:        err,
			}
			return
		}

	} else {
		if vals, err = storage.GetRedisAllRanking(storage.DataSources, rankingName, false); err != nil {
			responseCh <- httpResponse{
				statusCode: http.StatusInternalServerError,
				err:        err,
			}
			return
		}
	}

	var rankingData []UserResponseData
	rank, err := storage.GetUserRank(storage.DataSources, rankingName, info.UID)
	score, err := storage.GetScoreRedis(storage.DataSources, rankingName, info.UID)
	if err != nil || score <= 0 {
		rank = -1
		score = 0
	}

	if isServerRequest == "0" {
		myUser := UserResponseData{
			UID:   info.UID,
			Rank:  utils.Int64ToString(rank),
			Point: uint64(score),
		}
		rankingData = append(rankingData, myUser)
	}

	for index := 0; index < len(vals); index++ {
		var userData UserResponseData
		userData.UID = fmt.Sprintf("%v", vals[index].Member)
		userData.Rank = fmt.Sprintf("%v", index+1)
		userData.Point = uint64(vals[index].Score)
		rankingData = append(rankingData, userData)
	}
	if jsonData, err := json.Marshal(rankingData); err != nil {
		zap.L().Warn("handleGetRankingByEvent Type parse json error: ", zap.Error(err))
		responseCh <- httpResponse{
			statusCode: http.StatusInternalServerError,
			err:        err,
		}
	} else {
		responseCh <- httpResponse{
			statusCode:  http.StatusOK,
			contentType: "application/json",
			data:        jsonData,
			err:         nil,
		}
	}

}

// handleLoadUserEventData for init server load data from Database fill to redis
func handleLoadUserEventData() {
	_, err := storage.ClearAllRankingByKey(storage.DataSources, config.EventRankingKey)
	if err != nil {
		zap.S().Panic("Error handleLoadUserEventData clear all user data from Redis: ", err)
	}
	zap.L().Info("handleLoadUserGamePlayEventData clear all user data from Redis")

	dailyUserDataList, userDataErr := storage.GetAllUserEventDataFromDB(storage.DataSources)
	if userDataErr != nil {
		zap.L().Panic("GetDailyAllUserGamePlayEventDataFromDB get user data error: ", zap.Error(userDataErr))
	}

	for _, dailyData := range dailyUserDataList {
		rankingName := dailyData.EventType
		if err := storage.IncreaseScoreDataRedisByRankingKey(storage.DataSources, rankingName, utils.ToFloat64(dailyData.Amount), dailyData.UID, config.EventRankingKey); err != nil {
			zap.L().Panic("handleLoadUserGamePlayEventData dailyData increase redis error: ", zap.Error(err))
		}
	}
	zap.L().Info("LoadUserGamePlayEventData Done")
}

// handleClearRankingByKey for clear all data by key
func handleClearRankingByKey(key string, responseCh chan<- httpResponse) {
	if key != "" {
		if _, err := storage.ClearAllRankingByKey(storage.DataSources, key); err != nil {
			responseCh <- httpResponse{
				statusCode: http.StatusInternalServerError,
				err:        err,
			}
			return
		}
	} else {
		responseCh <- httpResponse{
			statusCode: http.StatusInternalServerError,
			err:        errors.New("invalid param"),
		}
		return
	}

	responseCh <- httpResponse{
		statusCode: http.StatusOK,
		err:        nil,
	}
}
