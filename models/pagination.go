package models

type (
	PaginationDirection uint
	URLGenerator        func(int, PaginationDirection) string
)

const (
	PaginationNext PaginationDirection = iota
	PaginationPrev
)

type ArticlesPaginator struct {
	Items        ArticleList
	CurrentPage  int
	PageCount    int
	URLGenerator URLGenerator
}

func (p *ArticlesPaginator) HasPrevious() bool {
	return p.CurrentPage > 1
}

func (p *ArticlesPaginator) HasNext() bool {
	return p.CurrentPage < p.PageCount
}

func (p *ArticlesPaginator) PreviousPageURL() string {
	return p.URLGenerator(p.CurrentPage, PaginationPrev)
}

func (p *ArticlesPaginator) NextPageURL() string {
	return p.URLGenerator(p.CurrentPage, PaginationNext)
}
