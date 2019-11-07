FROM node:12 AS client

WORKDIR /data
COPY . .

RUN apt-get update && apt-get install -y git && \
    npm ci && \
    npm run build && \
    git clone https://github.com/yasshi2525/RushHourResource.git && \
    mkdir -p ./assets/bundle/spritesheet && \
    cp -r RushHourResource/dist/* ./assets/bundle/spritesheet/

FROM golang:alpine as server

WORKDIR /work

COPY . .

RUN apk update && apk add --no-cache git && \
    go mod download && \
    mkdir -p ./dist && \
    go build -o ./dist/RushHour && \
    cp -R config ./dist && \
    cp -R assets ./dist && \
    cp -R templates ./dist

FROM alpine

ENV persist "false"
ENV admin_username "admin"
ENV admin_password "password"
ENV baseurl "http://localhost:8080/"
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

COPY --from=server --chown=rushhour:rushhour /work/dist/ ./
COPY --from=client --chown=rushhour:rushhour /data/assets/bundle /rushhour/assets
COPY --chown=rushhour:rushhour docker-entrypoint.sh .

RUN chmod u+x docker-entrypoint.sh

EXPOSE 8080

VOLUME [ "/rushhour/logs" ]

USER rushhour

ENV GIN_MODE "release"

ENTRYPOINT [ "./docker-entrypoint.sh" ]