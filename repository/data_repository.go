package repository

import (
	"github.com/FarmerChillax/fakeSSH/model"
	"gorm.io/gorm"
)

type DataRepository struct {
	CommonRepository[model.Data]
}

func NewDataRepository(db *gorm.DB) *DataRepository {
	return &DataRepository{
		CommonRepository: CommonRepository[model.Data]{
			db: db.Table(model.Data{}.TableName()),
		},
	}
}
