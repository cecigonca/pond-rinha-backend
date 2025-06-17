package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Transacao struct {
	Valor     int    `json:"valor"`
	Tipo      string `json:"tipo"`
	Descricao string `json:"descricao"`
}

func TransacaoHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "id inválido", 422)
			return
		}

		var t Transacao
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			http.Error(w, "body inválido", 422)
			return
		}

		if t.Tipo != "c" && t.Tipo != "d" || t.Valor <= 0 || len(t.Descricao) < 1 || len(t.Descricao) > 10 {
			http.Error(w, "payload inválido", 422)
			return
		}

		tx, _ := pool.Begin(context.Background())
		defer tx.Rollback(context.Background())

		var saldo, limite int
		err = tx.QueryRow(context.Background(),
			"SELECT saldo, limite FROM clientes WHERE id = $1 FOR UPDATE", id).
			Scan(&saldo, &limite)
		if err != nil {
			http.Error(w, "cliente não encontrado", 404)
			return
		}

		novoSaldo := saldo
		if t.Tipo == "d" {
			novoSaldo -= t.Valor
		} else {
			novoSaldo += t.Valor
		}

		if novoSaldo < -limite {
			http.Error(w, "saldo insuficiente", 422)
			return
		}

		_, err = tx.Exec(context.Background(),
			"INSERT INTO transacoes (cliente_id, valor, tipo, descricao) VALUES ($1, $2, $3, $4)",
			id, t.Valor, t.Tipo, t.Descricao)
		if err != nil {
			http.Error(w, "erro ao salvar", 500)
			return
		}

		_, err = tx.Exec(context.Background(),
			"UPDATE clientes SET saldo = $1 WHERE id = $2", novoSaldo, id)
		if err != nil {
			http.Error(w, "erro ao atualizar", 500)
			return
		}

		tx.Commit(context.Background())

		json.NewEncoder(w).Encode(map[string]int{
			"limite": limite,
			"saldo":  novoSaldo,
		})
	}
}