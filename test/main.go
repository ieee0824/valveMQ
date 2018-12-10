package main

import (
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	valve "github.com/ieee0824/valveMQ"
)

const TEST_NUM = 10

func main() {
	log.SetFlags(log.Llongfile)
	q := &valve.Queue{}
	q.SetLimit(50)
	enqStart := time.Now()
	for i := 0; i < TEST_NUM; i++ {
		q.Enqueue(&valve.Message{
			Body: "hoge",
		})
	}
	enqEnd := time.Now()

	eps := float64(TEST_NUM) / (float64(enqEnd.Sub(enqStart)) / float64(time.Second))
	fmt.Println("enqueue speed: ", int(eps))

	ac := make(chan bool)
	a := make([]int, 0, TEST_NUM/2)

	bc := make(chan bool)
	b := make([]int, 0, TEST_NUM/2)

	deqStart := time.Now()
	go func() {
		for i := 0; i < TEST_NUM/2; i++ {
			m, err := q.Dequeue()
			if err != nil {
				log.Println(err)
				continue
			}
			a = append(a, m.ID)
		}
		ac <- true
	}()

	go func() {
		for i := 0; i < TEST_NUM/2; i++ {
			m, err := q.Dequeue()
			if err != nil {
				log.Println(err)
				continue
			}
			b = append(b, m.ID)
		}
		bc <- true
	}()

	<-ac
	<-bc

	deqEnd := time.Now()

	dps := float64(TEST_NUM) / (float64(deqEnd.Sub(deqStart)) / float64(time.Second))
	fmt.Println("dequeue speed: ", int(dps))

	cnt := 0
	for _, va := range a {
		for _, vb := range b {
			if va == vb {
				cnt++
			}
		}
	}

	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(cnt)

}
