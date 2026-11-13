package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	isDevelop := flag.Bool("isDevelop", true, "if the environment is development")
	backupFile := flag.String("backupFilePath", "./testing.json", "where the backup file is stored")
	flag.Parse()

	addr := "localhost:8080"
	if !*isDevelop {
		addr = "1.1.1.1:8080"
	}
	log.Println("Startup - Listening on:", addr)
	mux := http.NewServeMux()
	handler := addLoggingMiddleware(mux)

	ib := &jsonFileBackup[ResItemCreate]{*backupFile}
	d, err := NewDependencies(ib)
	if err != nil {
		log.Fatalln("Error creating dependencies:", err)
	}
	setupV1Endpoints(mux, d, addr)
	log.Fatalln(http.ListenAndServe(addr, handler))
}
