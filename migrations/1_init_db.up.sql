-- USERS: gli utenti dell'app
CREATE TABLE users (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email       TEXT NOT NULL UNIQUE,
    password    TEXT NOT NULL,          -- hash bcrypt, mai la password in chiaro
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- WALLETS: il credito di ogni utente (tabella separata per chiarezza)
CREATE TABLE wallets (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    balance     NUMERIC(10, 2) NOT NULL DEFAULT 0.00
                    CHECK (balance >= 0),  -- il saldo non può andare sotto zero
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- DEVICES: le macchinette fisiche
CREATE TABLE devices (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT NOT NULL,           -- es. "Macchinetta Piano 2"
    api_key     TEXT NOT NULL UNIQUE,    -- chiave con cui la macchina si autentica
    location    TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- PRODUCTS: i prodotti disponibili
CREATE TABLE products (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT NOT NULL,
    price       NUMERIC(10, 2) NOT NULL CHECK (price > 0),
    stock       INTEGER NOT NULL DEFAULT 0 CHECK (stock >= 0),
    device_id   UUID NOT NULL REFERENCES devices(id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- OTP_TOKENS: i token monouso generati quando l'utente sceglie un prodotto
CREATE TABLE otp_tokens (
    token       TEXT PRIMARY KEY,        -- il token firmato con HMAC
    user_id     UUID NOT NULL REFERENCES users(id),
    product_id  UUID NOT NULL REFERENCES products(id),
    price       NUMERIC(10, 2) NOT NULL, -- prezzo al momento della generazione
    used        BOOLEAN NOT NULL DEFAULT FALSE,
    expires_at  TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- TRANSACTIONS: audit log immutabile di tutto quello che è successo
CREATE TABLE transactions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id),
    product_id  UUID NOT NULL REFERENCES products(id),
    device_id   UUID NOT NULL REFERENCES devices(id),
    token_hash  TEXT NOT NULL,           -- hash del token, non il token stesso
    amount      NUMERIC(10, 2) NOT NULL,
    status      TEXT NOT NULL,           -- 'success', 'failed', 'refunded'
    error_msg   TEXT,                    -- popolato solo in caso di errore
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- INDEX: velocizzano le query più frequenti
CREATE INDEX idx_otp_tokens_user_id    ON otp_tokens(user_id);
CREATE INDEX idx_transactions_user_id  ON transactions(user_id);
CREATE INDEX idx_products_device_id    ON products(device_id);
