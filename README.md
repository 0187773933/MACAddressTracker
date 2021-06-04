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