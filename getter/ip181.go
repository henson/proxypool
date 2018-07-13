package getter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/go-clog/clog"

	"github.com/Aiicy/ProxyPool/pkg/models"
	"github.com/parnurzeal/gorequest"
)

type ip181 struct {
	ErrorCode string   `json:"ERRORCODE"`
	Results   []Result `json:"RESULT"`
}

type Result struct {
	Postion string `json:"position"`
	Port    string `json:"port"`
	Ip      string `json:"ip"`
}

// IP181 get ip from ip181.com
func IP181() (result []*models.IP) {
	var ips ip181
	var results []Result

	pollURL := "http://www.ip181.com/"
	resp, _, errs := gorequest.New().Get(pollURL).End()
	if errs != nil {
		log.Println(errs)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	err := json.Unmarshal(body, &ips)

	if err != nil {
		fmt.Println(err)
	}

	results = ips.Results

	for i := 0; i < len(results); i++ {
		ip := models.NewIP()
		ip.Data = results[i].Ip + ":" + results[i].Port
		ip.Type1 = "http"
		clog.Info("[IP181] ip.Data: %s,ip.Type: %s", ip.Data, ip.Type1)
		result = append(result, ip)
	}

	clog.Info("IP181 done.")
	return
}
