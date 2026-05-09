# wxcloudrun-golang

这是一个以微信小程序为主入口的“职场人格测试”项目。  
当前正式保留的是小程序端 `miniprogram/`，旧的浏览器首页已经弃用，仅保留为说明页，不再作为产品界面继续维护。

## 当前状态

- 小程序主题：`工位人格研究所`
- 答题结构：`36` 题主测试 + `9` 题风格加试
- 结果产出：`12` 门派 × `9` 风格 = `108` 个职场角色
- 结果能力：角色结果页、差异化文案、分享文案、海报生成、海报保存
- 本地服务：仍保留 Go 服务入口和原有接口，便于继续扩展接口或云部署

## 目录

```text
.
├── main.go                    服务入口
├── index.html                 已弃用的 Web 说明页
├── miniprogram                微信小程序正式前端
├── miniprogram/data/quiz.js   题库、门派、风格、108 角色映射
├── service/study_service.go   当前接口与首页说明页处理
├── db                         数据层
└── README.md
```

## 本地运行

启动服务：

```bash
go run .
```

默认监听：

```text
http://localhost:8080
```

注意：

- 这个地址现在只展示“旧 Web 已弃用”的说明页
- 真正要查看产品，请在微信开发者工具中导入 `miniprogram/`

## 小程序打开方式

在微信开发者工具中导入：

```text
/Users/liuxin/Desktop/410/wxcloudrun-golang/miniprogram
```

如果后续需要真机调试或接云接口，再根据实际情况配置合法域名。

## 数据库

如果未配置 MySQL，服务会自动回退到内存模式，适合本地开发。

启用 MySQL 时可设置：

```bash
export MYSQL_ADDRESS=127.0.0.1:3306
export MYSQL_USERNAME=root
export MYSQL_PASSWORD=your-password
export MYSQL_DATABASE=golang_demo
go run .
```

## 说明

当前项目已经从“学习记录 Web 页面 + 小程序骨架”收敛为：

- 只保留小程序测试产品作为主体验
- 浏览器首页仅作为弃用提示页

如果你后面要继续迭代，建议优先围绕小程序端做题库、结果页、分享链路和数据沉淀。
