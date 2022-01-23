package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	server2 "seabattle/internal/service"
)

func main() {
	//go func() {
	//	for {
	//		mylogger.PrintGoroutines("interval")
	//		time.Sleep(time.Second * 3)
	//	}
	//}()

	log.SetFlags(log.LstdFlags | log.Llongfile)

	gameServer := server2.NewServer()

	router := mux.NewRouter()
	gameServer.RegisterHandler(router)

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	http.ListenAndServe(":3000", router)
}
