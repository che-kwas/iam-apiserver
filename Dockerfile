FROM golang:1.17.6-alpine3.15 AS build

ARG VERSION=latest
ARG OUTPUT=/iam-apiserver

WORKDIR /src
COPY . .

RUN go env -w GOPROXY=https://goproxy.cn,direct \
      && go mod tidy -compat=1.17 \
      && go build -ldflags "-X main.Version=${VERSION}" -o ${OUTPUT}/ ./... \
      && cp configs/iam-apiserver.yaml ${OUTPUT}/ \
      && rm -rf /src

# ================================

FROM alpine:3.15

ENV TZ Asia/Shanghai

RUN apk add tzdata && cp /usr/share/zoneinfo/${TZ} /etc/localtime \
      && echo ${TZ} > /etc/timezone \
      && apk del tzdata

COPY --from=build /iam-apiserver/* /

EXPOSE 8000 8001
ENTRYPOINT [ "/iam-apiserver" ]
CMD [ "-c", "/iam-apiserver.yaml" ]