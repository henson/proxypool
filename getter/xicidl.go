package getter

import (
	"log"
	"regexp"
	"strings"

	"github.com/henson/proxypool/pkg/models"
	"github.com/nladuo/go-phantomjs-fetcher"
)

// Xici get ip from xicidaili.com
func Xici() (result []*models.IP) {
	pollURL := "http://www.xicidaili.com/nn/"

	fetcher, err := phantomjs.NewFetcher(2017, nil)
	defer fetcher.ShutDownPhantomJSServer()
	if err != nil {
		log.Println(err.Error())
		return
	}
	jsScript := "function() {s=document.documentElement.outerHTML;document.write('<body></body>');document.body.innerText=s;}"
	jsRunAt := phantomjs.RUN_AT_DOC_END
	resp, err := fetcher.GetWithJS(pollURL, jsScript, jsRunAt)
	if err != nil {
		log.Println(err.Error())
		return
	}
	re, _ := regexp.Compile("<td>(\\d+\\.){3}\\d+</td>.+?(\\d{2,4})</td>")
	temp := re.FindAllString(strings.Replace(strings.Replace(resp.Content, "&lt;", "<", -1), "&gt;", ">", -1), -1)

	for _, v := range temp {
		v = strings.Replace(v, "<td>", "", -1)
		v = strings.Replace(v, "</td>", "", -1)
		v = strings.Replace(v, " ", "", -1)
		v = strings.Replace(v, "<br>", ":", -1)
		ip := models.NewIP()
		ip.Data = v
		ip.Type1 = "http"
		ip.Source = "xicidl"
		result = append(result, ip)
	}
	log.Println("Xici done.")
	return
}
