package getter

import (
	"log"

	"github.com/Aiicy/ProxyPool/pkg/models"
	"github.com/PuerkitoBio/goquery"
	"github.com/parnurzeal/gorequest"
)

// IP181 get ip from ip181.com
func IP181() (result []*models.IP) {
	pollURL := "http://www.ip181.com/"
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
	doc.Find("tr.warning").Each(func(i int, s *goquery.Selection) {
		ss := s.Find("td:nth-child(1)").Text()
		sss := s.Find("td:nth-child(2)").Text()
		ssss := s.Find("td:nth-child(4)").Text()
		ip := models.NewIP()
		ip.Data = ss + ":" + sss
		ip.Type = ssss
		result = append(result, ip)
	})

	log.Println("IP181 done.")
	return
}
