package paginator

import (
	"strconv"

	"github.com/jinzhu/gorm"
)

type Paginator struct {
	DB      *gorm.DB
	OrderBy []string
	Page    string
	PerPage string
}

type Data struct {
	TotalRecords int         `json:"total_records"`
	Records      interface{} `json:"records"`
	CurrentPage  string      `json:"current_page"`
	TotalPages   int64       `json:"total_pages"`
}

func (p *Paginator) Paginate(dataSource interface{}) *Data {
	db := p.DB

	if len(p.OrderBy) > 0 {
		for _, o := range p.OrderBy {
			db = db.Order(o)
		}
	}

	done := make(chan bool, 1)
	var output Data
	var count int
	var offset int64

	go countRecords(db, dataSource, done, &count)

	if p.Page == "1" {
		offset = 0
	} else {
		tmpPage, _ := strconv.ParseInt(p.Page, 10, 32)
		tmpPerPage, _ := strconv.ParseInt(p.PerPage, 10, 32)
		offset = (tmpPage - 1) * tmpPerPage
	}

	db.Limit(p.PerPage).Offset(offset).Find(dataSource)
	<-done

	output.TotalRecords = count
	output.Records = dataSource
	output.CurrentPage = p.Page
	output.TotalPages = getTotalPages(p.PerPage, count)

	return &output
}

func countRecords(db *gorm.DB, countDataSource interface{}, done chan bool, count *int) {
	db.Model(countDataSource).Count(count)
	done <- true
}

func getTotalPages(perPage string, totalRecords int) int64 {
	perPageInt, _ := strconv.ParseInt(perPage, 10, 32)
	totalPages := float64(totalRecords) / float64(perPageInt)
	// This stupid conversion is needed as golang does not have any round method.
	// Chance for creating a new library ??
	return int64(float64(totalPages) + float64(1.0))
}
