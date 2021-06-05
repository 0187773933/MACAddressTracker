#!/bin/bash
/usr/local/go/bin/go build -o macAddressTracker
chmod +x ./macAddressTracker
sudo cp ./macAddressTracker /usr/bin/