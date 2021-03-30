package repository

import (
	"com.jumbo/ginessential/common"
	"com.jumbo/ginessential/model"
	"com.jumbo/ginessential/vo"
	"github.com/jinzhu/gorm"
)

type IPostRepository interface {
	Create(requestPost vo.CreatePostRequest, userId uint)(*model.Post, error)
	SelectById(postId string)(*model.Post, error)
	Update(post model.Post, requestPost vo.CreatePostRequest)(*model.Post, error)
	DeleteByPost(post model.Post)error
	PageList(pageNum int, pageSize int) ([]model.Post, int)
}

type PostRepository struct {
	DB *gorm.DB
}

func NewPostRepository() IPostRepository {
	db := common.GetDB()
	db.AutoMigrate(model.Post{})
	return PostRepository{DB: db}
}

func (p PostRepository) Create(requestPost vo.CreatePostRequest, userId uint) (*model.Post, error) {
	// 创建文章
	post := model.Post{
		UserId:     userId,
		CategoryId: requestPost.CategoryId,
		Title:      requestPost.Title,
		Content:    requestPost.Content,
	}

	// 插入数据
	if err := p.DB.Create(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (p PostRepository) SelectById(postId string) (*model.Post, error) {
	var post model.Post
	if err := p.DB.Preload("Category").Where("id = ?", postId).First(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (p PostRepository) Update(post model.Post, requestPost vo.CreatePostRequest) (*model.Post, error) {
	if err := p.DB.Model(&post).Update(requestPost).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (p PostRepository) DeleteByPost(post model.Post) error {
	if err := p.DB.Delete(&post).Error; err != nil {
		return err
	}
	return nil
}

func (p PostRepository) PageList(pageNum int, pageSize int) ([]model.Post, int) {

	// 分页
	var posts []model.Post
	p.DB.Preload("Category").Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&posts)

	// 前端渲染分页需要知道总数
	var total int
	p.DB.Model(model.Post{}).Count(&total)

	return posts, total
}



