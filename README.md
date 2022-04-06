go install go-micro.dev/v4/cmd/micro@master

micro new service github.com/wsjcko/shopcategory

mkdir -p domain/{model,repository,service} 
mkdir -p protobuf/{pb,pbserver} 
mkdir -p proto/{pb,pbserver}
mkdir common

go mod edit --module=github.com/wsjcko/shopcategory
go mod edit --go=1.17  

gorm 有个根据创建表sql 生成model  : gormt

清除mod下载的包
go clean -modcache


### consul 微服务注册中心和配置中心
docker search --filter is-official=true --filter stars=3 consul
docker pull consul

## 生产环境要注意数据落盘  -v /data/consul:/data/consul
docker run -d -p 8500:8500 consul:latest 

### 注册中心
"github.com/asim/go-micro/plugins/registry/consul/v4"

### 配置中心
"github.com/asim/go-micro/plugins/config/source/consul/v4"

### consul数据库配置
http://127.0.0.1:8500/ui/dc1/kv/create

key: micro/config/mysql

{
  "host":"127.0.0.1",
  "user":"root",
  "pwd":"123456",
  "database":"shopdb",
  "port":3306
}


### jaeger 耶格 
[官方文档](https://www.jaegertracing.io/docs/1.32/)

docker run -d --name jaeger \
  -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
  -p 5775:5775/udp \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 14268:14268 \
  -p 9411:9411 \
  jaegertracing/all-in-one:latest

  http://127.0.0.1:16686/search


  服务注册consul、 链路追踪jaeger、集流量控制、熔断、容错，负载均衡等hystrix-go、
  监控告警Prometheus、日志接入ELK，到最后k8s部署

docker search hystrix
docker pull mlabouardy/hystrix-dashboard
docker run --name hystrix-dashboard -d -p 9002:9002 mlabouardy/hystrix-dashboard:latest