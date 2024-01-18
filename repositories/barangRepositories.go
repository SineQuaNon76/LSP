package repositories

import (
	"GoProject/models"

	"gorm.io/gorm"
)

type BarangRepositories interface {
	Create(barang *models.Barang) error
	FindAll() ([]models.Barang, error)
	FindById(id int) (*models.Barang, error)
	Edit(barang *models.Barang) error
	Delete(id int) error
}

type repositories struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) BarangRepositories {
	return &repositories{db: db}
}

func (repo *repositories) Create(barang *models.Barang) error {
	return repo.db.Create(barang).Error
}

func (repo *repositories) FindAll() ([]models.Barang, error) {
	var barang []models.Barang

	err := repo.db.Find(&barang).Error

	return barang, err
}

func (repo *repositories) FindById(id int) (*models.Barang, error) {
	var barangg models.Barang
	if err := repo.db.First(&barangg, id).Error; err != nil {
		return nil, err
	}
	return &barangg, nil
}

func (repo *repositories) Edit(barang *models.Barang) error {
	return repo.db.Save(barang).Error
}


func (repo *repositories) Delete(id int) error {
	return repo.db.Delete(&models.Barang{}, id).Error
}

