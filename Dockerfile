FROM node:12 AS client

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

RUN mkdir -p src/github.com/yasshi2525/RushHour
COPY . src/github.com/yasshi2525/RushHour

RUN apk update && apk add --no-cache git && \
    mkdir -p /rushhour/config && \
    cd src/github.com/yasshi2525/RushHour/app && \
    go mod download && \
    go build -o /rushhour/RushHour && \
    cp config/*.conf /rushhour/config

FROM alpine

ENV admin_username "admin"
ENV admin_password "password"
ENV baseurl "https://localhost:9000/"
ENV salt ""
ENV key "1234567890123456"
ENV state ""
ENV cookie kO0HKDOKQRLT6y9Vo0Uk69X2nxQ1p2Ln485wrYZmxiGiR7MDHa4TBxLvwLfWojcg
ENV db_spec "rushhourgo:rushhourgo@tcp(localhost:3306)/rushhourgo?parseTime=true&loc=Asia%2FTokyo"
ENV twitter_token ""
ENV twitter_secret ""
ENV google_client ""
ENV google_secret ""

RUN apk update && apk --no-cache add tzdata && \
    addgroup rushhour && adduser rushhour --disabled-password -G rushhour

WORKDIR /rushhour

COPY --from=server --chown=rushhour:rushhour /rushhour/ ./
COPY --from=client --chown=rushhour:rushhour /data/public/ src/github.com/yasshi2525/RushHour/public/
COPY --chown=rushhour:rushhour docker-entrypoint.sh .

RUN chmod u+x docker-entrypoint.sh

EXPOSE 9000

VOLUME [ "/rushhour/src/github.com/yasshi2525/RushHour/log" ]

USER rushhour

ENV GIN_MODE "release"

ENTRYPOINT [ "./docker-entrypoint.sh" ]