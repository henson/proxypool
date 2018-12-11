package phantomjs

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	ErrPhantomJSNotFound = "\"phantomjs\": executable file not found in $PATH"
	ErrFetcherJSNotFound = "cannot find ./phantomjs_fetcher.js or $GOPATH/github.com/src/nladuo/go-phantomjs-fetcher/phantomjs_fetcher.js"
)

const (
	RUN_AT_DOC_START = "document-start"
	RUN_AT_DOC_END   = "document-end"
)

const (
	type_WINDOWS = "windows os"
	type_UNIX    = "*nix os"
)

type Fetcher struct {
	client         *http.Client
	ProxyPort      string
	AllowRedirects bool
	//phantomJSPid   int
	//phantomJSHandlePtr uintptr
	DefaultOption *Option
	cmd           *exec.Cmd
}

func NewFetcher(port int, option *Option) (*Fetcher, error) {
	var fetcher Fetcher
	fetcher.ProxyPort = strconv.FormatInt(int64(port), 10)
	phantomJSPath, err := fetcher.checkPhantomJS()
	if err != nil {
		return nil, errors.New(ErrPhantomJSNotFound)
	}
	fetcherJSPath, err := fetcher.checkFetcherJS()
	if err != nil {
		return nil, err
	}
	err = fetcher.startPhantomJSServer(phantomJSPath, fetcherJSPath)
	if err != nil {
		return nil, err
	}
	time.Sleep(2 * time.Second)
	if option != nil {
		fetcher.DefaultOption = option
	} else {
		headers := make(map[string]string)
		headers["User-Agent"] = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.110 Safari/537.36"
		fetcher.DefaultOption = &Option{
			Headers:        headers,
			Timeout:        120,
			UseGzip:        true,
			AllowRedirects: true,
		}
	}

	return &fetcher, nil
}

//shutdown the phantomjs server in windows or linux
func (this *Fetcher) ShutDownPhantomJSServer() {
	this.cmd.Process.Kill()
}

func (this *Fetcher) startPhantomJSServer(phantomJSPath, fetcherJSPath string) error {
	args := []string{fetcherJSPath, this.ProxyPort}
	cmd := exec.Command("phantomjs", args...)
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	this.cmd = cmd
	return cmd.Start()
}

//send httpGet request by phantomjs
func (this *Fetcher) Get(url string) (*Response, error) {
	return this.GetWithJS(url, "", RUN_AT_DOC_START)
}

//send httpGet request by phantomjs with the js_script
func (this *Fetcher) GetWithJS(url, js_script, js_run_at string) (*Response, error) {
	return this.GetWithOption(url, js_script, js_run_at, this.DefaultOption)
}

type postData struct {
	LoadImages     bool              `json:"load_images"`
	Url            string            `json:"url"`
	Headers        map[string]string `json:"headers"`
	Timeout        int               `json:"timeout"`
	UseGzip        bool              `json:"use_gzip"`
	AllowRedirects bool              `json:"allow_redirects"`
	Method         string            `json:"method"`
	JsScript       string            `json:"js_script"`
	JsRunAt        string            `json:"js_run_at"`
}

//send httpGet request by phantomjs with the js_script and some option like headers
func (this *Fetcher) GetWithOption(url, js_script, js_run_at string, option *Option) (*Response, error) {
	_postData := postData{
		LoadImages:     false,
		Url:            url,
		Headers:        option.Headers,
		Timeout:        option.Timeout,
		UseGzip:        option.UseGzip,
		AllowRedirects: option.AllowRedirects,
		Method:         "GET",
		JsScript:       js_script,
		JsRunAt:        js_run_at,
	}

	data, err := json.Marshal(&_postData)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(data)
	res, err := http.Post("http://localhost:"+this.ProxyPort, "application/json;charset=utf-8", buffer)
	if err != nil {
		panic(err)
	}

	byte_data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	res.Body.Close()
	var response Response

	err = json.Unmarshal(byte_data, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

//check the existence of executable "phantomjs" in $PATH
func (this *Fetcher) checkPhantomJS() (string, error) {
	phantomJSPath, err := exec.LookPath("phantomjs")
	if err != nil {
		return "", errors.New(ErrPhantomJSNotFound)
	}
	return phantomJSPath, nil
}

// exePath returns the executable path.
func exePath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	return filepath.Abs(file)
}

//check the existence of
//$GOPATH/github.com/nladuo/go-phantomjs-fetcher/phantomjs_fetcher.js
func (this *Fetcher) checkFetcherJS() (string, error) {
	if this.DefaultOption != nil && len(this.DefaultOption.FetcherJsPath) > 0 {
		return this.DefaultOption.FetcherJsPath, nil
	}
	p, err := exePath()
	if err != nil {
		return "", err
	}

	fetcherJSPath := filepath.Join(filepath.Dir(p), "phantomjs_fetcher.js")
	if this.exist(fetcherJSPath) {
		return fetcherJSPath, nil
	}

	str := os.Getenv("GOPATH")
	var paths []string
	if this.getOSType() == type_UNIX {
		paths = strings.Split(str, ":")
	} else {
		paths = strings.Split(str, ";")
	}
	for _, path := range paths {
		fetcherJSPath := path + "/src/github.com/nladuo/go-phantomjs-fetcher/phantomjs_fetcher.js"
		if this.exist(fetcherJSPath) {
			return fetcherJSPath, nil
		}
	}
	return "", errors.New(ErrFetcherJSNotFound)
}

// get os_type, *nix or windows
func (this *Fetcher) getOSType() string {
	//in *nix os, contain ls in $PATH.
	_, err := exec.LookPath("ls")
	if err != nil {
		return type_WINDOWS
	}
	return type_UNIX
}

//check the file existence
func (this *Fetcher) exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
