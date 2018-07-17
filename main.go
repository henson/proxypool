package main

import (
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/henson/proxypool/api"
	"github.com/henson/proxypool/getter"
	"github.com/henson/proxypool/pkg/initial"
	"github.com/henson/proxypool/pkg/models"
	"github.com/henson/proxypool/pkg/storage"
)

func main() {

	//init the database
	initial.GlobalInit()

	runtime.GOMAXPROCS(runtime.NumCPU())
	ipChan := make(chan *models.IP, 2000)

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
		n := models.CountIPs()
		log.Printf("Chan: %v, IP: %v\n", len(ipChan), n)
		if len(ipChan) < 100 {
			go run(ipChan)
		}
		time.Sleep(10 * time.Minute)
	}
}

func run(ipChan chan<- *models.IP) {
	var wg sync.WaitGroup
	funs := []func() []*models.IP{
		//getter.Data5u,
		getter.IP66,
		//getter.KDL,
		//getter.GBJ,	//因为网站限制，无法正常下载数据
		//getter.Xici,
		//getter.XDL,
		getter.IP181,
		//getter.YDL,	//失效的采集脚本，用作系统容错实验
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
