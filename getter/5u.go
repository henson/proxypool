package getter

import (
	"log"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/henson/ProxyPool/models"
	"github.com/parnurzeal/gorequest"
)

// Data5u get ip from data5u.com
func Data5u() (result []*models.IP) {
	pollURL := "http://www.data5u.com/free/index.shtml"
	resp, _, errs := gorequest.New().Get(pollURL).End()
	if errs != nil {
		log.Println(errs)
		return
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Println(err.Error())
		return
	}
	doc.Find("body > div.wlist > li:nth-child(2) > ul").Each(func(i int, s *goquery.Selection) {
		node := strconv.Itoa(i + 1)
		ss := s.Find("ul:nth-child(" + node + ") > span:nth-child(1) > li").Text()
		sss := s.Find("ul:nth-child(" + node + ") > span:nth-child(2) > li").Text()
		ssss := s.Find("ul:nth-child(" + node + ") > span:nth-child(4) > li").Text()
		ip := models.NewIP()
		ip.Data = ss + ":" + sss
		ip.Type = ssss
		result = append(result, ip)
	})
	log.Println("Data5u done.")
	return
}
