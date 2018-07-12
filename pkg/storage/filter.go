package storage

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/Aiicy/ProxyPool/pkg/models"
	"github.com/parnurzeal/gorequest"
)

// CheckProxy .
func CheckProxy(ip *models.IP) {
	if CheckIP(ip) {
		ProxyAdd(ip)
	}
}

// CheckIP is to check the ip work or not
func CheckIP(ip *models.IP) bool {
	pollURL := "http://httpbin.org/get"
	var testIP string
	if ip.Type2 == "https" {
		testIP = "https://" + ip.Data
	} else {
		testIP = "http://" + ip.Data
	}
	//fmt.Println(testIP)
	resp, _, errs := gorequest.New().Proxy(testIP).Get(pollURL).End()
	if errs != nil {
		return false
	}
	if resp.StatusCode == 200 {
		return true
	}
	return false
}

// CheckProxyDB to check the ip in DB
func CheckProxyDB() {

	x := models.CountIPs()
	log.Println("Before check, DB has:", x, "records.")
	ips, err := models.GetAll()
	if err != nil {
		log.Println(err.Error())
		return
	}
	var wg sync.WaitGroup
	for _, v := range ips {
		wg.Add(1)
		go func(v *models.IP) {
			if !CheckIP(v) {
				ProxyDel(v)
			}
			wg.Done()
		}(v)
	}
	wg.Wait()
	x = models.CountIPs()
	log.Println("After check, DB has:", x, "records.")
}

// ProxyRandom .
func ProxyRandom() (ip *models.IP) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	ips, _ := models.GetAll()
	x := len(ips)

	return ips[r.Intn(x)]
}

// ProxyFind .
func ProxyFind(value string) (ip *models.IP) {
	ips, err := models.FindAll(value)
	x := len(ips)
	if err != nil {
		log.Println(err)
		return models.NewIP()
	}
	randomNum := RandInt(0, x)
	fmt.Printf("[proxyFind] random num = %d\n", randomNum)
	if randomNum == 0 {
		return models.NewIP()
	}
	return ips[randomNum]
}

// ProxyAdd .
func ProxyAdd(ip *models.IP) {
	models.InsertIps(ip)
}

// ProxyDel .
func ProxyDel(ip *models.IP) {
	models.DeleteIP(ip)
}
