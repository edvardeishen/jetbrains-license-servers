# jetbrains license servers

## 直接运行
```shell
cd ./jetbrains-license-servers
go run 
```

## docker
+ 第一步 make build : 编译一个linux-amd64的二进制
+ 第二步 make image : 使用上一步编译的二进制生成docker image
+ 运行： docker run --rm jetbrains_license_servers 



 find license server from internet



许多服务器在其他国家, 所以运行此程序以及运行jetbrains IDE时你需要确保你的网络环境足够开放,以获得更多的服务器支持
