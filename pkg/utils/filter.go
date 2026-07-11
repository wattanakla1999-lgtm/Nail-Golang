package utils

import (
	"fmt"

	"gorm.io/gorm"
)

func ApplyLikeFilters(query *gorm.DB, filters map[string]string) *gorm.DB {
	for column, value := range filters {
		if value != "" {
			query = query.Where(fmt.Sprintf("%s LIKE ?", column), "%"+value+"%")
		}
	}

	return query
}
