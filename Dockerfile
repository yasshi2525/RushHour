FROM node AS client

WORKDIR /data
COPY . .

RUN npm install && \
    npm run build

FROM golang:alpine as server

WORKDIR /go

RUN mkdir -p src/github.com/yasshi2525/RushHour
COPY . src/github.com/yasshi2525/RushHour
RUN sed -i -e "s|conf/game.conf|src/github.com/yasshi2525/RushHour/conf/game.conf|" src/github.com/yasshi2525/RushHour/app/services/config.go && \
    sed -i -e "s/&loc=Asia%2FTokyo//" src/github.com/yasshi2525/RushHour/conf/app.conf

RUN apk update && apk add --no-cache git

RUN go get gopkg.in/go-playground/validator.v9 && \
    go get github.com/BurntSushi/toml && \
    go get github.com/jinzhu/gorm && \
    go get github.com/go-sql-driver/mysql && \
    go get github.com/revel/revel && \
    go get github.com/revel/cmd/revel

RUN mkdir -p /rushhour && \
    cd /rushhour && \
    /go/bin/revel build -a github.com/yasshi2525/RushHour

FROM alpine

RUN apk update && apk --no-cache add tzdata

WORKDIR /rushhour

COPY --from=server /rushhour/ ./
COPY --from=client /data/public/ src/github.com/yasshi2525/RushHour/public/

EXPOSE 9000

ENTRYPOINT [ "./run.sh" ]