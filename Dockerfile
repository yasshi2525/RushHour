FROM node AS client

WORKDIR /data
COPY . .

RUN apt-get update && apt-get install -y git && \
    npm ci && \
    npm run build && \
    git clone https://github.com/yasshi2525/RushHourResource.git && \
    mkdir -p public/spritesheet && \
    cp -r RushHourResource/dist/* public/spritesheet/

FROM golang:alpine as server

WORKDIR /go

ENV baseurl "https://localhost:9000/"
ENV salt ""
ENV key "1234567890123456"
ENV state ""
ENV twitter_token ""
ENV twitter_secret ""
ENV google_client ""
ENV google_secret ""

RUN mkdir -p src/github.com/yasshi2525/RushHour
COPY . src/github.com/yasshi2525/RushHour

RUN apk update && apk add --no-cache git
RUN sed -i -e "s|conf/game.conf|src/github.com/yasshi2525/RushHour/conf/game.conf|" src/github.com/yasshi2525/RushHour/app/services/config.go && \
    sed -i -e "s|conf/secret.conf|src/github.com/yasshi2525/RushHour/conf/secret.conf|" src/github.com/yasshi2525/RushHour/app/services/secret.go

RUN go get gopkg.in/go-playground/validator.v9 && \
    go get github.com/BurntSushi/toml && \
    go get github.com/jinzhu/gorm && \
    go get github.com/go-sql-driver/mysql && \
    go get github.com/revel/revel && \
    go get github.com/revel/cmd/revel && \
    go get github.com/gomodule/oauth1/oauth && \
    go get golang.org/x/oauth2 && \
    go get google.golang.org/api/oauth2/v2

RUN mkdir -p /rushhour && \
    cd /rushhour && \
    /go/bin/revel build -m prod -a github.com/yasshi2525/RushHour

FROM alpine

ENV APP_SECRET kO0HKDOKQRLT6y9Vo0Uk69X2nxQ1p2Ln485wrYZmxiGiR7MDHa4TBxLvwLfWojcg
RUN apk update && apk --no-cache add tzdata

WORKDIR /rushhour

COPY --from=server /rushhour/ ./
COPY --from=client /data/public/ src/github.com/yasshi2525/RushHour/public/
COPY docker-entrypoint.sh .
RUN chmod u+x docker-entrypoint.sh

EXPOSE 9000

ENTRYPOINT [ "./docker-entrypoint.sh" ]