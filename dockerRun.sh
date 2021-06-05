#!/bin/bash
name="mac-address-tracker"
sudo docker rm $name -f || echo ""
sudo docker build -t $name .
id=$(sudo docker run -dit --restart='always' \
--name $name \
--net=host \
--privileged \
-v ${PWD}/config.json:"/home/morphs/.config/personal/mac_address_tracker.json" \
$name)
sudo docker logs -f $id

# sudo docker rm $name -f || echo ""
# sudo docker build -t $name .
# sudo docker run -it --net=host \
# -e MAC_LOCATION_NAME="1234 Merily Way" \
# -e MAC_CRON_STRING="@every 5m" \
# -e MAC_SERVER_PORT="1234" \
# -e MAC_SAVED_RECORD_TOTAL="100" \
# -e MAC_NETWORK_HARDWARE_INTERFACE_NAME="en0" \
# -e MAC_REDIS_HOST="11.22.33.44" \
# -e MAC_REDIS_PORT="6379" \
# -e MAC_REDIS_DB="0" \
# -e MAC_REDIS_PASSWORD="asdf" \
# -e MAC_REDIS_PREFIX="MACS." \
# $name

# id=$(sudo docker run -dit --restart='always' \
# --name public-homebridge \
# --net=host \
# public-homebridge)
# sudo docker logs -f $id
