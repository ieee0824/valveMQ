package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ieee0824/valveMQ"
)

func main() {
	r := gin.Default()
	q := valve.NewQueue()

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

	if err := r.Run(); err != nil {
		log.Fatalln(err)
	}
}
