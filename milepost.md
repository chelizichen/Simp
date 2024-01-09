# 里程碑

## 2024.1.6

验证进程间通信、配置文件更新。

备注: 进程间通信用STD有点蠢了感觉

准备添加TCP套接字通信

## 2024.1.7

验证打包后文件的执行。添加Web界面，添加调试接口。

备注：先从HTTP形式的微服务开始。

## 2024.1.8

命令行打包

````shell
# 构建
go build -o service_go
# 压缩
tar -cvf SimpTestServer.tar.gz ./simp.yaml ./service_go
````

获取包列表、服务列表、创建服务。

## 2024.1.9

界面优化 完善重启、状态接口
构建命令可能将 service目录覆盖 ，所以采用 service_go 为二进制文件名