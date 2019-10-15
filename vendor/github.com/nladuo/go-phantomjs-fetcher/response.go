package phantomjs

import (
	"net/http"
)

type Response struct {
	OrigUrl        string            `json:"orig_url"`
	Url            string            `json:"url"`
	Headers        map[string]string `json:"headers"`
	StatusCode     int               `json:"status_code"`
	Content        string            `json:"content"`
	Cookies        []http.Cookie     `json:"cookies"`
	Time           float64           `json:"time"`
	JsScriptResult string
}
