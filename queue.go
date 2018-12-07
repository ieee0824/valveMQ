package valve

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Queue struct {
	setting Setting
}

func (q *Queue) Enqueue(m *Message) error {
	now := time.Now()

	m.CreatedAt = &now

	tx := db.MustBegin()

	if _, err := tx.Exec(`INSERT INTO message (body, expire) VALUES (?, ?)`, m.Body, m.Expire); err != nil {
		return err
	}
	return tx.Commit()
}

func (q *Queue) Dequeue() (*Message, error) {
	tx := db.MustBegin()

	hash := fmt.Sprintf("%X", sha256.Sum256([]byte(strconv.Itoa(rand.Int()))))

	ret := &Message{}
	//SELECT * FROM message ORDER BY id LIMIT 1
	if _, err := tx.Exec(`UPDATE
			message
		SET
			flag = 1,
			hash = ?
		WHERE flag = 0 ORDER BY id LIMIT 1`, hash); err != nil {
		return nil, err
	}

	if err := tx.Get(ret, `SELECT id, body, created_at FROM message WHERE hash = ?`, hash); err != nil {
		return nil, err
	}

	if _, err := tx.Exec(`DELETE FROM message where id = ?`, ret.ID); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return ret, nil
}
