package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ApplyLikeFilters(c *gin.Context, query *gorm.DB, filters map[string]string) *gorm.DB {
	for queryParam, column := range filters {
		value := c.Query(queryParam)
		if value != "" {
			query = query.Where(fmt.Sprintf("%s LIKE ?", column), "%"+value+"%")
		}
	}

	return query
}
