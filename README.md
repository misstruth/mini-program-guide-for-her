# wxcloudrun-golang

这是一个用 Go 写的轻量前后端项目，页面风格仿照微信小程序，用来展示“小程序开发流程说明”。内容按需求梳理、前端页面、后端接口、测试上线四个阶段组织，适合作为送给对象看的说明页或演示项目。

## 当前功能

- 前端是单页卡片式界面，手机端优先，视觉上接近小程序内容页
- 后端提供 `GET /api/guide`，返回页面所需的全部文案数据
- 服务启动不再依赖 MySQL，拉下来即可直接运行

## 本地运行

```bash
go run .
```

服务默认监听 `:80`。如果你在本地没有权限监听 80 端口，可以自行改 [main.go](/Users/liuxin/Desktop/410/wxcloudrun-golang/main.go) 里的端口。

## API

### `GET /api/guide`

返回页面渲染需要的所有内容，包括：

- 应用标题和副标题
- 开发流程阶段
- 开发清单
- 提醒文案
- 常见问题
- 页脚留言

响应示例：

```json
{
  "code": 0,
  "data": {
    "app": {
      "name": "小程序开发流程卡片"
    }
  }
}
```

## 目录

```text
.
├── index.html                前端页面
├── main.go                   服务入口
├── service/guide_service.go  页面数据接口
├── service/counter_service.go 保留的模板代码，当前未使用
└── db                        原模板数据库代码，当前未接入主流程
```

## License

[MIT](./LICENSE)
