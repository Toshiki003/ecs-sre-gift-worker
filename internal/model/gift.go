package model

import "time"

// GiftStatus はギフトの状態を表す。
// 要件4.2: CREATED → QUEUED → SENT / FAILED
type GiftStatus string

const (
	GiftStatusCreated GiftStatus = "CREATED"
	GiftStatusQueued  GiftStatus = "QUEUED"
	GiftStatusSent    GiftStatus = "SENT"
	GiftStatusFailed  GiftStatus = "FAILED"
)

// Gift はギフトエンティティ（要件4.3: gifts テーブル）。
type Gift struct {
	ID        string     `json:"gift_id"`
	Code      string     `json:"code"`
	Amount    int        `json:"amount"`
	Message   string     `json:"message,omitempty"`
	Status    GiftStatus `json:"status"`
	LastError string     `json:"last_error,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// SendRequest はギフト送信リクエスト（要件4.3: send_requests テーブル）。
type SendRequest struct {
	GiftID    string    `json:"gift_id"`
	RequestID string    `json:"request_id"`
	CreatedAt time.Time `json:"created_at"`
}
