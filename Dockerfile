# build stage
FROM golang:1.16-alpine AS build-env
ARG APP_NAME=goapp

RUN apk add --no-cache curl bash git openssh
    
COPY . /go/src/github.com/govinda-attal/eshop
WORKDIR /go/src/github.com/govinda-attal/eshop
RUN go build -o dist/$APP_NAME

# final stage
FROM alpine:3.12
RUN apk -U add ca-certificates

WORKDIR /app
COPY --from=build-env /go/src/github.com/govinda-attal/eshop/dist/$APP_NAME /app/
COPY --from=build-env /go/src/github.com/govinda-attal/eshop/configs/app-cfg.yaml /app/configs/app-cfg.yaml
COPY --from=build-env /go/src/github.com/govinda-attal/eshop/scripts/db /app/scripts/db

EXPOSE 9080