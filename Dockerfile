FROM golang:1.18.3-alpine3.16 AS build

ARG VERSION=latest
ARG OUTPUT=/iam-apiserver

WORKDIR /src
COPY . .

RUN go env -w GOPROXY=https://goproxy.cn,direct \
      && go mod tidy -compat=1.18 \
      && go build -ldflags "-X main.Version=${VERSION}" -o ${OUTPUT}/ ./... \
      && cp configs/iam-apiserver.yaml ${OUTPUT}/ \
      && rm -rf /src

# ================================

FROM alpine:3.16

ENV TZ Asia/Shanghai

RUN apk add tzdata && cp /usr/share/zoneinfo/${TZ} /etc/localtime \
      && echo ${TZ} > /etc/timezone \
      && apk del tzdata

COPY --from=build /iam-apiserver/iam-apiserver /opt/iam/bin/
COPY --from=build /iam-apiserver/iam-apiserver.yaml /etc/iam/

EXPOSE 8000 8001
ENTRYPOINT [ "/opt/iam/bin/iam-apiserver" ]
CMD [ "-c", "/etc/iam/iam-apiserver.yaml" ]
