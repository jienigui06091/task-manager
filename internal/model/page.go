package model


type Page struct {
	List     any   `json:"list"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
	Total    int64 `json:"total"`
}

// func (p *Page) GetPageSize() int {
// 	return p.PageSize
// }

// func (p *Page) GetPage() int {
// 	return p.Page
// }

// func (p *Page) SetTotal(total int64) {
// 	p.Total = total
// }

// func (p *Page) SetList(list any) {
// 	p.List = list
// }
