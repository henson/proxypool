package getter

import (
	"github.com/Aiicy/htmlquery"
	"github.com/henson/proxypool/pkg/models"
	clog "unknwon.dev/clog/v2"
)

// KDL get ip from kuaidaili.com
func KDL() (result []*models.IP) {
	pollURL := "http://www.kuaidaili.com/free/inha/"
	doc, err := htmlquery.LoadURL(pollURL)
	if err != nil {
		clog.Warn(err.Error())
		return
	}
	trNode, err := htmlquery.Find(doc, "//table[@class='table table-bordered table-striped']//tbody//tr")
	if err != nil {
		clog.Warn(err.Error())
	}
	for i := 0; i < len(trNode); i++ {
		tdNode, err_ := htmlquery.Find(trNode[i], "//td")
		if err_ != nil {
			clog.Warn(err_.Error())
			continue
		}
		ip := htmlquery.InnerText(tdNode[0])
		port := htmlquery.InnerText(tdNode[1])
		Type := htmlquery.InnerText(tdNode[3])
		speed := htmlquery.InnerText(tdNode[5])

		IP := models.NewIP()
		IP.Data = ip + ":" + port
		if Type == "HTTPS" {
			IP.Type1 = "https"
			IP.Type2 = "https"
		} else if Type == "HTTP" {
			IP.Type1 = "http"
		}
		IP.Source = "KDL"
		IP.Speed = extractSpeed(speed)
		result = append(result, IP)
	}

	clog.Info("[kuaidaili] done")
	return
}
