package query

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type FuzzyQuery struct {
	*gorm.DB
}

func (q *FuzzyQuery)Precise(enable bool, key string, value interface{}) {
	if !enable {
		return
	}
	q.DB = q.DB.Where(fmt.Sprintf("%s = ?", key), value)
}

func (q *FuzzyQuery)Fuzzy(enable bool, key string, value interface{}) {
	if !enable {
		return
	}
	q.DB = q.DB.Where(fmt.Sprintf("%s like ?", key), fmt.Sprintf("%v%%", value))
}

func (q *FuzzyQuery)Get() *gorm.DB {
	return q.DB
}

func (q *FuzzyQuery)GetWhere(like bool) func (bool, string, interface{}) {
	if like {
		return q.Fuzzy
	} else {
		return q.Precise
	}
}