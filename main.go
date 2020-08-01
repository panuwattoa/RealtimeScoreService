package main

import (
	"fmt"
	"log"
	"net/http"
	"rangkingserver/config"
	"rangkingserver/ranking"
	"rangkingserver/storage"

	"go.uber.org/zap"
)

func main() {
	var logger *zap.Logger
	var err error

	switch config.ServerType {
	case "Production":
		logger, err = zap.NewProduction()
		if err != nil {
			log.Fatal("cannot initialize logger: ", err)
		}
	case "Development":
		logger, err = zap.NewDevelopment()
		if err != nil {
			log.Fatal("cannot initialize logger: ", err)
		}
	}
	zap.ReplaceGlobals(logger)
	zap.L().Debug("start ranking server: ", zap.String("server-type", config.ServerType))
	defer logger.Sync()

	storage.DataSources = storage.NewDataSource()
	defer storage.DataSources.Close()
	// init event loop
	ranking.InitHandler()
	ranking.InitRankingSystemData()
	// http handle

	http.Handle("/saveGamePlayRanking", withCors(ranking.SaveRankingByEvent))
	http.Handle("/getRankingByEvent", withCors(ranking.GetRankingByEvent))
	http.Handle("/clearRankingByKey", withCors(ranking.ClearRankingByKey))
	switch config.ServerType {
	case "Production":
		log.Fatal(http.ListenAndServeTLS("0.0.0.0:8444", "certs/fullchain.pem", "certs/privkey.pem", nil))
	case "Development":
		log.Fatal(http.ListenAndServeTLS("0.0.0.0:8444", "certs/fullchain_ds.pem", "certs/privkey_ds.pem", nil))
	}
}

func handle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Please add choices before spin.")
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
}

func withCors(handler func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		addCors(&w)
		handler(w, r)
	})
}

func addCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
