package pagination

import (
	"math"

	"gorm.io/gorm"
)

type Pagination struct {
	Limit      int         `json:"limit,omitempty"`
	Page       int         `json:"page,omitempty"`
	Sort       string      `json:"sort,omitempty"`
	Search     string      `json:"search,omitempty"`
	Column     *string     `json:"column,omitempty"`
	TotalRows  int64       `json:"total_rows"`
	TotalPages int         `json:"total_pages"`
	Rows       interface{} `json:"rows"`
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination) GetSort() string {
	switch p.Sort {
	case "desc":
		p.Sort = "Id desc"
	case "asc":
		p.Sort = "Id asc"
	default:
		p.Sort = "Id desc"
	}
	return p.Sort
}

func (p *Pagination) GetSearch() string {
	return p.Search
}

func (p *Pagination) GetColumn() string {
	if p.Column == nil {
		return "name"
	}
	return *p.Column
}

func Paginate(value interface{}, pagination *Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db = db.Model(value)
	if search := pagination.GetSearch(); search != "" {

		db = db.Where(pagination.GetColumn()+" LIKE ?", "%"+pagination.GetSearch()+"%")
	}
	db.Count(&totalRows)

	pagination.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.GetLimit())))
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		if search := pagination.GetSearch(); search != "" {
			return db.Where(pagination.GetColumn()+" LIKE ?", "%"+pagination.GetSearch()+"%").Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort()).Omit("password", "verification_code")
		}
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort()).Omit("password", "verification_code")
	}
}
