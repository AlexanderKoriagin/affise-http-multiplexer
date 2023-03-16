package service

type Checker interface {
	UrlsQty(urls []string) bool
}
