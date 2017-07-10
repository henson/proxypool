package main

import (
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/henson/ProxyPool/api"
	"github.com/henson/ProxyPool/getter"
	"github.com/henson/ProxyPool/models"
	"github.com/henson/ProxyPool/storage"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	ipChan := make(chan *models.IP, 2000)
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

func run(ipChan chan<- *models.IP) {
	var wg sync.WaitGroup
	funs := []func() []*models.IP{
		getter.Data5u,
		getter.IP66,
		getter.KDL,
		getter.GBJ,
		getter.Xici,
		getter.XDL,
		//getter.IP181, 	//无内容的采集脚本，用作系统容错实验
		//getter.YDL,		//失效的采集脚本，用作系统容错实验
		getter.PLP,
	}
	for _, f := range funs {
		wg.Add(1)
		go func(f func() []*models.IP) {
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
