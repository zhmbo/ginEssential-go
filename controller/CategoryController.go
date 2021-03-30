package controller

import (
	"com.jumbo/ginessential/repository"
	"com.jumbo/ginessential/response"
	"com.jumbo/ginessential/vo"
	"github.com/gin-gonic/gin"
	"strconv"
)

type ICategoryController interface {
	RestController
	PageList(ctx *gin.Context)
}

type CategoryController struct {
	Repository repository.ICategoryRepository
}

func NewCategoryController() ICategoryController  {
	repository := repository.NewCategoryRepository()
	return CategoryController{Repository: repository}
}

func (c CategoryController) Create(ctx *gin.Context) {
	// 绑定 body中的参数
	var requestCategory vo.CreateCategoryRequest
	if err := ctx.ShouldBind(&requestCategory); err != nil {
		response.Fail(ctx,"数据验证错误", nil)
		return
	}
	category, err := c.Repository.Create(requestCategory.Name)
	if  err != nil {
		panic(err)
		return
	}

	response.Success(ctx, gin.H{"category": category}, "")
}

func (c CategoryController) Update(ctx *gin.Context) {
	// 绑定 body中的参数
	var requestCategory vo.CreateCategoryRequest
	if err := ctx.ShouldBind(&requestCategory); err != nil {
		response.Fail(ctx,"数据验证错误", nil)
		return
	}

	// 获取path 中参数
	categoryId, _ := strconv.Atoi(ctx.Params.ByName("id")) // Atoi string强转int
	updateCategory, err := c.Repository.SelectById(categoryId)
	if err != nil {
		response.Fail(ctx, "分类不存在", nil)
		return
	}

	// 更新分类 update类型：map，struct，name：value
	category, err := c.Repository.Update(*updateCategory, requestCategory.Name)
	if err != nil {
		response.Fail(ctx, "更新失败", nil)
		return
	}

	response.Success(ctx, gin.H{"category": category}, "修改成功")
}

func (c CategoryController) Show(ctx *gin.Context) {
	// 获取path 中参数
	categoryId, _ := strconv.Atoi(ctx.Params.ByName("id")) // Atoi string强转int

	category, err := c.Repository.SelectById(categoryId)
	if err != nil {
		response.Fail(ctx, "分类不存在", nil)
		return
	}

	response.Success(ctx, gin.H{"category": category}, "")
}

func (c CategoryController) Delete(ctx *gin.Context) {
	// 获取path 中参数
	categoryId, _ := strconv.Atoi(ctx.Params.ByName("id")) // Atoi string强转int

	if err := c.Repository.DeleteById(categoryId); err != nil {
		response.Fail(ctx, "删除失败，请重试", nil)
		return
	}

	response.Success(ctx, nil, "删除成功")
}

func (c CategoryController) PageList(ctx *gin.Context) {
	// 获取分页参数
	pageNum,_ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	categories, total := c.Repository.PageList(pageNum, pageSize)

	response.Success(ctx, gin.H{"categories": categories, "total": total}, "成功")
}

