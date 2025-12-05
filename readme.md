# readme.md

## 项目概述

基于 Go 1.25.0 的微服务单体仓库，实现 管理平台。

## 目录结构

- `backend` 后端服务目录
- `constants` 常量定义目录
- `models` 数据模型目录
- `protocol` 协议定义目录
- `main.go` 后端服务入口文件


## 服务实现
- `backend\services\platform_service` 平台服务实现目录
- `backend\services\user_service` 用户服务实现目录
- `app.go`  后端服务入口
- `app_middleware.go` 后端服务中间件

## 需要实现功能

- 用户管理与认证
- 平台配置管理
- 日志记录与监控
- 数据库交互与缓存
- 模型管理服务