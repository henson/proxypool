package getter

import (
	clog "unknwon.dev/clog/v2"

	"github.com/Aiicy/htmlquery"
	"github.com/henson/proxypool/pkg/models"
)

//PLP get ip from proxylistplus.com
func PLPSSL() (result []*models.IP) {
	pollURL := "https://list.proxylistplus.com/SSL-List-1"
	doc, _ := htmlquery.LoadURL(pollURL)
	trNode, err := htmlquery.Find(doc, "//div[@class='hfeed site']//table[@class='bg']//tbody//tr")
	if err != nil {
		clog.Warn(err.Error())
	}
	for i := 3; i < len(trNode); i++ {
		tdNode, _ := htmlquery.Find(trNode[i], "//td")
		ip := htmlquery.InnerText(tdNode[1])
		port := htmlquery.InnerText(tdNode[2])
		Type := htmlquery.InnerText(tdNode[6])

		IP := models.NewIP()
		IP.Data = ip + ":" + port

		if Type == "yes" {
			IP.Type1 = "https"
			IP.Type2 = ""

		} else if Type == "no" {
			IP.Type1 = "http"
		}

		clog.Info("[PLP SSL] ip.Data = %s,ip.Type = %s,%s", IP.Data, IP.Type1, IP.Type2)

		result = append(result, IP)
	}

	clog.Info("PLP SSL done.")
	return
}
