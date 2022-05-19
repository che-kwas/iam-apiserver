#!/usr/bin/bash

# iam-apiserver smoke testing

SERVER_ADDR=localhost:8000

Header="-HContent-Type: application/json"
CCURL="curl -s -XPOST" # Create
UCURL="curl -s -XPUT" # Update
RCURL="curl -s -XGET" # Retrieve
DCURL="curl -s -XDELETE" # Delete

test::admin_login()
{
  basicToken="-HAuthorization: Basic Y2hlOmNoZS1rd2FzLmdpdGVlLmlv"
  ${CCURL} "${Header}" "${basicToken}" http://${SERVER_ADDR}/login | grep -Po '(?<=token":")(.+)(?=")'
}

test::create_admin()
{
  # ${CCURL} "${Header}" http://${SERVER_ADDR}/v1/users \
  #   -d'{"password":"che-kwas.gitee.io","username":"che","email":"che@kwas.com","phone":"17700001111","isAdmin":true}'


  token="-HAuthorization: Bearer $(test::admin_login)"
  ${DCURL} -v "${token}" http://${SERVER_ADDR}/v1/users/tom
}

test::user()
{
  echo -e '\033[32m/v1/user test begin========\033[0m'

  token="-HAuthorization: Bearer $(test::admin_login)"

  # 1. 如果有tom、jerry、john用户先清空
  echo -e '\033[32m1. delete users\033[0m'
  ${DCURL} "${token}" http://${SERVER_ADDR}/v1/users/tom
  ${DCURL} "${token}" http://${SERVER_ADDR}/v1/users/jerry
  ${DCURL} "${token}" http://${SERVER_ADDR}/v1/users/john

  # 2. 创建tom、jerry、john用户
  echo -e '\033[32m2. create users\033[0m'
  ${CCURL} "${Header}" "${token}" http://${SERVER_ADDR}/v1/users \
    -d'{"password":"tomtom","username":"tom","email":"tom@gmail.com","phone":"1812884xxxx"}'
  ${CCURL} "${Header}" "${token}" http://${SERVER_ADDR}/v1/users \
    -d'{"password":"jerryjerry","username":"jerry","email":"jerry@gmail.com","phone":"1812884xxxx"}'
  ${CCURL} "${Header}" "${token}" http://${SERVER_ADDR}/v1/users \
    -d'{"password":"johnjohn","username":"john","email":"john@gmail.com","phone":"1812884xxxx"}'

  # 3. 列出所有用户
  echo -e '\033[32m3. list users\033[0m'
  ${RCURL} "${token}" "http://${SERVER_ADDR}/v1/users?offset=0&limit=10"

  # 4. 获取tom用户的详细信息
  echo -e '\033[32m4. get tom\033[0m'
  ${RCURL} "${token}" http://${SERVER_ADDR}/v1/users/tom

  # 5. 修改tom用户
  echo -e '\033[32m5. update tom\033[0m'
  ${UCURL} "${Header}" "${token}" http://${SERVER_ADDR}/v1/users/tom \
    -d'{"username":"tom","email":"tom_modified@gmail.com","phone":"1812884xxxx"}'

  # 6. 删除tom用户
  echo -e '\033[32m6. delete tom\033[0m'
  ${DCURL} "${token}" http://${SERVER_ADDR}/v1/users/tom

  # 7. 批量删除用户
  echo -e '\033[32m7. delete users\033[0m'
  ${DCURL} "${token}" "http://${SERVER_ADDR}/v1/users?name=jerry&name=john"

  echo -e '\033[32m/v1/user test end==========\033[0m'
}

test::secret()
{
  echo -e '\033[32m/v1/secret test begin========\033[0m'

  token="-HAuthorization: Bearer $(test::admin_login)"

  # 1. 如果有secret0密钥先清空
  ${DCURL} "${token}" http://${SERVER_ADDR}/v1/secrets/secret0

  # 2. 创建secret0密钥
  ${CCURL} "${Header}" "${token}" http://${SERVER_ADDR}/v1/secrets \
    -d'{"metadata":{"name":"secret0"},"expires":0,"description":"admin secret"}'

  # 3. 列出所有密钥
  ${RCURL} "${token}" http://${SERVER_ADDR}/v1/secrets

  # 4. 获取secret0密钥的详细信息
  ${RCURL} "${token}" http://${SERVER_ADDR}/v1/secrets/secret0

  # 5. 修改secret0密钥
  ${UCURL} "${Header}" "${token}" http://${SERVER_ADDR}/v1/secrets/secret0 \
    -d'{"expires":0,"description":"admin secret(modified)"}'

  # 6. 删除secret0密钥
  ${DCURL} "${token}" http://${SERVER_ADDR}/v1/secrets/secret0

  echo -e '\033[32m/v1/secret test end==========\033[0m'
}

test::policy()
{
  echo -e '\033[32m/v1/policy test begin========\033[0m'

  token="-HAuthorization: Bearer $(test::admin_login)"

  # 1. 如果有policy0策略先清空
  ${DCURL} "${token}" http://${SERVER_ADDR}/v1/policies/policy0

  # 2. 创建policy0策略
  ${CCURL} "${Header}" "${token}" http://${SERVER_ADDR}/v1/policies \
    -d'{"metadata":{"name":"policy0"},"policy":{"description":"One policy to rule them all.","subjects":["users:<peter|ken>","users:maria","groups:admins"],"actions":["delete","<create|update>"],"effect":"allow","resources":["resources:articles:<.*>","resources:printer"],"conditions":{"remoteIPAddress":{"type":"CIDRCondition","options":{"cidr":"192.168.0.1/16"}}}}}'

  # 3. 列出所有策略
  ${RCURL} "${token}" http://${SERVER_ADDR}/v1/policies

  # 4. 获取policy0策略的详细信息
  ${RCURL} "${token}" http://${SERVER_ADDR}/v1/policies/policy0

  # 5. 修改policy0策略
  ${UCURL} "${Header}" "${token}" http://${SERVER_ADDR}/v1/policies/policy0 \
    -d'{"policy":{"description":"One policy to rule them all(modified).","subjects":["users:<peter|ken>","users:maria","groups:admins"],"actions":["delete","<create|update>"],"effect":"allow","resources":["resources:articles:<.*>","resources:printer"],"conditions":{"remoteIPAddress":{"type":"CIDRCondition","options":{"cidr":"192.168.0.1/16"}}}}}'

  # 6. 删除policy0策略
  ${DCURL} "${token}" http://${SERVER_ADDR}/v1/policies/policy0

  echo -e '\033[32m/v1/policy test end==========\033[0m'
}

test::create_admin
# test::user
# test::secret
# test::policy
