package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Pagination struct {
	Page   int
	Limit  int
	Offset int
}

func GetPagination(c *gin.Context) Pagination {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	return Pagination{
		Page:   page,
		Limit:  limit,
		Offset: (page - 1) * limit,
	}
}

func Paginate(c *gin.Context, query *gorm.DB, dest any) (Pagination, int64, error) {
	pagination := GetPagination(c)
	var total int64

	if err := query.Count(&total).Error; err != nil {
		return pagination, total, err
	}

	err := query.
		Order("id ASC").
		Offset(pagination.Offset).
		Limit(pagination.Limit).
		Find(dest).Error

	return pagination, total, err
}
