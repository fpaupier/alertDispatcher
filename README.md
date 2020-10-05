# Alert Dispatcher

This code runs on your IoT device and regularly polls the local storage to see if there are events
to publish to the Kafka topics.

# Install the project

First, make sure your Pi has SQLite3 installed.
```shell script
sudo apt-get install sqlite3
```

## Packages
Run the following command to install the dependencies:
````shell script
go get github.com/golang/protobuf/proto
go get github.com/mattn/go-sqlite3
```` 

- We use [protocol buffers](https://developers.google.com/protocol-buffers) to have a consistent schema definition for messages.
- The [`go-sqlite3`](https://github.com/mattn/go-sqlite3) package offers a driver for interacting with a SQLite3 database from our Go code.  

## Credentials
Rename the `credentials-template.go` file into `credentials.go` and make sure your version control system ignores it.
Then, add your API key and secret into the file.

### Sending credentials to your Pi
On Mac and Linux, you can use the `scp` (secure copy) command to send files to your Raspberry Pi.
_Example: To send the `credentials.go` file from your development machine to your Pi:_
```shell script
scp path/to/credentials.go pi@YOUR.PI.IP.ADDRESS:/path/to/credentials.go
```   
Replace ``YOUR.PI.IP.ADDRESS`` by the IP address of your Pi. 
Run `ifconfig -a` to list your device's network information. If you're on a WiFi network, look for the `wlan` section and the 
`inet` line.

**Note: Make sur ssh is enabled on your Raspberry Pi.**
Find out how to enable SSH on your Raspberry Pi by checking the [official documentation](https://www.raspberrypi.org/documentation/remote-access/ssh/).

