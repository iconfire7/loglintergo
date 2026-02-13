package main

import (
	"fmt"
	"log/slog"
)

type config struct {
	user string
	host string
	port string
}

func main() {
	lg := slog.Logger{}
	lg.Info("add logger")

	// cfg := config.Load()
	lg.Info("load config")
	lg.Debug(fmt.Sprintf("%+v", config{}))

	// db := repo.New(cfg)
	passwd := "postgres"
	lg.Info("database connected")
	lg.Debug("database password" + passwd)

	//http.ListenAndServe("8080", mux)
	slog.Info("Server started in port 8080")
}
