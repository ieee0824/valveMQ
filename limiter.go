package valve

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type Log struct {
	LastDequeueTime time.Time `json:"last_dequeue_time" sql:"last_dequeue_time"`
	hash            string
}

func (l *Log) Block() error {
	hash := fmt.Sprintf("%X", sha256.Sum256([]byte(strconv.Itoa(rand.Int()))))
	tx := db.MustBegin()
	if _, err := tx.Exec(`
	UPDATE
		log
	SET
		hash = ?
	WHERE hash = ""`, hash); err != nil {
		tx.Rollback()
		return err
	}
	l.hash = hash

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	if err := db.Get(&l.LastDequeueTime, `SELECT last_dequeue_time FROM log WHERE hash = ?`, hash); err != nil {
		return err
	}

	return nil
}

func (l *Log) Nop() error {
	tx := db.MustBegin()
	if _, err := tx.Exec(`
	UPDATE
		log
	SET
		hash = ?
	WHERE hash = ?`, "", l.hash); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	l.hash = ""
	l.LastDequeueTime = time.Date(0, 0, 0, 0, 0, 0, 0, time.Local)
	return nil

}

func (l *Log) Free() error {
	tx := db.MustBegin()
	if _, err := tx.Exec(`
	UPDATE
		log
	SET
		last_dequeue_time = NOW(6),
		hash = ?
	WHERE hash = ?`, "", l.hash); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	l.hash = ""
	l.LastDequeueTime = time.Date(0, 0, 0, 0, 0, 0, 0, time.Local)
	return nil
}
