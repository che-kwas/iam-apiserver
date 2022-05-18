#!/usr/bin/bash

SERVER_ADDR=localhost:8000

Header="-HContent-Type: application/json"
CCURL="curl -XPOST" # Create
UCURL="curl -XPUT" # Update
RCURL="curl -XGET" # Retrieve
DCURL="curl -XDELETE" # Delete

test::login()
{
  baseToken=""

  ${CCURL} "${Header}" http://${SERVER_ADDR}/login \
    -d'{"username":"admin","password":"Admin@2021"}' | grep -Po 'token[" :]+\K[^"]+'
}

test::user()
{
  echo -e '\033[32m/v1/user test begin========\033[0m'

  token="-HAuthorization: Bearer $(test::login)"

  # 1. 如果有colin、mark、john用户先清空
  ${DCURL} "${token}" http://${SERVER_ADDR}/v1/users/colin; echo
  ${DCURL} "${token}" http://${SERVER_ADDR}/v1/users/mark; echo
  ${DCURL} "${token}" http://${SERVER_ADDR}/v1/users/john; echo

  # 2. 创建colin、mark、john用户
  ${CCURL} "${Header}" "${token}" http://${SERVER_ADDR}/v1/users \
    -d'{"password":"User@2021","metadata":{"name":"colin"},"nickname":"colin","email":"colin@foxmail.com","phone":"1812884xxxx"}'; echo

  # 3. 列出所有用户
  ${RCURL} "${token}" "http://${SERVER_ADDR}/v1/users?offset=0&limit=10"; echo

  # 4. 获取colin用户的详细信息
  ${RCURL} "${token}" http://${SERVER_ADDR}/v1/users/colin; echo

  # 5. 修改colin用户
  ${UCURL} "${Header}" "${token}" http://${SERVER_ADDR}/v1/users/colin \
    -d'{"nickname":"colin","email":"colin_modified@foxmail.com","phone":"1812884xxxx"}'; echo

  # 6. 删除colin用户
  ${DCURL} "${token}" http://${SERVER_ADDR}/v1/users/colin; echo

  # 7. 批量删除用户
  ${DCURL} "${token}" "http://${SERVER_ADDR}/v1/users?name=mark&name=john"; echo

  echo -e '\033[32m/v1/user test end==========\033[0m'
}

test::secret()
{
  echo -e '\033[32m/v1/secret test begin========\033[0m'

  token="-HAuthorization: Bearer $(test::login)"

  # 1. 如果有secret0密钥先清空
  ${DCURL} "${token}" http://${SERVER_ADDR}/v1/secrets/secret0; echo

  # 2. 创建secret0密钥
  ${CCURL} "${Header}" "${token}" http://${SERVER_ADDR}/v1/secrets \
    -d'{"metadata":{"name":"secret0"},"expires":0,"description":"admin secret"}'; echo

  # 3. 列出所有密钥
  ${RCURL} "${token}" http://${SERVER_ADDR}/v1/secrets; echo

  # 4. 获取secret0密钥的详细信息
  ${RCURL} "${token}" http://${SERVER_ADDR}/v1/secrets/secret0; echo

  # 5. 修改secret0密钥
  ${UCURL} "${Header}" "${token}" http://${SERVER_ADDR}/v1/secrets/secret0 \
    -d'{"expires":0,"description":"admin secret(modified)"}'; echo

  # 6. 删除secret0密钥
  ${DCURL} "${token}" http://${SERVER_ADDR}/v1/secrets/secret0; echo

  echo -e '\033[32m/v1/secret test end==========\033[0m'
}

test::policy()
{
  echo -e '\033[32m/v1/policy test begin========\033[0m'

  token="-HAuthorization: Bearer $(test::login)"

  # 1. 如果有policy0策略先清空
  ${DCURL} "${token}" http://${SERVER_ADDR}/v1/policies/policy0; echo

  # 2. 创建policy0策略
  ${CCURL} "${Header}" "${token}" http://${SERVER_ADDR}/v1/policies \
    -d'{"metadata":{"name":"policy0"},"policy":{"description":"One policy to rule them all.","subjects":["users:<peter|ken>","users:maria","groups:admins"],"actions":["delete","<create|update>"],"effect":"allow","resources":["resources:articles:<.*>","resources:printer"],"conditions":{"remoteIPAddress":{"type":"CIDRCondition","options":{"cidr":"192.168.0.1/16"}}}}}'; echo

  # 3. 列出所有策略
  ${RCURL} "${token}" http://${SERVER_ADDR}/v1/policies; echo

  # 4. 获取policy0策略的详细信息
  ${RCURL} "${token}" http://${SERVER_ADDR}/v1/policies/policy0; echo

  # 5. 修改policy0策略
  ${UCURL} "${Header}" "${token}" http://${SERVER_ADDR}/v1/policies/policy0 \
    -d'{"policy":{"description":"One policy to rule them all(modified).","subjects":["users:<peter|ken>","users:maria","groups:admins"],"actions":["delete","<create|update>"],"effect":"allow","resources":["resources:articles:<.*>","resources:printer"],"conditions":{"remoteIPAddress":{"type":"CIDRCondition","options":{"cidr":"192.168.0.1/16"}}}}}'; echo

  # 6. 删除policy0策略
  ${DCURL} "${token}" http://${SERVER_ADDR}/v1/policies/policy0; echo

  echo -e '\033[32m/v1/policy test end==========\033[0m'
}

test::user
test::secret
test::policy
