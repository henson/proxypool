package getter

import (
	"log"
	"strings"

	"github.com/parnurzeal/gorequest"
)

// IP66 get ip from 66ip.cn
func IP66() (result []string) {
	pollURL := "http://www.66ip.cn/mo.php?tqsl=100&submit=%CC%E1++%C8%A1"
	_, body, errs := gorequest.New().Get(pollURL).End()
	if errs != nil {
		log.Println(errs)
		return
	}

	body = strings.Split(body, "c.js'></script>")[1]
	body = strings.Split(body, "</div>")[0]
	body = strings.TrimSpace(body)
	body = strings.Replace(body, "	", "", -1)
	temp := strings.Split(body, "<br />")
	for index := 0; index < len(temp[:len(temp)-1]); index++ {
		result = append(result, strings.TrimSpace(temp[index]))
	}
	log.Println("IP66 done.")
	return
}
