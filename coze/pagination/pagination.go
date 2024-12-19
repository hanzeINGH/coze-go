package pagination

import "github.com/chyroc/go-ptr"

type PageRequest struct {
	PageToken string `json:"page_token,omitempty"`
	PageNum   int    `json:"page_num,omitempty"`
	PageSize  int    `json:"page_size,omitempty"`
}

type PageResponse[T any] struct {
	HasMore bool   `json:"has_more"`
	Total   int    `json:"total"`
	Data    []T    `json:"data"`
	LastID  string `json:"last_id,omitempty"`
	NextID  string `json:"next_id,omitempty"`
	LogID   string `json:"log_id,omitempty"`
}

type basePager[T any] struct {
	pageFetcher    PageFetcher[T]
	pageSize       int
	currentPage    *PageResponse[T]
	currentIndex   int
	currentPageNum int
	err            error
	cur            T
}

func (p *basePager[T]) Err() error {
	return p.err
}

func (p *basePager[T]) Items() []T {
	return ptr.Value(p.currentPage).Data
}

func (p *basePager[T]) Current() T {
	return p.cur
}

func (p *basePager[T]) Total() int {
	return ptr.Value(p.currentPage).Total
}

func (p *basePager[T]) HasMore() bool {
	return ptr.Value(p.currentPage).HasMore
}

// PageFetcher 接口
type PageFetcher[T any] func(request *PageRequest) (*PageResponse[T], error)

// PageNumBasedPager 实现
type PageNumBasedPager[T any] struct {
	basePager[T]
}

func NewPageNumBasedPager[T any](fetcher PageFetcher[T], pageSize, pageNum int) (*PageNumBasedPager[T], error) {
	if pageNum <= 0 {
		pageNum = 1
	}
	paginator := &PageNumBasedPager[T]{basePager: basePager[T]{pageFetcher: fetcher, pageSize: pageSize, currentPageNum: pageNum}}
	err := paginator.fetchNextPage()
	return paginator, err
}

func (p *PageNumBasedPager[T]) fetchNextPage() error {
	request := &PageRequest{PageNum: p.currentPageNum, PageSize: p.pageSize}
	var err error
	p.currentPage, err = p.pageFetcher(request)
	if err != nil {
		return err
	}
	p.currentIndex = 0
	p.currentPageNum++
	return nil
}

func (p *PageNumBasedPager[T]) Next() bool {
	if p.currentIndex < len(ptr.Value(p.currentPage).Data) {
		p.cur = p.currentPage.Data[p.currentIndex]
		p.currentIndex++
		return true
	}
	if p.currentPage.HasMore {
		err := p.fetchNextPage()
		if err != nil {
			return false
		}
		if len(p.currentPage.Data) == 0 {
			return false
		}
		p.cur = p.currentPage.Data[p.currentIndex]
		p.currentIndex++
		return true
	}
	return false
}

// TokenBasedPager 实现
type TokenBasedPager[T any] struct {
	basePager[T]
	pageToken string
}

func NewTokenBasedPager[T any](fetcher PageFetcher[T], pageSize int, nextID string) (*TokenBasedPager[T], error) {
	paginator := &TokenBasedPager[T]{basePager: basePager[T]{pageFetcher: fetcher, pageSize: pageSize}}
	err := paginator.fetchNextPage()
	return paginator, err
}

func (p *TokenBasedPager[T]) fetchNextPage() error {
	request := &PageRequest{PageToken: p.pageToken, PageSize: p.pageSize}
	var err error
	p.currentPage, err = p.pageFetcher(request)
	if err != nil {
		p.err = err
		return err
	}
	p.currentIndex = 0
	p.pageToken = p.currentPage.NextID
	return nil
}

func (p *TokenBasedPager[T]) Next() bool {
	if p.currentIndex < len(ptr.Value(p.currentPage).Data) {
		p.cur = p.currentPage.Data[p.currentIndex]
		p.currentIndex++
		return true
	}
	if p.currentPage.HasMore && p.pageToken != "" {
		err := p.fetchNextPage()
		if err != nil {
			return false
		}
		if len(p.currentPage.Data) == 0 {
			return false
		}
		p.cur = p.currentPage.Data[p.currentIndex]
		p.currentIndex++
		return true
	}
	return false
}
