package vars

import (
	"sync"

	"github.com/FarmerChillax/fakeSSH/utils"
	"gorm.io/gorm"
)

var (
	db     *gorm.DB
	dbOnce sync.Once
)

func GetDB() *gorm.DB {
	if db == nil {
		dbOnce.Do(func() {
			db = utils.NewDB()
		})
	}
	return db
}
