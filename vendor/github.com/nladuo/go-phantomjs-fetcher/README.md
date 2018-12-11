# go-phantomjs-fetcher
[pyspider phantomjs fetcher](https://github.com/binux/pyspider/tree/master/pyspider/fetcher) clone in golang.

## Installation
### Install PhantomJS
You can download the phantomjs executable binary [here](http://phantomjs.org/download.html). And add it to your $PATH.
### Clone the Source
``` shell
go get github.com/PuerkitoBio/goquery           # used in example
go get github.com/nladuo/go-phantomjs-fetcher
```

## Example
```shell
cd $GOPATH/src/github.com/nladuo/go-phantomjs-fetcher
go run ./example/mock_baidu_search.go
```
![mock_baidu_search](./example/mock_baidu_search.png)
## LICENSE
MIT