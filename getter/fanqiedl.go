package getter

import (
	"github.com/Aiicy/htmlquery"
	"github.com/henson/proxypool/pkg/models"
	"golang.org/x/net/html"
	clog "unknwon.dev/clog/v2"
)

// FQDL get ip from https://www.fanqieip.com/
func FQDL() (result []*models.IP) {
	pollURL := "https://www.fanqieip.com/free/1"
	doc, _ := htmlquery.LoadURL(pollURL)
	trNode, err := htmlquery.Find(doc, "//table[@class='layui-table']//tbody//tr")
	if err != nil {
		clog.Warn(err.Error())
	}
	for i := 0; i < len(trNode); i++ {
		tdNode, _ := htmlquery.Find(trNode[i], "//td")
		ip := extractTextFromDivNode(tdNode[0])
		port := extractTextFromDivNode(tdNode[1])
		Type := "http"
		speed := htmlquery.InnerText(tdNode[4])

		IP := models.NewIP()
		IP.Data = ip + ":" + port
		IP.Type1 = Type
		IP.Source = "fanqieip"
		IP.Speed = extractSpeed(speed)
		result = append(result, IP)
	}

	clog.Info("[fanqiedl] done")
	return
}

func extractTextFromDivNode(node *html.Node) string {
	divNode, _ := htmlquery.Find(node, "//div")
	divOut := htmlquery.InnerText(divNode[0])
	return divOut
}
