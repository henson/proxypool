package getter

import (
	"log"
	"strings"

	"github.com/henson/ProxyPool/models"
	"github.com/parnurzeal/gorequest"
)

// IP66 get ip from 66ip.cn
func IP66() (result []*models.IP) {
	pollURL := "http://www.66ip.cn/mo.php?tqsl=100"
	_, body, errs := gorequest.New().Get(pollURL).End()
	if errs != nil {
		log.Println(errs)
		return
	}
	body = strings.Split(body, "var sogou_ad_width=960;")[1]
	body = strings.Split(body, "</script>")[1]
	body = strings.Split(body, "</div>")[0]
	body = strings.TrimSpace(body)
	body = strings.Replace(body, "	", "", -1)
	temp := strings.Split(body, "<br />")
	for index := 0; index < len(temp[:len(temp)-1]); index++ {
		ip := models.NewIP()
		ip.Data = strings.TrimSpace(temp[index])
		ip.Type = "http"
		result = append(result, ip)
	}
	log.Println("IP66 done.")
	return
}
