package valve

import "time"

type Message struct {
	ID        int        `json:"id" db:"id"`
	Body      string     `json:"body" db:"body"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	Expire    *time.Time `json:"expire" db:"expire"`
	RequestID string     `json:"request_id" db:"request_id"`
}
