package valve

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Queue struct {
	setting  Setting
	limitter *Log
}

func NewQueue() *Queue {
	return &Queue{
		limitter: &Log{
			LastDequeueTime: time.Date(0, 0, 0, 0, 0, 0, 0, time.Local),
		},
	}
}

func (q *Queue) SetLimit(n uint) {
	q.setting.Limit = limit(n)
}

func (q *Queue) Enqueue(m *Message) error {
	now := time.Now()

	m.CreatedAt = &now

	if _, err := db.Exec(`INSERT INTO message (body, expire, request_id) VALUES (?, ?, ?)`, m.Body, m.Expire, m.RequestID); err != nil {
		return err
	}
	return nil
}

func (q *Queue) Dequeue() (*Message, error) {
	tx := db.MustBegin()
	defer tx.Commit()
	ok, err := q.limitter.Block(q.setting.Limit)
	if err != nil {
		q.limitter.Nop()
		return nil, err
	}
	if !ok {
		q.limitter.Nop()
		return nil, errors.New("limit")
	}

	hash := fmt.Sprintf("%X", sha256.Sum256([]byte(strconv.Itoa(rand.Int()))))

	ret := &Message{}
	if _, err := tx.Exec(`UPDATE
			message
		SET
			flag = 1,
			hash = ?
		WHERE flag = 0 ORDER BY id LIMIT 1`, hash); err != nil {
		q.limitter.Nop()
		tx.Rollback()
		return nil, err
	}

	if err := tx.Get(ret, `SELECT id, body, created_at, request_id FROM message WHERE hash = ?`, hash); err != nil {
		q.limitter.Nop()
		tx.Rollback()
		return nil, err
	}

	if _, err := tx.Exec(`DELETE FROM message where id = ?`, ret.ID); err != nil {
		q.limitter.Nop()
		tx.Rollback()
		return nil, err
	}

	return ret, q.limitter.Free()
}
