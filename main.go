package main

import (
	"com.jumbo/ginessential/common"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"os"
)

func main()  {
	// 读取配置
	InitConfig()

	// 数据库
	db := common.InitDB()
	defer db.Close()

	// 路由
	r := gin.Default()
	r = CollectRoute(r)
	port := viper.GetString("server.port")
	if port != "" {
		panic(r.Run(":" + port))
	}
	panic(r.Run())
}

func InitConfig()  {
	workDic, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDic + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}