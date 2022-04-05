go install go-micro.dev/v4/cmd/micro@master

micro new service shopcategory

mkdir -p domain/{model,repository,service} 
mkdir -p protobuf/{pb,pbserver} 
mkdir -p proto/{pb,pbserver}
mkdir common

go mod edit --module=github.com/wsjcko/shopcategory
go mod edit --go=1.17  

gorm 有个根据创建表sql 生成model  : gormt