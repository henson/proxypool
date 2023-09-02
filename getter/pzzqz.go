package getter

import (
	"github.com/antchfx/htmlquery"
	"github.com/henson/proxypool/pkg/models"
	clog "unknwon.dev/clog/v2"
)

// PZZQZ get ip from http://pzzqz.com/
func PZZQZ() (result []*models.IP) {
	pollURL := "http://pzzqz.com/"
	doc, err := htmlquery.LoadURL(pollURL)
	if err != nil {
		clog.Warn(err.Error())
		return
	}
	trNode, err := htmlquery.QueryAll(doc, "//table[@class='table table-hover']//tbody//tr")
	if err != nil {
		clog.Warn(err.Error())
	}
	for i := 0; i < len(trNode); i++ {
		tdNode, err_ := htmlquery.QueryAll(trNode[i], "//td")
		if err_ != nil {
			clog.Warn(err_.Error())
			return
		}
		ip := htmlquery.InnerText(tdNode[0])
		port := htmlquery.InnerText(tdNode[1])
		Type := htmlquery.InnerText(tdNode[4])

		IP := models.NewIP()
		IP.Data = ip + ":" + port
		IP.Type1 = Type
		IP.Source = "pzzqz"
		result = append(result, IP)
	}

	clog.Info("[pzzqz] done")
	return
}
