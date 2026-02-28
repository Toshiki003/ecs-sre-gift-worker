-- 要件4.3: gifts テーブル + send_requests テーブル
-- T-2.4 で実行基盤（マイグレーションツール or 手動適用）を整備する

CREATE TABLE IF NOT EXISTS gifts (
    id         UUID PRIMARY KEY,
    code       VARCHAR NOT NULL UNIQUE,
    amount     INTEGER NOT NULL CHECK (amount > 0),
    message    VARCHAR(500),
    status     VARCHAR(20) NOT NULL DEFAULT 'CREATED',
    last_error TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS send_requests (
    gift_id    UUID NOT NULL REFERENCES gifts(id),
    request_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (gift_id, request_id)
);
