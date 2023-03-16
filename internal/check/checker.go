package check

type Checker struct {
	Features Features
}

func Init(maxUrlQuantity int) *Checker {
	return &Checker{Features: &service{urlQuantity: maxUrlQuantity}}
}

func (c *Checker) UrlsQty(urls []string) bool {
	return c.Features.CheckUrl(urls)
}
