package main

import (
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	server2 "seabattle/internal/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln(err)
	}

	level, err := log.ParseLevel(os.Getenv("SEABATTLE_LOG_LEVEL"))
	if err != nil {
		log.Fatalln(err)
	}

	log.SetOutput(os.Stdout)
	log.SetLevel(level)

	log.Debugln("Log level is", level)

	gameServer := server2.NewServer()

	router := mux.NewRouter()
	gameServer.RegisterHandler(router)

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	addr := os.Getenv("SEABATTLE_ADDR")
	log.Infoln("Starting server on address", addr)
	http.ListenAndServe(addr, router)
}
