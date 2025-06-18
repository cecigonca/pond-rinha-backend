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

type Extrato struct {
    Saldo struct {
        Total       int       `json:"total"`
        DataExtrato time.Time `json:"data_extrato"`
        Limite      int       `json:"limite"`
    } `json:"saldo"`
    UltimasTransacoes []TransacaoExtrato `json:"ultimas_transacoes"`
}

type TransacaoExtrato struct {
    Valor       int       `json:"valor"`
    Tipo        string    `json:"tipo"`
    Descricao   string    `json:"descricao"`
    RealizadaEm time.Time `json:"realizada_em"`
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

func ExtratoHandler(pool *pgxpool.Pool) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        idStr := chi.URLParam(r, "id")
        id, err := strconv.Atoi(idStr)
        if err != nil {
            http.Error(w, "id inválido", 422)
            return
        }
        var saldo, limite int
        err = pool.QueryRow(context.Background(),
            "SELECT saldo, limite FROM clientes WHERE id = $1", id).
            Scan(&saldo, &limite)
        if err != nil {
            http.Error(w, "cliente não encontrado", 404)
            return
        }
        rows, err := pool.Query(context.Background(),
            `SELECT valor, tipo, descricao, realizada_em
             FROM transacoes
             WHERE cliente_id = $1
             ORDER BY realizada_em DESC
             LIMIT 10`, id)
        if err != nil {
            http.Error(w, "erro ao buscar transações", 500)
            return
        }
        defer rows.Close()
        var transacoes []TransacaoExtrato
        for rows.Next() {
            var t TransacaoExtrato
            if err := rows.Scan(&t.Valor, &t.Tipo, &t.Descricao, &t.RealizadaEm); err != nil {
                http.Error(w, "erro ao processar transações", 500)
                return
            }
            transacoes = append(transacoes, t)
        }
        extrato := Extrato{}
        extrato.Saldo.Total = saldo
        extrato.Saldo.DataExtrato = time.Now().UTC()
        extrato.Saldo.Limite = limite
        extrato.UltimasTransacoes = transacoes
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(extrato)
    }
}