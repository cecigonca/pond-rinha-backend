package main

import (
	"log"
	"net/http"
	"rinha/api"
	"rinha/database"

	"github.com/go-chi/chi/v5"
)

func main() {
	dbpool := database.Connect()
	defer dbpool.Close()

	r := chi.NewRouter()

	r.Post("/clientes/{id}/transacoes", api.TransacaoHandler(dbpool))
	r.Get("/clientes/{id}/extrato", api.ExtratoHandler(dbpool))

	log.Println("API rodando na porta 8080")
	http.ListenAndServe(":8080", r)
}