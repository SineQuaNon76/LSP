package service

import (
	"GoProject/models"
	"GoProject/repositories"
)

type BarangService interface {
	Create(barang *models.Barang) error
	FindAll() ([]models.Barang, error)
	FindById(id int) (*models.Barang, error)
	Edit(barang *models.Barang) error
	Delete(id int) error
}

type barangService struct {
	barangRepo repositories.BarangRepositories
}

func NewService(repo repositories.BarangRepositories) BarangService {
	return &barangService{
		barangRepo: repo,
	}
}

func (serv *barangService) Create(barang *models.Barang) error {
	err := serv.barangRepo.Create(barang)
	if err != nil {
		return err
	}
	return nil
}


func (serv *barangService) FindAll() ([]models.Barang, error){
	return serv.barangRepo.FindAll()
}

func (serv *barangService) FindById(id int) (*models.Barang, error){
	return serv.barangRepo.FindById(id)
}

func (serv *barangService) Edit(barang *models.Barang) error {
	return serv.barangRepo.Edit(barang)
}

func (serv *barangService) Delete(id int) error {
	return serv.barangRepo.Delete(id)
}