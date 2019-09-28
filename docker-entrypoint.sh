#!/bin/ash

sed -i -e "s|conf/game.conf|src/github.com/yasshi2525/RushHour/conf/game.conf|" src/github.com/yasshi2525/RushHour/app/services/config.go
sed -i -e "s|conf/secret.conf|src/github.com/yasshi2525/RushHour/conf/secret.conf|" src/github.com/yasshi2525/RushHour/app/services/secret.go
sed -i -e "s/&loc=Asia%2FTokyo//" src/github.com/yasshi2525/RushHour/conf/app.conf
sed -i -e "s|^baseurl = .*$|baseurl = \"${baseurl}\"|" src/github.com/yasshi2525/RushHour/conf/secret.conf
sed -i -e "s/__SALT__/${salt}/" src/github.com/yasshi2525/RushHour/conf/secret.conf
sed -i -e "s/______KEY_______/${key}/" src/github.com/yasshi2525/RushHour/conf/secret.conf
sed -i -e "s/__STATE__/${state}/" src/github.com/yasshi2525/RushHour/conf/secret.conf
sed -i -e "s/__TWITTER_TOKEN__/${twitter_token}/" src/github.com/yasshi2525/RushHour/conf/secret.conf
sed -i -e "s/__TWITTER_SECRET__/${twitter_secret}/" src/github.com/yasshi2525/RushHour/conf/secret.conf
sed -i -e "s/__GOOGLE_CLIENT__/${google_client}/" src/github.com/yasshi2525/RushHour/conf/secret.conf
sed -i -e "s/__GOOGLE_SECRET__/${google_secret}/" src/github.com/yasshi2525/RushHour/conf/secret.conf

./run.sh
