package main

import (
	"github.com/asim/go-micro/plugins/registry/consul/v4"
	ratelimit4 "github.com/asim/go-micro/plugins/wrapper/ratelimiter/uber/v4"
	opentracing4 "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v4"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/opentracing/opentracing-go"
	cli2 "github.com/urfave/cli/v2"
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
	MICRO_SERVICE_NAME   = "go.micro.service.shop.category"
	MICRO_VERSION        = "latest"
	MICRO_ADDRESS        = "0.0.0.0:8086"
	MICRO_QPS            = 100
	DOCKER_HOST          = "127.0.0.1"
	MICRO_CONSUL_ADDRESS = "127.0.0.1:8500"
	MICRO_JAEGER_ADDRESS = "127.0.0.1:6831"
)

func SetDockerHost(host string) {
	DOCKER_HOST = host
	MICRO_CONSUL_ADDRESS = host + ":8500"
	MICRO_JAEGER_ADDRESS = host + ":6831"
}

func main() {
	function := micro.NewFunction(
		micro.Flags(
			&cli2.StringFlag{ //micro 多个选项 --ip
				Name:  "ip",
				Usage: "docker Host IP(ubuntu)",
				Value: "0.0.0.0",
			},
		),
	)

	function.Init(
		micro.Action(func(c *cli2.Context) error {
			ipstr := c.Value("ip").(string)
			if len(ipstr) > 0 { //后续校验IP
				log.Info("docker Host IP(ubuntu)1111", ipstr)
			}
			SetDockerHost(ipstr)
			return nil
		}),
	)

	log.Info("DOCKER_HOST ", DOCKER_HOST)

	//配置中心
	consulConfig, err := common.GetConsulConfig(MICRO_CONSUL_ADDRESS, "/micro/config")
	if err != nil {
		log.Error(err)
	}
	//注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			MICRO_CONSUL_ADDRESS,
		}
	})
	//链路追踪
	t, io, err := common.NewTracer(MICRO_SERVICE_NAME, MICRO_JAEGER_ADDRESS)
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	// Create service
	srv := micro.NewService(
		micro.Name(MICRO_SERVICE_NAME),
		micro.Version(MICRO_VERSION),
		//这里设置地址和需要暴露的端口
		micro.Address(MICRO_ADDRESS),
		//添加consul 作为注册中心
		micro.Registry(consulRegistry),
		//绑定链路追踪 服务端绑定handle 客户端绑定client
		micro.WrapHandler(opentracing4.NewHandlerWrapper(opentracing.GlobalTracer())),
		//添加限流
		micro.WrapHandler(ratelimit4.NewHandlerWrapper(MICRO_QPS)),
	)

	//获取mysql配置,路径中不带前缀
	mysqlInfo := common.GetMysqlFromConsul(consulConfig, "mysql")

	//连接数据库
	db, err := gorm.Open("mysql", mysqlInfo.User+":"+mysqlInfo.Pwd+"@tcp("+mysqlInfo.Host+":"+mysqlInfo.Port+")/"+mysqlInfo.Database+"?charset=utf8&parseTime=True&loc=Local")
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
