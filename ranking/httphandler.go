package ranking

import (
	"encoding/json"
	"fmt"
	"gamerangkingserver/storage"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"
)

type httpResponse struct {
	statusCode  int
	contentType string
	data        []byte
	err         error
}

type userBody struct {
	UID           string `json:"uid"`
	Name          string `json:"name"`
	Level         string `json:"level"`
	ProfilePicID  string `json:"profile_pic_id"`
	AccountTypeID string `json:"account_type_id"`
	CrownID       string `json:"crown_id"`
	FrameID       string `json:"frame_id"`
	EventType     string `json:"event_type"`
	GameMode      string `json:"game_mode"`
	SubTitle      string `json:"sub_title"`
	Amount        string `json:"amount"`
	RoomPrivacy   string `json:"room_privacy"`
}

// InitRankingSystemData prepare data when starter
func InitRankingSystemData() {
	eventCh <- initRankingSystemDataEvent{}
}

// SaveRankingByEvent save rank via event type
func SaveRankingByEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		zap.L().Warn("SaveRankingByEvent method is not POST")
		http.Error(w, "SaveRankingByEvent method is not POST", http.StatusMethodNotAllowed)
		return
	}
	receiveResponseCh := make(chan httpResponse)
	var info userBody
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "err Body %v", err)
		return
	}
	err = json.Unmarshal(reqBody, &info)
	if err != nil {
		fmt.Fprintf(w, "err Unmarshal %v", err)
		return
	}

	eventCh <- sendRequestSaveRankingEvent{
		responseCh: receiveResponseCh,
		info: storage.UserData{
			UID:       info.UID,
			EventType: info.EventType,
			Amount:    info.Amount,
			Name:      info.Name,
		},
	}

	responseData := <-receiveResponseCh
	if responseData.err != nil {
		http.Error(w, responseData.err.Error(), responseData.statusCode)
		return
	}

	w.WriteHeader(responseData.statusCode)

	if responseData.contentType != "" {
		w.Header().Set("Content-Type", responseData.contentType)
	}
	if len(responseData.data) > 0 {
		w.Write(responseData.data)
	}

}

// GetRankingByEvent get ranking by event type gameMode and subtitle rate
func GetRankingByEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		zap.L().Warn("SaveWorldRanking method is not GET")
		http.Error(w, "SaveWorldRanking method is not GET", http.StatusMethodNotAllowed)
		return
	}
	receiveResponseCh := make(chan httpResponse)

	UID := r.FormValue("uid")
	eventType := r.FormValue("eventType")
	gameMode := r.FormValue("gameMode")
	subtitle := r.FormValue("subTitle")
	rankingDuration := r.FormValue("rankingDuration")
	serverRequest := r.FormValue("isServerRequest")

	if eventType == "" || gameMode == "" || subtitle == "" || rankingDuration == "" || serverRequest == "" {
		http.Error(w, "Invalid param", http.StatusNoContent)
		return
	}
	// eventType ex. 1 =  PlayCount
	// serverRequest ex.  1 or 0
	// in case name of ranking is 11

	eventCh <- getRankingByEvent{
		responseCh: receiveResponseCh,
		info: storage.UserData{
			UID:             UID,
			EventType:       eventType,
			RankingDuration: rankingDuration,
		},
		isServerRequest: serverRequest,
	}

	responseData := <-receiveResponseCh
	if responseData.err != nil {
		http.Error(w, responseData.err.Error(), responseData.statusCode)
		return
	}

	w.WriteHeader(responseData.statusCode)

	if responseData.contentType != "" {
		w.Header().Set("Content-Type", responseData.contentType)
	}
	if len(responseData.data) > 0 {
		w.Write(responseData.data)
	}

}

// ClearRankingByKey clear ranking by key ex. daily or weekly
func ClearRankingByKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		zap.L().Warn("ClearRankingBykey method is not GET")
		http.Error(w, "ClearRankingBykey method is not GET", http.StatusMethodNotAllowed)
		return
	}
	key := r.FormValue("rankingkey")
	if key == "" {
		http.Error(w, "Invalid param", http.StatusNoContent)
		return
	}
	receiveResponseCh := make(chan httpResponse)

	eventCh <- clearRankingByEvent{
		responseCh: receiveResponseCh,
		rankingKey: key,
	}

	responseData := <-receiveResponseCh
	if responseData.err != nil {
		http.Error(w, responseData.err.Error(), responseData.statusCode)
		return
	}

	w.WriteHeader(responseData.statusCode)

	if responseData.contentType != "" {
		w.Header().Set("Content-Type", responseData.contentType)
	}
	if len(responseData.data) > 0 {
		w.Write(responseData.data)
	}
}
