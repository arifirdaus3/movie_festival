package model

type Pagination struct {
	Limit int `json:"limit" query:"limit"`
	Page  int `json:"page" query:"page"`
}

func (p *Pagination) Default() {
	if p.Limit > 1000 || p.Limit <= 0 {
		p.Limit = 1000
	}
	if p.Page <= 0 {
		p.Page = 1
	}
}
