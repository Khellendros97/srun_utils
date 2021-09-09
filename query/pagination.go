package query

type Pagination struct {
	PageNum int `form:"pageNum" json:"pageNum"`
	PageSize int `form:"pageSize" json:"pageSize"`
}

func (p *Pagination) Limit() int {
	return p.PageSize
}

func (p *Pagination) Offset() int {
	return (p.PageNum - 1) * p.PageSize
}