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

ARG baseurl
ARG salt
ARG key
ARG state
ARG twitter_token
ARG twitter_secret
ARG google_client
ARG google_secret

RUN mkdir -p src/github.com/yasshi2525/RushHour
COPY . src/github.com/yasshi2525/RushHour
RUN sed -i -e "s|conf/game.conf|src/github.com/yasshi2525/RushHour/conf/game.conf|" src/github.com/yasshi2525/RushHour/app/services/config.go && \
    sed -i -e "s|conf/secret.conf|src/github.com/yasshi2525/RushHour/conf/secret.conf|" src/github.com/yasshi2525/RushHour/app/services/secret.go && \
    sed -i -e "s/&loc=Asia%2FTokyo//" src/github.com/yasshi2525/RushHour/conf/app.conf && \
    sed -i -e "s|^baseurl = .*$|baseurl = \"baseurl = ${baseurl}\"|" src/github.com/yasshi2525/RushHour/conf/secret.conf && \
    sed -i -e "s/__SALT__/\"${salt}\"/" src/github.com/yasshi2525/RushHour/conf/secret.conf && \
    sed -i -e "s/______KEY_______/\"${key}\"/" src/github.com/yasshi2525/RushHour/conf/secret.conf && \
    sed -i -e "s/__STATE__/\"${state}\"/" src/github.com/yasshi2525/RushHour/conf/secret.conf && \
    sed -i -e "s/__TWITTER_TOKEN__/\"${twitter_token}\"/" src/github.com/yasshi2525/RushHour/conf/secret.conf && \
    sed -i -e "s/__TWITTER_SECRET__/\"${twitter_secret}\"/" src/github.com/yasshi2525/RushHour/conf/secret.conf && \
    sed -i -e "s/__GOOGLE_CLIENT__/\"${google_client}\"/" src/github.com/yasshi2525/RushHour/conf/secret.conf && \
    sed -i -e "s/__GOOGLE_SECRET__/\"${google_secret}\"/" src/github.com/yasshi2525/RushHour/conf/secret.conf

RUN apk update && apk add --no-cache git

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
    /go/bin/revel build -m dev -a github.com/yasshi2525/RushHour

FROM alpine

ENV APP_SECRET kO0HKDOKQRLT6y9Vo0Uk69X2nxQ1p2Ln485wrYZmxiGiR7MDHa4TBxLvwLfWojcg
RUN apk update && apk --no-cache add tzdata

WORKDIR /rushhour

COPY --from=server /rushhour/ ./
COPY --from=client /data/public/ src/github.com/yasshi2525/RushHour/public/

EXPOSE 9000

ENTRYPOINT [ "./run.sh" ]