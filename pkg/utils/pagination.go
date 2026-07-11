package utils

import (
	"strconv"

	"gorm.io/gorm"
)

type Pagination struct {
	Page   int
	Limit  int
	Offset int
}

func NewPagination(pageParam, limitParam string) Pagination {
	page, _ := strconv.Atoi(pageParam)
	limit, _ := strconv.Atoi(limitParam)

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

func Paginate(query *gorm.DB, pagination Pagination, dest any) (int64, error) {
	var total int64

	if err := query.Count(&total).Error; err != nil {
		return total, err
	}

	err := query.
		Order("id ASC").
		Offset(pagination.Offset).
		Limit(pagination.Limit).
		Find(dest).Error

	return total, err
}
