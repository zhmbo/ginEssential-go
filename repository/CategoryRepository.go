package repository

import (
	"com.jumbo/ginessential/common"
	"com.jumbo/ginessential/model"
	"github.com/jinzhu/gorm"
)

type CategoryRepository struct {
	DB *gorm.DB
}

type ICategoryRepository interface {
	Create(name string) (*model.Category, error)
	Update(category model.Category, name string) (*model.Category, error)
	SelectById(id int) (*model.Category, error)
	DeleteById(id int) error
	PageList(pageNum int, pageSize int) ([]model.Category, int)
}

func NewCategoryRepository() ICategoryRepository {
	db := common.GetDB()
	db.AutoMigrate(model.Category{})
	return CategoryRepository{DB: db}
}

func (c CategoryRepository) Create(name string) (*model.Category, error) {
	category := model.Category{
		Name: name,
	}

	if err := c.DB.Create(&category).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

func (c CategoryRepository) Update(category model.Category, name string) (*model.Category, error) {
	if err := c.DB.Model(&category).Update("name", name).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

func (c CategoryRepository) SelectById(id int) (*model.Category, error) {
	var category model.Category
	if err := c.DB.First(&category, id).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

func (c CategoryRepository) DeleteById(id int) error {
	if err := c.DB.Delete(model.Category{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (c CategoryRepository) PageList(pageNum int, pageSize int) ([]model.Category, int) {

	// 分页
	var categories []model.Category
	c.DB.Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&categories)

	// 前端渲染分页需要知道总数
	var total int
	c.DB.Model(model.Category{}).Count(&total)

	return categories, total
}