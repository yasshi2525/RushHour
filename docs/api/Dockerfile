FROM golang:alpine as base

WORKDIR /work

COPY . .

RUN apk update && apk add --no-cache git wget && \
    # go get swagに失敗するため、一時的にコメントアウト
    #     go get -u github.com/swaggo/swag/cmd/swag && \
    wget https://github.com/swaggo/swag/releases/download/v1.6.3/swag_1.6.3_Linux_x86_64.tar.gz && \
    tar xvf swag_1.6.3_Linux_x86_64.tar.gz && \
    ./swag init -o docs/api

FROM swaggerapi/swagger-ui

COPY --from=base /work/docs/api/ /app/