## 1. client-go 简介

​ client-go 是一个调用 kubernetes 集群资源对象 API 的客户端，即通过 client-go 实现对 kubernetes 集群中资源对象（包括 deployment、service、ingress、replicaSet、pod、namespace、node 等）的增删改查等操作。大部分对 kubernetes 进行前置 API 封装的二次开发都通过 client-go 这个第三方包来实现。

​ client-go 官方文档：https://github.com/kubernetes/client-go

## 2. client-go 的使用

### 2.1 示例代码

```shell
git clone https://github.com/huweihuang/client-go.git
cd client-go
#保证本地HOME目录有配置kubernetes集群的配置文件
go run client-go.go
```
