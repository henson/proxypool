package getter

import (
	"log"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/parnurzeal/gorequest"
)

// IP181 get ip from ip181.com
func IP181() (result []string) {
	pollURL := "http://www.ip181.com"
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

	doc.Find("body > div:nth-child(3) > div.panel.panel-info > div.panel-body > div > div:nth-child(2) > table > tbody > tr").Each(func(i int, s *goquery.Selection) {
		node := strconv.Itoa(i + 1)
		sf, _ := s.Find("tr:nth-child(" + node + ") > td:nth-child(1)").Html()
		ff, _ := s.Find("tr:nth-child(" + node + ") > td:nth-child(2)").Html()
		result = append(result, (sf + ":" + ff))
	})
	log.Println("IP181 done.")
	return result[1:]
}
