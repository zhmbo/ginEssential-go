package controller

import (
	"com.jumbo/ginessential/common"
	"com.jumbo/ginessential/dto"
	"com.jumbo/ginessential/model"
	"com.jumbo/ginessential/response"
	"com.jumbo/ginessential/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

type IUserController interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	Info(ctx *gin.Context)
}

type UserController struct {
	DB *gorm.DB
}

func NewUserController() IUserController {
	db := common.DB
	// 自动迁移
	db.AutoMigrate(&model.User{})
	return UserController{DB: db}
}

func (u UserController) Register(ctx *gin.Context) {
	// 使用map 获取请求参数
	//var requestMap = make(map[string]string)
	//json.NewDecoder(ctx.Request.Body).Decode(&requestMap)

	// 结构体获取请求数据
	var requestUser model.User
	//json.NewDecoder(ctx.Request.Body).Decode(&requestUser)

	// gin Bind获取请求参数
	ctx.Bind(&requestUser)

	// 获取参数
	name := requestUser.Name
	telephone := requestUser.Telephone
	password := requestUser.Password

	// 验证手机号
	if len(telephone) != 11 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位")
		return
	}

	// 验证密码
	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}

	// 名称为空，随机名称
	if len(name) == 0 {
		name = util.RandomString(10)
	}

	log.Println(name, telephone, password)
	//判断手机号是否存在
	if isTelephoneExist(u.DB, telephone) {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户已存在")
		return
	}

	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "加密错误")
		return
	}
	// 创建用户
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hasedPassword),
	}
	u.DB.Create(&newUser)


	// 发放token
	token, err := common.ReleaseTken(newUser)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "系统错误")
		log.Panicf("token generate error : %v", err)
		return
	}

	response.Success(ctx, gin.H{"token": token}, "注册成功")
}

func (u UserController) Login(ctx *gin.Context) {
	// 获取参数
	var requestUser model.User
	// gin Bind获取请求参数
	ctx.Bind(&requestUser)

	// 获取参数
	telephone := requestUser.Telephone
	password := requestUser.Password

	log.Println("手机号：", telephone)

	// 验证手机号
	if len(telephone) != 11 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位")
		return
	}

	// 验证密码
	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}

	// 判断手机号是否存在
	var user model.User
	u.DB.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户不存在")
		return
	}

	// 判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		response.Fail(ctx, "密码错误", nil)
		return
	}

	// 发放token
	token, err := common.ReleaseTken(user)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "系统错误")
		log.Panicf("token generate error : %v", err)
		return
	}
	// 返回结果
	response.Success(ctx, gin.H{"token": token}, "登陆成功")
}

func (u UserController) Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	data := gin.H{"user": dto.ToUserDto(user.(model.User))}
	response.Success(ctx, data, "")
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}