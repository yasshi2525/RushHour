#!/bin/ash

BASEDIR="/rushhour/src/github.com/yasshi2525/RushHour/conf"

sed -i -e "s/__ADMIN_USERNAME__/${admin_username}/" ${BASEDIR}/secret.conf
sed -i -e "s/__ADMIN_PASSWORD__/${admin_password}/" ${BASEDIR}/secret.conf
sed -i -e "s|^baseurl = .*$|baseurl = \"${baseurl}\"|" ${BASEDIR}/secret.conf
sed -i -e "s/__SALT__/${salt}/" ${BASEDIR}/secret.conf
sed -i -e "s/______KEY_______/${key}/" ${BASEDIR}/secret.conf
sed -i -e "s/__STATE__/${state}/" ${BASEDIR}/secret.conf
sed -i -e "s/__TWITTER_TOKEN__/${twitter_token}/" ${BASEDIR}/secret.conf
sed -i -e "s/__TWITTER_SECRET__/${twitter_secret}/" ${BASEDIR}/secret.conf
sed -i -e "s/__GOOGLE_CLIENT__/${google_client}/" ${BASEDIR}/secret.conf
sed -i -e "s/__GOOGLE_SECRET__/${google_secret}/" ${BASEDIR}/secret.conf
sed -i -e "s/__GITHUB_CLIENT__/${github_client}/" ${BASEDIR}/secret.conf
sed -i -e "s/__GITHUB_SECRET__/${github_secret}/" ${BASEDIR}/secret.conf

./run.sh
