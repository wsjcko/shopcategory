go install go-micro.dev/v4/cmd/micro@master

micro new service shopcategory

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