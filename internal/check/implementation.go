package check

type Features interface {
	CheckUrl(urls []string) bool
}

type service struct {
	urlQuantity int // maximum number of url
}

func (s *service) CheckUrl(urls []string) bool {
	return len(urls) > 0 && len(urls) <= s.urlQuantity
}
