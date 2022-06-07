[![CI](https://github.com/che-kwas/iam-apiserver/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/che-kwas/iam-apiserver/actions/workflows/ci.yaml)

# iam-apiserver

## 接口文档

- [通用说明](./docs/generic.md)
- [认证相关接口](./docs/authentication.md)
- [用户相关接口](./docs/user.md)
- [密钥相关接口](./docs/secret.md)
- [授权策略相关接口](./docs/policy.md)

## API 概览

### 认证相关接口

| 接口名称                                              | 接口功能  |
| ----------------------------------------------------- | --------- |
| [POST /login](./docs/authentication.md#1-用户登录)    | 用户登录  |
| [POST /logout](./docs/authentication.md#2-用户登出)   | 用户登出  |
| [POST /refresh](./docs/authentication.md#2-刷新Token) | 刷新Token |

### 用户相关接口

| 接口名称                                                             | 接口功能     |
| -------------------------------------------------------------------- | ------------ |
| [POST /v1/users](./docs/user.md#1-创建用户)                          | 创建用户     |
| [DELETE /v1/users](./docs/user.md#2-批量删除用户)                    | 批量删除用户 |
| [DELETE /v1/users/:name](./docs/user.md#3-删除用户)                  | 删除用户     |
| [PUT /v1/users/:name/change_password](./docs/user.md#4-修改用户密码) | 修改用户密码 |
| [PUT /v1/users/:name](./docs/user.md#5-修改用户属性)                 | 修改用户属性 |
| [GET /v1/users/:name](./docs/user.md#6-查询用户信息)                 | 查询用户信息 |
| [GET /v1/users](./docs/user.md#7-查询用户列表)                       | 查询用户列表 |

### 密钥相关接口

| 接口名称                                                              | 接口功能            |
| --------------------------------------------------------------------- | ------------------- |
| [POST /v1/secrets](./docs/secret.md#1-创建密钥)                       | 创建密钥            |
| [DELETE /v1/secrets/:name](./docs/secret.md#2-删除密钥)               | 删除密钥            |
| [PUT /v1/secrets/:name](./docs/secret.md#3-修改密钥属性)              | 修改密钥属性        |
| [GET /v1/secrets/:name](./docs/secret.md#4-查询密钥信息)              | 查询密钥信息        |
| [GET /v1/secrets/:name/token](./docs/secret.md#5-获取密钥签名的token) | 获取密钥签名的token |
| [GET /v1/secrets](./docs/secret.md#6-查询密钥列表)                    | 查询密钥列表        |

### 策略相关接口

| 接口名称                                                      | 接口功能         |
| ------------------------------------------------------------- | ---------------- |
| [POST /v1/policies](./docs/policy.md#1-创建授权策略)          | 创建授权策略     |
| [DELETE /v1/policies](./docs/policy.md#2-批量删除授权策略)    | 批量删除授权策略 |
| [DELETE /v1/policies/:name](./docs/policy.md#3-删除授权策略)  | 删除授权策略     |
| [PUT /v1/policies/:name](./docs/policy.md#4-修改授权策略属性) | 修改授权策略属性 |
| [GET /v1/policies/:name](./docs/policy.md#5-查询授权策略信息) | 查询授权策略信息 |
| [GET /v1/policies](./docs/policy.md#6-查询授权策略列表)       | 查询授权策略列表 |

## TODO

- swagger
- test
