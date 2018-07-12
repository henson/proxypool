package getter

import (
	"fmt"
	"log"

	"github.com/Aiicy/ProxyPool/pkg/models"
	"github.com/Aiicy/htmlquery"
)

//PLP get ip from proxylistplus.com
func PLP() (result []*models.IP) {
	pollURL := "https://list.proxylistplus.com/Fresh-HTTP-Proxy-List-1"
	doc, _ := htmlquery.LoadURL(pollURL)
	trNode, err := htmlquery.Find(doc, "//div[@class='hfeed site']//table[@class='bg']//tbody//tr")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("tdNode = %v\n", trNode)
	for i := 3; i < len(trNode); i++ {
		tdNode, _ := htmlquery.Find(trNode[i], "//td")
		ip := htmlquery.InnerText(tdNode[1])
		port := htmlquery.InnerText(tdNode[2])
		Type := htmlquery.InnerText(tdNode[6])

		Ip := models.NewIP()
		Ip.Data = ip + ":" + port

		if Type == "yes" {
			Ip.Type1 = "http"
			Ip.Type2 = "https"

		} else if Type == "no" {
			Ip.Type1 = "http"
		}

		fmt.Printf("[PLP] ip.Data = %s,ip.Type = %s,%s\n", Ip.Data, Ip.Type1, Ip.Type2)

		result = append(result, Ip)
	}

	log.Println("PLP done.")
	return
}
