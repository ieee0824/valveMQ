package valve

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type Log struct {
	m               sync.Mutex
	LastDequeueTime time.Time `json:"last_dequeue_time" sql:"last_dequeue_time"`
	hash            string
}

func (l *Log) GetHash() string {
	return l.hash
}

func (l *Log) Block(lim limit) (bool, error) {
	tx := db.MustBegin()
	l.m.Lock()
	defer tx.Commit()
	defer l.m.Unlock()

	if l.hash != "" {
		return false, nil
	}

	hash := fmt.Sprintf("%X", sha256.Sum256([]byte(strconv.Itoa(rand.Int()))))

	//fmt.Println(fmt.Sprintf("00:00:%02f", float64(lim.DqSpan()/time.Millisecond)/1000))
	//fmt.Println(int64(lim.DqSpan() / time.Millisecond))
	if _, err := tx.Exec(`
	UPDATE
		log
	SET
		hash = ?
	WHERE 
		id = 1
	AND
		? < ROUND(UNIX_TIMESTAMP(NOW(4)) * 1000) - ROUND(UNIX_TIMESTAMP(last_dequeue_time) * 1000)
	AND
		hash = ""`, hash, int64(lim.DqSpan()/time.Millisecond)); err != nil {
		tx.Rollback()
		return false, err
	}

	var n int
	if err := tx.Get(&n, "SELECT COUNT(id) FROM log WHERE id = 1 AND hash = ?", hash); err != nil {
		tx.Rollback()
		return false, err
	}

	l.hash = hash
	return n == 1, nil
}

func (l *Log) Nop() error {
	tx := db.MustBegin()
	l.m.Lock()
	defer func() {
		tx.Commit()
		l.m.Unlock()
	}()

	if l.hash == "" {
		return nil
	}

	if _, err := tx.Exec(`
	UPDATE
		log
	SET
		hash = ?
	WHERE hash = ?`, "", l.hash); err != nil {
		tx.Rollback()
		return err
	}
	l.hash = ""
	return nil
}

func (l *Log) Free() error {
	tx := db.MustBegin()
	l.m.Lock()

	defer func() {
		tx.Commit()
		l.m.Unlock()
	}()

	if l.hash == "" {
		return nil
	}

	if _, err := tx.Exec(`
	UPDATE
		log
	SET
		last_dequeue_time = NOW(4),
		hash = ?
	WHERE hash = ?`, "", l.hash); err != nil {
		tx.Rollback()
		return err
	}
	l.hash = ""
	l.LastDequeueTime = time.Date(0, 0, 0, 0, 0, 0, 0, time.Local)
	return nil
}
