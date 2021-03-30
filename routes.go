package main

import (
	"com.jumbo/ginessential/controller"
	"com.jumbo/ginessential/middleware"
	"github.com/gin-gonic/gin"
)

func CollectRoute(r *gin.Engine) *gin.Engine  {
	// 路由
	r.Use(middleware.CORSMiddleware(),middleware.RecoverMiddleware())

	// 路由分组-用户
	userRoutes := r.Group("/api/auth")
	userController := controller.NewUserController()
	userRoutes.POST("/register", userController.Register)
	userRoutes.POST("/login", userController.Login)
	userRoutes.GET("/info", middleware.AuthMiddleware(), userController.Info)

	// 路由分组-分类
	categoryRoutes := r.Group("/api/categories")
	categoryController := controller.NewCategoryController()
	categoryRoutes.POST("", categoryController.Create)
	categoryRoutes.PUT("/:id",categoryController.Update)
	categoryRoutes.GET("/:id",categoryController.Show)
	categoryRoutes.DELETE("/:id",categoryController.Delete)
	categoryRoutes.POST("page/list", categoryController.PageList)

	// 路由分组-文章
	postRoutes := r.Group("/api/posts")
	postRoutes.Use(middleware.AuthMiddleware())
	postController := controller.NewPostController()
	postRoutes.POST("", postController.Create)
	postRoutes.PUT("/:id",postController.Update)
	postRoutes.GET("/:id",postController.Show)
	postRoutes.DELETE("/:id",postController.Delete)
	postRoutes.POST("page/list", postController.PageList)

	return r
}