package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/t1nyb0x/tracktaste/handler"
)

func StartServer() {
	r := mux.NewRouter()

	// ルート
	r.HandleFunc("/fetch-artist", handler.FetchArtistHandler).Methods("POST")
	r.HandleFunc("/searchTrack", handler.SearchTrack).Methods("GET")
	// r.HandleFunc("/audio-recommendation", handler.AudioRecommendationHandler).Methods("POST")

	log.Println("サーバー起動中: http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("サーバー起動に失敗しました: %v", err)
	}
	log.Println("サーバー終了")
}
