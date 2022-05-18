# 通用说明

IAM 系统 API 严格遵循 REST 标准进行设计，采用 JSON 格式进行数据传输，使用 JWT Token 进行 API 认证。

## 1. 公共参数

公共参数，是每个接口都需要传入的，在每个接口文档中，不再一一说明。

IAM API 接口公共参数如下：

| 参数名称      | 位置   | 类型   | 必选 | 描述                          |
| ------------- | ------ | ------ | ---- | ----------------------------- |
| Content-Type  | Header | String | 是   | 固定值：application/json      |
| Authorization | Header | String | 是   | JWT Token，值以 `Bearer` 开头 |

## 2. 返回结果

| 名称         | 位置   | 描述                |
| ------------ | ------ | ------------------- |
| HTTP状态码   | Header | HTTP状态码          |
| X-Request-Id | Header | 用于识别一次Request |
| 响应数据     | Body   | JSON格式的响应数据  |


### 2.1 成功返回结果

成功时返回的 HTTP 状态码为`200`，在 Body 中返回数据，以下是创建密钥 API 接口的返回结果：

```json
{
  "metadata": {
    "id": 24,
    "name": "secretdemo",
    "createdAt": "2020-09-20T10:17:58.108812081+08:00",
    "updatedAt": "2020-09-20T10:17:58.108812081+08:00"
  },
  "username": "admin",
  "secretID": "k5jZYMJCAk4jGH1nqgszTn6hPaZ8aZbKO0ZO",
  "secretKey": "cKdfmDJlTELfumu3SpLPf0k0SXQDqvdJ",
  "expires": 0,
  "description": "admin secret"
}
```

### 2.2 失败返回结果

失败时返回的 HTTP 状态码是 `400、401、403、404、500` 中的一个，以下是创建重复密钥时，API 接口返回的错误结果：

```json
{
  "code": 100101,
  "message": "Database error"
}
```

## 3. 错误码

| Identifier          | Code   | HTTP Code | Description                  |
| ------------------- | ------ | --------- | ---------------------------- |
| ErrSuccess          | 100001 | 200       | OK                           |
| ErrUnknown          | 100002 | 500       | Internal server error        |
| ErrBadParams        | 100003 | 400       | Bad request parameters       |
| ErrNotFound         | 100004 | 404       | Not found                    |
| ErrPasswordInvalid  | 100101 | 401       | Password invalid             |
| ErrHeaderInvalid    | 100102 | 401       | Authorization header invalid |
| ErrSignatureInvalid | 100103 | 401       | Signature invalid            |
| ErrTokenInvalid     | 100104 | 401       | Token invalid                |
| ErrTokenExpired     | 100105 | 401       | Token expired                |
| ErrPermissionDenied | 100106 | 403       | Permission denied            |
| ErrDatabase         | 100201 | 500       | Database error               |
| ErrUserNotFound     | 110001 | 404       | User not found               |
| ErrUserAlreadyExist | 110002 | 400       | User already exist           |
| ErrSecretNotFound   | 110101 | 404       | Secret not found             |
| ErrReachMaxCount    | 110102 | 400       | Secret reach the max count   |
| ErrPolicyNotFound   | 110201 | 404       | Policy not found             |
