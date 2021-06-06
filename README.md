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
#!/bin/bash
name="mac-address-tracker"
sudo docker rm $name -f || echo ""
sudo docker build -t $name .
id=$(sudo docker run -dit --restart='always' \
--name $name \
--net=host \
-v ${PWD}/config.json:"/home/morphs/.config/personal/mac_address_tracker.json" \
$name)
sudo docker logs -f $id
```

### Can't Figure out Environment Variables from docker run --> go binary --> os.Getenv()
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

```bash
redis-cli -n 0 -p 6379 --no-auth-warning -a asdf
```

```bash
set MACS.SEEN.a1:b2:c3:d4:a2:b2  "{\"device_name\":\"asdf cellphone\",\"current_time_string\":null,\"records\":null,\"transitions\":null}"
```

```bash
echo "nameserver 8.8.8.8" > sudo tee /etc/resolv.conf && \
sudo apt-get install --reinstall resolvconf network-manager libnss-resolve
```


```bash
sudo nano /etc/resolv.conf
nameserver 192.168.1.1
nameserver 8.8.8.8
nameserver 8.8.4.4
```

```bash
sudo nano /etc/hosts
127.0.0.1 localhost
127.0.1.1 mediabox
```


```bash
pm2 start --interpreter none --name MacAddressTracker /home/morphs/DOCKER_IMAGES/MACAddressTracker/macAddressTracker -- /home/morphs/DOCKER_IMAGES/MACAddressTracker/config.json
```