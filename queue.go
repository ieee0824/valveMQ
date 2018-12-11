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

	if _, err := db.Exec(`INSERT INTO message (body, expire) VALUES (?, ?)`, m.Body, m.Expire); err != nil {
		return err
	}
	return nil
}

func (q *Queue) Dequeue() (*Message, error) {
	if err := q.limitter.Block(); err != nil {
		return nil, err
	}
	now := time.Now()
	if now.Sub(q.limitter.LastDequeueTime) < q.setting.Limit.DqSpan() {
		if err := q.limitter.Nop(); err != nil {
			return nil, err
		}
		return nil, errors.New("It took band limitation.")
	}
	tx := db.MustBegin()

	hash := fmt.Sprintf("%X", sha256.Sum256([]byte(strconv.Itoa(rand.Int()))))

	ret := &Message{}
	if _, err := tx.Exec(`UPDATE
			message
		SET
			flag = 1,
			hash = ?
		WHERE flag = 0 ORDER BY id LIMIT 1`, hash); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := db.Get(ret, `SELECT id, body, created_at FROM message WHERE hash = ?`, hash); err != nil {
		return nil, err
	}

	if _, err := db.Exec(`DELETE FROM message where id = ?`, ret.ID); err != nil {
		return nil, err
	}

	if err := q.limitter.Free(); err != nil {
		return nil, err
	}
	return ret, nil
}
