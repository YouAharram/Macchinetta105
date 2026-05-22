package config

import (
    "log"
    "os"
    "github.com/joho/godotenv"
    "github.com/kelseyhightower/envconfig"
)

type Config struct {
    DBHost     string `envconfig:"DB_HOST"     required:"true"`
    DBPort     string `envconfig:"DB_PORT"     default:"5432"`
    DBUser     string `envconfig:"DB_USER"     required:"true"`
    DBPassword string `envconfig:"DB_PASSWORD" required:"true"`
    DBName     string `envconfig:"DB_NAME"     required:"true"`
    RedisAddr  string `envconfig:"REDIS_ADDR"  default:"localhost:6379"`
    ServerPort string `envconfig:"SERVER_PORT" default:"8080"`
    JWTSecret  string `envconfig:"JWT_SECRET"  required:"true"`
    HMACSecret string `envconfig:"HMAC_SECRET" required:"true"`
}

// DatabaseURL costruisce la stringa di connessione per pgx
func (c *Config) DatabaseURL() string {
    return "postgres://" + c.DBUser + ":" + c.DBPassword +
        "@" + c.DBHost + ":" + c.DBPort + "/" + c.DBName +
        "?sslmode=disable"
}

func Load() *Config {
    // Carica il file .env se esiste (in sviluppo)
    // In produzione le variabili vengono iniettate direttamente dall'ambiente
    if _, err := os.Stat(".env"); err == nil {
        if err := godotenv.Load(); err != nil {
            log.Println("Attenzione: impossibile caricare .env")
        }
    }

    var cfg Config
    if err := envconfig.Process("", &cfg); err != nil {
        log.Fatal("Errore configurazione: ", err)
    }
    return &cfg
}
