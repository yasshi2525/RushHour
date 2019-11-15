FROM node:12 AS client

ENV baseurl "http://localhost:8080"

WORKDIR /data
COPY . .

ENV RES_VRS "0.1.0"

RUN npm ci && \
    npm run build && \
    curl -LsS https://github.com/yasshi2525/RushHourResource/archive/v${RES_VRS}.tar.gz | tar zx && \
    mkdir -p ./assets/bundle/spritesheet && \
    cp -r RushHourResource-${RES_VRS}/dist/* ./assets/bundle/spritesheet/

FROM golang:alpine as server

WORKDIR /work

COPY . .

RUN apk update && apk add --no-cache git && \
    go mod download && \
    mkdir -p ./dist/config && \
    go build -o ./dist/RushHour && \
    cp -R config/*.conf ./dist/config && \
    cp -R assets ./dist && \
    cp -R templates ./dist

FROM alpine

ENV persist "false"
ENV admin_username "admin"
ENV admin_password "password"
ENV baseurl "http://localhost:8080"
ENV salt ""
ENV key "1234567890123456"
ENV state ""
ENV cookie kO0HKDOKQRLT6y9Vo0Uk69X2nxQ1p2Ln485wrYZmxiGiR7MDHa4TBxLvwLfWojcg
ENV db_spec "rushhourgo:rushhourgo@tcp(localhost:3306)/rushhourgo?parseTime=true&loc=Asia%2FTokyo"
ENV twitter_token ""
ENV twitter_secret ""
ENV google_client ""
ENV google_secret ""
ENV github_client ""
ENV github_secret ""

RUN apk update && apk --no-cache add tzdata && \
    addgroup rushhour && adduser rushhour --disabled-password -G rushhour

WORKDIR /rushhour

COPY --from=server --chown=rushhour:rushhour /work/dist/ ./
COPY --from=client --chown=rushhour:rushhour /data/assets/bundle/ /rushhour/assets/bundle/
COPY --chown=rushhour:rushhour docker-entrypoint.sh .

RUN chmod u+x docker-entrypoint.sh

EXPOSE 8080

VOLUME [ "/rushhour/logs" ]

USER rushhour

ENV GIN_MODE "release"

ENTRYPOINT [ "./docker-entrypoint.sh" ]