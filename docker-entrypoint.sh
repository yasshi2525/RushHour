#!/bin/ash

sed -i -e "s/__ADMIN_USERNAME__/${admin_username}/" /tmp/conf/secret.conf
sed -i -e "s/__ADMIN_PASSWORD__/${admin_password}/" /tmp/conf/secret.conf
sed -i -e "s|^baseurl = .*$|baseurl = \"${baseurl}\"|" /tmp/conf/secret.conf
sed -i -e "s/__SALT__/${salt}/" /tmp/conf/secret.conf
sed -i -e "s/______KEY_______/${key}/" /tmp/conf/secret.conf
sed -i -e "s/__STATE__/${state}/" /tmp/conf/secret.conf
sed -i -e "s/__TWITTER_TOKEN__/${twitter_token}/" /tmp/conf/secret.conf
sed -i -e "s/__TWITTER_SECRET__/${twitter_secret}/" /tmp/conf/secret.conf
sed -i -e "s/__GOOGLE_CLIENT__/${google_client}/" /tmp/conf/secret.conf
sed -i -e "s/__GOOGLE_SECRET__/${google_secret}/" /tmp/conf/secret.conf
sed -i -e "s/__GITHUB_CLIENT__/${github_client}/" /tmp/conf/secret.conf
sed -i -e "s/__GITHUB_SECRET__/${github_secret}/" /tmp/conf/secret.conf

TARGET="src/github.com/yasshi2525/RushHour/conf"

old_pwd="$(pwd)"

cd "/tmp/conf"
find -mindepth 1 | while read line; do
    filename="$(basename "$line")"
    if [ ! -e "$TARGET/$filename" ]; then
        echo "copy defalut config $filename"
        cp "$line" "$TARGET"
    fi
    done
cd "$old_pwd"

./run.sh
