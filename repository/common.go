package repository

import (
	"errors"
	"time"

	"github.com/FarmerChillax/fakeSSH/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrModelNotBeNil         = errors.New("model not be nil")
	ErrIdLessThanZero        = errors.New("id not less than zero")
	ErrModelsNotBeEmptySlice = errors.New("models not be empty slice")
	ErrConditionNotBeNil     = errors.New("condition not be nil")
	ErrIdsNotBeEmpty         = errors.New("ids not be empty array")
)

type Model interface {
	TableName() string
}

type CommonRepository[M Model] struct {
	db *gorm.DB
}

// 创建记录
func (r *CommonRepository[M]) Create(model *M) error {
	if model == nil {
		return ErrModelNotBeNil
	}

	db := r.db.Create(model)
	if err := db.Error; err != nil {
		return err
	}
	return nil
}

// 批量创建记录
func (r *CommonRepository[M]) CreateBatch(models []*M) error {
	if models == nil {
		return ErrModelNotBeNil
	}
	if len(models) <= 0 {
		return ErrModelsNotBeEmptySlice
	}

	db := r.db.Create(models)
	if err := db.Error; err != nil {
		return err
	}
	return nil
}

// 通过条件更新记录
func (r *CommonRepository[M]) Update(model *M, condition *M, selectFields ...string) error {
	if model == nil {
		return ErrModelNotBeNil
	}

	if condition == nil {
		return ErrConditionNotBeNil
	}

	db := r.db
	if len(selectFields) > 0 {
		db = db.Select(selectFields)
	}

	db = db.Where(condition).Updates(model)
	if err := db.Error; err != nil {
		return err
	}
	return nil
}

// 通过 ID 更新记录
// 使用 struct 进行更新时，GORM 只会更新非零值的字段。selectFields 用于指定更新的字段，避免无法更新零值。
// more see: https://gorm.io/zh_CN/docs/update.html#%E6%9B%B4%E6%96%B0%E9%80%89%E5%AE%9A%E5%AD%97%E6%AE%B5
func (r *CommonRepository[M]) UpdateById(model *M, id int64, selectFields ...string) error {
	if model == nil {
		return ErrModelNotBeNil
	}

	if id <= 0 {
		return ErrIdLessThanZero
	}

	db := r.db
	if len(selectFields) > 0 {
		db = db.Select(selectFields)
	}

	db = db.Where("id = ?", id).Updates(model)
	if err := db.Error; err != nil {
		return err
	}
	return nil
}

// 创建及更新记录。根据数据的唯一键，不存在即创建记录，存在则更新记录
// more see: https://gorm.io/zh_CN/docs/create.html#Upsert-%E5%8F%8A%E5%86%B2%E7%AA%81
func (r *CommonRepository[M]) Upsert(models []*M, fields ...string) error {
	if models == nil {
		return ErrModelNotBeNil
	}
	if len(models) <= 0 {
		return ErrModelsNotBeEmptySlice
	}

	db := r.db
	if len(fields) > 0 {
		db = db.Clauses(clause.OnConflict{
			DoUpdates: clause.AssignmentColumns(fields),
		})
	} else {
		db = db.Clauses(clause.OnConflict{
			UpdateAll: true,
		})
	}

	db = db.Create(&models)
	if err := db.Error; err != nil {
		return err
	}
	return nil
}

// 通过条件删除记录
func (r *CommonRepository[M]) Delete(condition *M, deleteUser ...string) error {
	if condition == nil {
		return ErrConditionNotBeNil
	}

	var deleteUser0 string
	if len(deleteUser) > 0 {
		deleteUser0 = deleteUser[0]
	}

	db := r.db.Where(condition).Updates(map[string]interface{}{
		"deleted_time": time.Now().Unix(),
		"delete_user":  deleteUser0,
	})
	if err := db.Error; err != nil {
		return err
	}
	return nil
}

// 通过 IDs 删除记录
func (r *CommonRepository[M]) DeleteByIds(ids []int64, deleteUser ...string) error {
	if len(ids) <= 0 {
		return ErrIdsNotBeEmpty
	}

	var deleteUser0 string
	if len(deleteUser) > 0 {
		deleteUser0 = deleteUser[0]
	}

	db := r.db.Where(`id IN ?`, ids).Updates(map[string]interface{}{
		"deleted_time": time.Now().Unix(),
		"delete_user":  deleteUser0,
	})
	if err := db.Error; err != nil {
		return err
	}
	return nil
}

// 通过 ID 删除记录
func (r *CommonRepository[M]) DeleteById(id int64, deleteUser ...string) error {
	if id <= 0 {
		return ErrIdLessThanZero
	}

	err := r.DeleteByIds([]int64{id}, deleteUser...)
	if err != nil {
		return err
	}
	return nil
}

// 通过条件查询记录
func (r *CommonRepository[M]) Get(condition *M, fields ...string) (*M, error) {
	var record M
	db := r.db.Where(condition)
	if len(fields) > 0 {
		db = db.Select(fields)
	}
	db = db.First(&record)
	if err := db.Error; err != nil {
		return nil, err
	}
	return &record, nil
}

// 通过 ID 查询记录
func (r *CommonRepository[M]) GetById(id int64, fields ...string) (*M, error) {
	if id <= 0 {
		return nil, ErrIdLessThanZero
	}

	var record M
	db := r.db.Where("id = ?", id)
	if len(fields) > 0 {
		db = db.Select(fields)
	}
	db = db.First(&record)
	if err := db.Error; err != nil {
		return nil, err
	}
	return &record, nil
}

// 通过条件查询是否存在记录
func (r *CommonRepository[M]) Exist(condition *M) (bool, error) {
	_, err := r.Get(condition)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

// 通过条件查询是否存在记录，排除指定 IDs
func (r *CommonRepository[M]) ExistExcludeIds(condition *M, ids []int64) (bool, error) {
	if len(ids) <= 0 {
		return false, ErrIdsNotBeEmpty
	}

	var record M
	db := r.db.Where("id NOT IN ?", ids).Where(condition)
	db = db.First(&record)
	if err := db.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

// 通过条件查询是否存在记录，排除指定 ID
func (r *CommonRepository[M]) ExistExcludeId(condition *M, id int64) (bool, error) {
	if id <= 0 {
		return false, ErrIdLessThanZero
	}
	return r.ExistExcludeIds(condition, []int64{id})
}

// 通过条件查询记录列表
func (r *CommonRepository[M]) List(condition *M, fields ...string) ([]*M, error) {
	var records []*M
	db := r.db.Where(condition)
	if len(fields) > 0 {
		db = db.Select(fields)
	}
	db = db.Find(&records)
	if err := db.Error; err != nil {
		return records, err
	}
	return records, nil
}

// 通过条件查询记录列表（包含删除的记录）
func (r *CommonRepository[M]) ListAll(condition *M, fields ...string) ([]*M, error) {
	var records []*M
	db := r.db.Unscoped().Where(condition)
	if len(fields) > 0 {
		db = db.Select(fields)
	}
	db = db.Find(&records)
	if err := db.Error; err != nil {
		return records, err
	}
	return records, nil
}

// 通过条件分页查询记录列表
func (r *CommonRepository[M]) Page(condition *M, pager *model.Pager, fields ...string) ([]*M, error) {
	records, err := r.PageOrder(condition, pager, "modify_time desc", fields...)
	if err != nil {
		return nil, err
	}
	return records, nil
}

// 通过条件分页查询记录列表（指定排序）
func (r *CommonRepository[M]) PageOrder(condition *M, pager *model.Pager, order string, fields ...string) ([]*M, error) {
	var records []*M
	db := r.db.Where(condition)
	if len(fields) > 0 {
		db = db.Select(fields)
	}
	if len(order) > 0 {
		db = db.Order(order)
	}
	db.Count(&pager.TotalRows)
	db = db.Offset(getOffset(pager)).Limit(int(pager.PageSize)).Find(&records)
	if err := db.Error; err != nil {
		return records, err
	}
	return records, nil
}

func getOffset(pager *model.Pager) int {
	if pager.Page <= 0 {
		return 0
	}
	return int((pager.Page - 1) * pager.PageSize)
}

func (r *CommonRepository[M]) FirstOrCreate(condition *M) (*M, error) {
	var model M
	db := r.db.FirstOrCreate(&model, condition)
	if err := db.Error; err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *CommonRepository[M]) Count(condition *M, fields ...string) (int64, error) {
	var count int64
	db := r.db
	if len(fields) > 0 {
		db = db.Select(fields)
	}
	db = db.Where(condition).Count(&count)
	if err := db.Error; err != nil {
		return count, err
	}
	return count, nil
}
