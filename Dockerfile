FROM node AS client

WORKDIR /data
COPY . .

RUN apt-get update && \
    wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb && \
    (dpkg -i google-chrome*.deb || apt-get -y -f install) && \
    npm install && \
    npm run resource-build && \
    npm run resource-run && \
    npm run build

FROM ubuntu AS resource

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
    apt-get install -y wget libgl1 libglib2.0-0 expect && \
    wget https://www.codeandweb.com/download/texturepacker/5.1.0/TexturePacker-5.1.0-ubuntu64.deb && \
    dpkg -i TexturePacker-5.1.0-ubuntu64.deb && \
    expect -c "\
        set timeout 20;\
        spawn TexturePacker --version;\
        expect \"Please type 'agree' if you agree with the terms above:\";\
        send \"agree\n\";\
        expect eof;\
    "

WORKDIR /data

COPY --from=client /data/resources/ ./resources/

RUN mkdir -p public/spritesheet && \
    chmod u+x resources/pack.sh && \
    ls -la && ls -la resources && \
    resources/pack.sh

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
    /go/bin/revel build -m prod -a github.com/yasshi2525/RushHour

FROM alpine

ENV APP_SECRET kO0HKDOKQRLT6y9Vo0Uk69X2nxQ1p2Ln485wrYZmxiGiR7MDHa4TBxLvwLfWojcg
RUN apk update && apk --no-cache add tzdata

WORKDIR /rushhour

COPY --from=server /rushhour/ ./
COPY --from=client /data/public/ src/github.com/yasshi2525/RushHour/public/
COPY --from=resource /data/public/spritesheet/ src/github.com/yasshi2525/RushHour/public/spritesheet/

EXPOSE 9000

ENTRYPOINT [ "./run.sh" ]