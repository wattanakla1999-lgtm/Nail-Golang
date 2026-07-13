package repository

import (
	"nailly-back-end/internal/model"
	"nailly-back-end/pkg/utils"

	"gorm.io/gorm"
)

type NailTechnicianFilter struct {
	TechnicianName string
	Phone          string
	Specialty      string
}

type NailTechnicianRepository struct {
	db *gorm.DB
}

func NewNailTechnicianRepository(db *gorm.DB) *NailTechnicianRepository {
	return &NailTechnicianRepository{db: db}
}

func (r *NailTechnicianRepository) FindAll(filter NailTechnicianFilter, pagination utils.Pagination) ([]model.NailTechnician, int64, error) {
	var technicians []model.NailTechnician

	query := r.db.Model(&model.NailTechnician{})
	query = utils.ApplyLikeFilters(query, map[string]string{
		"technician_name": filter.TechnicianName,
		"phone":           filter.Phone,
		"specialty":       filter.Specialty,
	})

	total, err := utils.Paginate(query, pagination, &technicians)
	if err != nil {
		return nil, 0, err
	}

	return technicians, total, nil
}

func (r *NailTechnicianRepository) FindByID(id string) (model.NailTechnician, error) {
	var technician model.NailTechnician
	err := r.db.First(&technician, id).Error
	return technician, err
}

func (r *NailTechnicianRepository) Create(technician *model.NailTechnician) error {
	return r.db.Create(technician).Error
}

func (r *NailTechnicianRepository) Update(technician *model.NailTechnician, input model.NailTechnician) error {
	return r.db.Model(technician).Updates(input).Error
}

func (r *NailTechnicianRepository) Delete(technician *model.NailTechnician) error {
	return r.db.Delete(technician).Error
}
