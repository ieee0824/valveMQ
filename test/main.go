package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	valve "github.com/ieee0824/valveMQ"
)

const TEST_NUM = 2048

func main() {
	log.SetFlags(log.Llongfile)
	q := &valve.Queue{}
	for i := 0; i < TEST_NUM; i++ {
		q.Enqueue(&valve.Message{
			Body: "hoge",
		})
	}

	ac := make(chan bool)
	a := []int{}

	bc := make(chan bool)
	b := []int{}

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

	cnt := 0
	for _, va := range a {
		for _, vb := range b {
			if va == vb {
				cnt++
			}
		}
	}

	fmt.Println(b)
	fmt.Println(cnt)

}
