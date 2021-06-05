# Requirements

- arp
- nmap


# Misc

- https://github.com/0187773933/RedisManagerUtils/blob/master/manager/manager.go
- https://stackoverflow.com/questions/8978670/what-do-windows-interface-names-look-like
  - network_map = scanner.ScanLocalNetwork( "ethernet_0" )



```bash
redis-cli -n 0 -p 6379 --no-auth-warning -a asdf KEYS "MACS.SEEN.*" | xargs redis-cli -n 0 -p 6379 --no-auth-warning -a asdf DEL
```

```bash
sudo crontab -e
```
```bash
*/5 * * * * /bin/bash -l -c 'su morphs -c "/usr/local/bin/macAddressTracker"' >/dev/null 2>&1
```
```bash
tail -f /var/log/syslog | grep 'CRON'
```

```bash
name="mac-address-tracker"
id=$(sudo docker run -dit --restart='always' \
--name public-homebridge \
--net=host \
-e MAC_LOCATION_NAME="1234 Merily Way" \
-e MAC_CRON_STRING="@every 5m" \
-e MAC_SERVER_PORT="1234" \
-e MAC_SAVED_RECORD_TOTAL="100" \
-e MAC_NETWORK_HARDWARE_INTERFACE_NAME="en0" \
-e MAC_REDIS_HOST="11.22.33.44" \
-e MAC_REDIS_PORT="6379" \
-e MAC_REDIS_DB="0" \
-e MAC_REDIS_PASSWORD="asdf" \
-e MAC_REDIS_PREFIX="MACS." \
public-homebridge)
sudo docker logs -f $id
```