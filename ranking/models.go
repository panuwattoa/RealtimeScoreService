package ranking


// UserResponseData is response data to client
type UserResponseData struct {
	UID   string `json:"uid"`
	Name  string `json:"name"`
	Rank  string `json:"rank"`
	Point uint64 `json:"point"`
}

// UserRankData  response data for get reward
type UserRankData struct {
	UID         string `json:"uid"`
	Rank        string `json:"rank"`
	Point       uint32 `json:"point"`
	RankingName string `json:"ranking_name"`
}
