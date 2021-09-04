package query

type Pagination struct {
	PageNum int `form:"page" json:"pageNum"`
	PageSize int `form:"size" json:"pageSize"`
}

func (p *Pagination) Limit() int {
	return p.PageSize
}

func (p *Pagination) Offset() int {
	return (p.PageNum - 1) * p.PageSize
}