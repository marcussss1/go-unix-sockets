package main

import (
	"fmt"
	"log"
	"math/rand"
	"my-server/client"
	"my-server/server"
	"sync"
)

func main() {
	srv, err := server.New(8321, [4]byte{127, 0, 0, 1}, 5)
	if err != nil {
		log.Fatal(fmt.Errorf("create server: %w", err))
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go startServer(wg, srv)
	wg.Wait()

	wg.Add(2)
	go cliSend(wg, []string{"1"}, 5)
	go cliSend(wg, []string{"2"}, 5)
	wg.Wait()
}

func startServer(wg *sync.WaitGroup, srv *server.Server) {
	wg.Done()
	err := srv.Start()
	fmt.Println(err)
}

func cliSend(wg *sync.WaitGroup, msgs []string, countRequests int) {
	defer wg.Done()

	mx := &sync.Mutex{}
	wg2 := &sync.WaitGroup{}
	wg2.Add(countRequests)

	for i := 0; i < countRequests; i++ {
		go func() {
			defer wg2.Done()
			defer mx.Unlock()

			mx.Lock()
			cli, err := client.New(8321, [4]byte{127, 0, 0, 1})
			if err != nil {
				fmt.Printf("create client: %v\n", err)
				return
			}

			resp, err := cli.SendMsg(msgs[rand.Intn(len(msgs))])
			if err != nil {
				fmt.Printf("send message: %v\n", err)
				return
			}

			fmt.Println(resp)
		}()
	}
	wg2.Wait()
}
