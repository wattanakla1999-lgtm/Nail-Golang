package service

import (
	"nailly-back-end/internal/apperror"
	"nailly-back-end/internal/model"
	"nailly-back-end/internal/repository"
	"nailly-back-end/pkg/utils"
)

type NailTechnicianService struct {
	repo *repository.NailTechnicianRepository
}

func NewNailTechnicianService(repo *repository.NailTechnicianRepository) *NailTechnicianService {
	return &NailTechnicianService{repo: repo}
}

func (s *NailTechnicianService) GetNailTechnicians(filter repository.NailTechnicianFilter, pagination utils.Pagination) ([]model.NailTechnician, int64, error) {
	return s.repo.FindAll(filter, pagination)
}

func (s *NailTechnicianService) GetNailTechnicianByID(id string) (model.NailTechnician, error) {
	return s.repo.FindByID(id)
}

func (s *NailTechnicianService) CreateNailTechnician(input model.NailTechnician) (model.NailTechnician, error) {
	if input.TechnicianName == "" {
		return model.NailTechnician{}, apperror.BadRequest("technician name is required", apperror.ErrValidation)
	}
	if input.ExperienceYears < 0 {
		return model.NailTechnician{}, apperror.BadRequest("experience years must be greater than or equal to 0", apperror.ErrValidation)
	}

	if err := s.repo.Create(&input); err != nil {
		return model.NailTechnician{}, err
	}

	return input, nil
}

func (s *NailTechnicianService) UpdateNailTechnician(id string, input model.NailTechnician) (model.NailTechnician, error) {
	technician, err := s.repo.FindByID(id)
	if err != nil {
		return model.NailTechnician{}, err
	}

	if err := s.repo.Update(&technician, input); err != nil {
		return model.NailTechnician{}, err
	}

	return technician, nil
}

func (s *NailTechnicianService) DeleteNailTechnician(id string) error {
	technician, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(&technician)
}
