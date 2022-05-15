[![CI](https://github.com/che-kwas/iam-apiserver/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/che-kwas/iam-apiserver/actions/workflows/ci.yaml)

# iam-apiserver

## 测试

### 健康检查

```sh
# request
curl -X GET http://localhost:8000/healthz

# response
{"status":"OK"}
```

### 登录

```sh
# request
curl -X POST http://localhost:8000/login \
    -H "Authorization: Basic dGVzdDo3NzQ0MTEK"

# response
{"status":"OK"}
```
