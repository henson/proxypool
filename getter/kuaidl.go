package getter

import (
	"log"
	"strconv"
	"strings"

	"github.com/henson/proxypool/pkg/models"
	"github.com/PuerkitoBio/goquery"
	"github.com/nladuo/go-phantomjs-fetcher"
)

// KDL get ip from kuaidaili.com
func KDL() (result []*models.IP) {
	pollURL := "http://www.kuaidaili.com/proxylist/"
	//create a fetcher which seems to a httpClient
	fetcher, err := phantomjs.NewFetcher(2016, nil)
	defer fetcher.ShutDownPhantomJSServer()
	if err != nil {
		log.Println(err.Error())
		return
	}
	//inject the javascript you want to run in the webpage just like in chrome console.
	jsScript := "function() {s=document.documentElement.outerHTML;document.write('<body></body>');document.body.innerText=s;}"
	//run the injected js_script at the end of loading html
	jsRunAt := phantomjs.RUN_AT_DOC_END
	//send httpGet request with injected js

	for i := 1; i <= 10; i++ {
		resp, err := fetcher.GetWithJS(pollURL+strconv.Itoa(i), jsScript, jsRunAt)
		if err != nil {
			log.Println(err.Error())
			return
		}

		//select search results by goquery
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.Content))
		if err != nil {
			log.Println(err.Error())
			return
		}
		doc.Find("#index_free_list > table > tbody > tr").Each(func(i int, s *goquery.Selection) {
			node := strconv.Itoa(i + 1)
			sf, _ := s.Find("tr:nth-child(" + node + ") > td:nth-child(1)").Html()
			ff, _ := s.Find("tr:nth-child(" + node + ") > td:nth-child(2)").Html()
			hh, _ := s.Find("tr:nth-child(" + node + ") > td:nth-child(4)").Html()
			ip := models.NewIP()
			ip.Data = sf + ":" + ff
			ip.Type1 = strings.ToLower(hh)
			result = append(result, ip)
		})
	}
	log.Println("KDL done.")
	return
}
