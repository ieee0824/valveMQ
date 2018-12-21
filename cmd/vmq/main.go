package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	_ "github.com/go-sql-driver/mysql"
	valve "github.com/ieee0824/valveMQ"
)

func sha256Sum(data string) string {
	bytes := sha256.Sum256([]byte(data))
	return hex.EncodeToString(bytes[:])
}

func main() {
	log.SetFlags(log.Llongfile)
	r := gin.Default()
	q := valve.NewQueue()
	cfg := valve.NewConfig()
	q.SetLimit(cfg.DequeueLimit)

	if err := valve.DBInit(cfg); err != nil {
		log.Fatalln(err)
	}

	r.POST("/enqueue", func(ctx *gin.Context) {
		msg := &valve.Message{}
		d := json.NewDecoder(ctx.Request.Body)
		defer ctx.Request.Body.Close()
		if err := d.Decode(msg); err != nil {
			log.Println(err)
			ctx.JSON(http.StatusBadRequest, nil)
			return
		}

		msg.RequestID = sha256Sum(msg.Body + fmt.Sprintf("%d", time.Now().UnixNano()))

		if err := q.Enqueue(msg); err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, nil)
			return
		}
		ctx.JSON(http.StatusOK, msg.RequestID)
	})

	r.GET("/dequeue", func(ctx *gin.Context) {
		cancel := time.NewTicker(30 * time.Second)
		interval := time.NewTicker(100 * time.Millisecond)
		defer cancel.Stop()
		defer interval.Stop()
		for {
			select {
			case <-cancel.C:
				ctx.JSON(http.StatusNotFound, nil)
				return
			case <-interval.C:
				msg, err := q.Dequeue()
				if err != nil {
					if err.Error() != "limit" &&
						err.Error() != "sql: no rows in result set" {
						log.Println(err)
						ctx.JSON(http.StatusNotFound, nil)
						return
					}
				} else {
					ctx.JSON(http.StatusOK, msg)
					return
				}
			}
		}

	})

	if err := r.Run(fmt.Sprintf(":%s", cfg.APIPort)); err != nil {
		log.Fatalln(err)
	}
}
