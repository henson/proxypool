package storage

import (
	"sync"

	"github.com/Henson/ProxyPool/pkg/models"
	"github.com/go-clog/clog"
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
	var pollURL string
	var testIP string
	if ip.Type2 == "https" {
		testIP = "https://" + ip.Data
		pollURL = "https://httpbin.org/get"
	} else {
		testIP = "http://" + ip.Data
		pollURL = "http://httpbin.org/get"
	}
	//fmt.Println(testIP)
	resp, _, errs := gorequest.New().Proxy(testIP).Get(pollURL).End()
	if errs != nil {
		clog.Warn("[CheckIP] testIP = %s, pollURL = %s: Error = %v", testIP, pollURL, errs)
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
	clog.Info("Before check, DB has: %d records.", x)
	ips, err := models.GetAll()
	if err != nil {
		clog.Warn(err.Error())
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
	clog.Info("After check, DB has: %d records.", x)
}

// ProxyRandom .
func ProxyRandom() (ip *models.IP) {

	ips, err := models.GetAll()
	x := len(ips)
	if err != nil {
		clog.Warn(err.Error())
		return models.NewIP()
	}
	randomNum := RandInt(0, x)

	return ips[randomNum]
}

// ProxyFind .
func ProxyFind(value string) (ip *models.IP) {
	ips, err := models.FindAll(value)
	x := len(ips)
	if err != nil {
		clog.Warn(err.Error())
		return models.NewIP()
	}
	randomNum := RandInt(0, x)
	clog.Info("[proxyFind] random num = %d", randomNum)
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
