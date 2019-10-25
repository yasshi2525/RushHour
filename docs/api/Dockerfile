FROM golang:alpine as base

WORKDIR /go

COPY /app ./src/github.com/yasshi2525/RushHour/app

RUN apk update && apk add --no-cache git && \
    go get -u github.com/swaggo/swag/cmd/swag && \
    cd ./src/github.com/yasshi2525/RushHour && \
    swag init -g ./app/controllers/v1/api.go

FROM swaggerapi/swagger-ui

COPY --from=base /go/src/github.com/yasshi2525/RushHour/docs/ /app/