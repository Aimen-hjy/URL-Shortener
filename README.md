# URL Shortener
使用Golang开发的一个简单的URL短链接服务。（Go beginner project）
## 开发计划
- [x] V1-基本功能：生成短链接，用内存存储映射（map），重定向到原始URL
- [x] V2-Web界面：学习使用Gin框架进行开发
- [ ] V3-持久化存储：使用数据库存储短链接和原始URL的映射关系
## 使用说明
先在本地安装Go环境，安装Gin框架，SQlite3数据库驱动，然后直接运行main.go文件，访问http://localhost:8080
```bash
go get -u github.com/gin-gonic/gin
go get github.com/mattn/go-sqlite3
```
## 参考仓库
> 本项目使用了以下仓库
> - [Gin框架](https://github.com/gin-gonic/gin)
> - [SQLite3](https://github.com/sqlite/sqlite)