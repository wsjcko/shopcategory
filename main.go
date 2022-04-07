package main

import (
	"github.com/asim/go-micro/plugins/registry/consul/v4"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/wsjcko/shopcategory/common"
	"github.com/wsjcko/shopcategory/domain/repository"
	"github.com/wsjcko/shopcategory/domain/service"
	"github.com/wsjcko/shopcategory/handler"
	pb "github.com/wsjcko/shopcategory/protobuf/pb"
	"go-micro.dev/v4"
	log "go-micro.dev/v4/logger"
	"go-micro.dev/v4/registry"
)

var (
	MICRO_SERVICE_NAME = "go.micro.service.shop.category"
	MICRO_VERSION      = "latest"
	MICRO_ADDRESS      = "127.0.0.1:8082"
)

func main() {
	//配置中心
	consulConfig, err := common.GetConsulConfig("127.0.0.1", 8500, "/micro/config")
	if err != nil {
		log.Error(err)
	}
	//注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})

	// Create service
	srv := micro.NewService(
		micro.Name(MICRO_SERVICE_NAME),
		micro.Version(MICRO_VERSION),
		//这里设置地址和需要暴露的端口
		micro.Address(MICRO_ADDRESS),
		//添加consul 作为注册中心
		micro.Registry(consulRegistry),
	)

	//获取mysql配置,路径中不带前缀
	mysqlInfo := common.GetMysqlFromConsul(consulConfig, "mysql")

	//连接数据库
	db, err := gorm.Open("mysql", mysqlInfo.User+":"+mysqlInfo.Pwd+"@/"+mysqlInfo.Database+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Error(err)
	}
	defer db.Close()
	//禁止复表
	db.SingularTable(true)

	//初始化创建表
	// err = repository.NewCategoryRepository(db).InitTable() //gorm 创建表  只需执行一次
	// if err != nil {
	// 	log.Fatal(err)
	// }
	categoryService := service.NewCategoryService(repository.NewCategoryRepository(db))
	srv.Init()

	// Register handler
	err = pb.RegisterShopCategoryHandler(srv.Server(), &handler.ShopCategory{CategoryService: categoryService})

	// Run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
