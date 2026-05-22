package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/jackc/pgx/v5/pgxpool"

    "github.com/tuonome/vending-server/internal/config"
)

type APIResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, payload APIResponse) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(payload)
}

func main() {
    // 1. Carica configurazione
    cfg := config.Load()

    // 2. Connettiti a Postgres
    // pgxpool è un pool di connessioni: invece di aprire/chiudere
    // una connessione per ogni richiesta, mantiene un gruppo pronto
    db, err := pgxpool.New(context.Background(), cfg.DatabaseURL())
    if err != nil {
        log.Fatal("Impossibile connettersi al database: ", err)
    }
    defer db.Close() // chiudi il pool quando il server si spegne

    // Verifica che la connessione funzioni davvero
    if err := db.Ping(context.Background()); err != nil {
        log.Fatal("Database non raggiungibile: ", err)
    }
    log.Println("✓ Connesso a PostgreSQL")

    // 3. Router
    r := chi.NewRouter()
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)

    r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
        // Qui verifichiamo anche il DB, non solo che il server risponda
        if err := db.Ping(r.Context()); err != nil {
            writeJSON(w, http.StatusServiceUnavailable, APIResponse{
                Error: "database non raggiungibile",
            })
            return
        }
        writeJSON(w, http.StatusOK, APIResponse{
            Success: true,
            Data:    map[string]string{"status": "ok"},
        })
    })

    log.Printf("Server in ascolto su :%s\n", cfg.ServerPort)
    if err := http.ListenAndServe(":"+cfg.ServerPort, r); err != nil {
        log.Fatal(err)
    }
}
