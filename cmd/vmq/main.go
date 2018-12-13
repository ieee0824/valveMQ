package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/go-sql-driver/mysql"
	valve "github.com/ieee0824/valveMQ"
)

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

		if err := q.Enqueue(msg); err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, nil)
			return
		}
		ctx.JSON(http.StatusOK, true)
	})

	r.GET("/dequeue", func(ctx *gin.Context) {
		msg, err := q.Dequeue()
		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusNotFound, nil)
			return
		}
		ctx.JSON(http.StatusOK, msg)
	})

	if err := r.Run(fmt.Sprintf(":%s", cfg.APIPort)); err != nil {
		log.Fatalln(err)
	}
}
