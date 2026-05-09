# 工位人格研究所

一个偏搞笑、偏真实、专戳打工人痛点的职场人格测试小程序。

当前项目已经收敛为微信小程序 MVP，核心是：

- `36` 题主测试 + `9` 题风格加试
- `12` 类工作驱动力 × `9` 种表达风格 = `108` 个职场角色
- 支持结果页、分享文案、海报生成与保存

## 当前结构

```text
.
├── main.go                    Go 服务入口
├── index.html                 本地预览说明页
├── miniprogram                微信小程序正式前端
├── miniprogram/data/quiz.js   题库、维度、角色映射
├── service/study_service.go   页面与接口逻辑
├── db                         数据层
└── README.md
```

## 本地启动

```bash
go run .
```

默认地址：

```text
http://localhost:8080
```

这个地址主要用于本地说明和前端调试预览。正式体验入口仍然是小程序。

## 打开小程序

用微信开发者工具导入下面这个目录：

```text
/Users/liuxin/Desktop/410/wxcloudrun-golang/miniprogram
```

小程序当前主题名为：`工位人格研究所`

## 数据库说明

如果没有配置 MySQL，服务会自动回退到内存模式，适合本地开发。

需要连 MySQL 时可设置：

```bash
export MYSQL_ADDRESS=127.0.0.1:3306
export MYSQL_USERNAME=root
export MYSQL_PASSWORD=your-password
export MYSQL_DATABASE=golang_demo
go run .
```

## 现在这版在做什么

- 只保留小程序测试作为主产品
- 浏览器页面仅承担预览和说明作用
- 后端保留扩展接口能力，方便继续加埋点、结果存档、分享数据等

## 后续适合继续迭代的方向

- 扩题库，把结果分布拉得更细
- 做用户答题记录与画像沉淀
- 加分享裂变、海报模板和结果榜单
- 接通云开发或正式后端存储
