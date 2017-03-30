package main

import (
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/henson/ProxyPool/api"
	"github.com/henson/ProxyPool/getter"
	"github.com/henson/ProxyPool/storage"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	ipChan := make(chan string, 1000)
	conn := storage.NewStorage()

	// Start HTTP
	go func() {
		api.Run()
	}()

	// Check the IPs in DB
	go func() {
		storage.CheckProxyDB()
	}()

	// Check the IPs in channel
	for i := 0; i < 50; i++ {
		go func() {
			for {
				storage.CheckProxy(<-ipChan)
			}
		}()
	}

	// Start getters to scraper IP and put it in channel
	for {
		x := conn.Count()
		log.Printf("Chan: %v, IP: %v\n", len(ipChan), x)
		if len(ipChan) < 100 {
			go run(ipChan)
		}
		time.Sleep(10 * time.Minute)
	}
}

func run(ipChan chan<- string) {
	var wg sync.WaitGroup
	funs := []func() []string{
		getter.IP66,
		getter.KDL,
		getter.GBJ,
		getter.Xici,
		getter.IP181,
		getter.YDL,
	}
	for _, f := range funs {
		wg.Add(1)
		go func(f func() []string) {
			temp := f()
			for _, v := range temp {
				ipChan <- v
			}
			wg.Done()
		}(f)
	}
	wg.Wait()
	log.Println("All getters finished.")
}
