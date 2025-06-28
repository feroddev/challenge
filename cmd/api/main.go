package main

import (
	"fmt"
	"log"
	"net/http"

	"challengeV3/api"
	"challengeV3/configs"
)

func main() {
	config := configs.NewConfig()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Bem-vindo à API do projeto!")
	})

	http.HandleFunc("/health", api.HealthCheck)

	fmt.Printf("Servidor iniciado na porta %s\n", config.Server.Port)
	log.Fatal(http.ListenAndServe(":"+config.Server.Port, nil))
}
