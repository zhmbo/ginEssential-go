package controller

import (
	"com.jumbo/ginessential/model"
	"com.jumbo/ginessential/repository"
	"com.jumbo/ginessential/response"
	"com.jumbo/ginessential/vo"
	"github.com/gin-gonic/gin"
	"strconv"
)

type IPostController interface {
	RestController
	PageList(ctx *gin.Context)
}

type PostController struct {
	Repository repository.IPostRepository
}

func NewPostController() IPostController {
	repository := repository.NewPostRepository()
	return PostController{Repository: repository}
}

func (p PostController) Create(ctx *gin.Context) {
	// 绑定 body中的参数
	var requestPost vo.CreatePostRequest
	if err := ctx.ShouldBind(&requestPost); err != nil {
		response.Fail(ctx,"数据验证错误", nil)
		return
	}

	// 获取登录用户 user
	user, _ := ctx.Get("user")

	post, err := p.Repository.Create(requestPost, user.(model.User).ID)
	if  err != nil {
		panic(err)
		return
	}

	response.Success(ctx, gin.H{"post": post}, "创建成功")
}

func (p PostController) Update(ctx *gin.Context) {
	// 绑定 body中的参数
	var requestPost vo.CreatePostRequest
	if err := ctx.ShouldBind(&requestPost); err != nil {
		response.Fail(ctx,"数据验证错误", nil)
		return
	}

	// 获取path 中参数
	postId := ctx.Params.ByName("id")
	updatePost, err := p.Repository.SelectById(postId)
	if err != nil {
		response.Fail(ctx, "文章不存在", nil)
		return
	}

	// 判断当前用户是否为文章作者
	// 获取登录用户 user
	user, _ := ctx.Get("user")
	userId := user.(model.User).ID
	if userId != updatePost.UserId {
		response.Fail(ctx, "您没有权限修改此文章，请勿非法操作", nil)
		return
	}

	// 更新文章
	post, err := p.Repository.Update(*updatePost, requestPost)
	if err != nil {
		response.Fail(ctx, "更新失败", nil)
		return
	}

	response.Success(ctx, gin.H{"post": post}, "修改成功")
}

func (p PostController) Show(ctx *gin.Context) {
	// 获取path 中参数
	postId := ctx.Params.ByName("id")

	post, err := p.Repository.SelectById(postId)
	if err != nil {
		response.Fail(ctx, "文章不存在", nil)
		return
	}

	response.Success(ctx, gin.H{"post": post}, "")
}

func (p PostController) Delete(ctx *gin.Context) {
	// 获取path 中参数
	postId := ctx.Params.ByName("id")
	post, err := p.Repository.SelectById(postId)
	if err != nil {
		response.Fail(ctx, "文章不存在", nil)
		return
	}

	// 判断当前用户是否为文章作者
	// 获取登录用户 user
	user, _ := ctx.Get("user")
	userId := user.(model.User).ID
	if userId != post.UserId {
		response.Fail(ctx, "您没有权限修改此文章，请勿非法操作", nil)
		return
	}

	if err := p.Repository.DeleteByPost(*post); err != nil {
		response.Fail(ctx, "删除失败，请重试", nil)
		return
	}

	response.Success(ctx, nil, "删除成功")
}

func (p PostController) PageList(ctx *gin.Context) {
	// 获取分页参数
	pageNum,_ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	posts, total := p.Repository.PageList(pageNum, pageSize)

	response.Success(ctx, gin.H{"posts": posts, "total": total}, "成功")
}