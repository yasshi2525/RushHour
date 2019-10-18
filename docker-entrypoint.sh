#!/bin/ash

sed -i -e "s/&loc=Asia%2FTokyo//" src/github.com/yasshi2525/RushHour/conf/app.conf
sed -i -e "s/__ADMIN_USERNAME__/${admin_username}/" src/github.com/yasshi2525/RushHour/conf/secret.conf
sed -i -e "s/__ADMIN_PASSWORD__/${admin_password}/" src/github.com/yasshi2525/RushHour/conf/secret.conf
sed -i -e "s|^baseurl = .*$|baseurl = \"${baseurl}\"|" src/github.com/yasshi2525/RushHour/conf/secret.conf
sed -i -e "s/__SALT__/${salt}/" src/github.com/yasshi2525/RushHour/conf/secret.conf
sed -i -e "s/______KEY_______/${key}/" src/github.com/yasshi2525/RushHour/conf/secret.conf
sed -i -e "s/__STATE__/${state}/" src/github.com/yasshi2525/RushHour/conf/secret.conf
sed -i -e "s/__TWITTER_TOKEN__/${twitter_token}/" src/github.com/yasshi2525/RushHour/conf/secret.conf
sed -i -e "s/__TWITTER_SECRET__/${twitter_secret}/" src/github.com/yasshi2525/RushHour/conf/secret.conf
sed -i -e "s/__GOOGLE_CLIENT__/${google_client}/" src/github.com/yasshi2525/RushHour/conf/secret.conf
sed -i -e "s/__GOOGLE_SECRET__/${google_secret}/" src/github.com/yasshi2525/RushHour/conf/secret.conf
sed -i -e "s/__GITHUB_CLIENT__/${github_client}/" src/github.com/yasshi2525/RushHour/conf/secret.conf
sed -i -e "s/__GITHUB_SECRET__/${github_secret}/" src/github.com/yasshi2525/RushHour/conf/secret.conf

./run.sh
