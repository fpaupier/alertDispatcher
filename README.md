[![Go Report Card](https://goreportcard.com/badge/github.com/fpaupier/alertDispatcher)](https://goreportcard.com/report/github.com/fpaupier/alertDispatcher)

# Alert Dispatcher

This code runs on your IoT device and regularly polls the local storage to see if there are events
to publish to the Kafka topics.

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


## Installing Kafka
The [Confluent Kafka Go Client](https://github.com/confluentinc/confluent-kafka-go) requires a manual installation, and a few tricks
to run on a Raspberry Pi. First, install the package

````shell script
go get github.com/confluentinc/confluent-kafka-go/kafka
```` 

It takes a few seconds, and it's likely you get the following error message:
```
# github.com/confluentinc/confluent-kafka-go/kafka
/usr/bin/ld: ../go/src/github.com/confluentinc/confluent-kafka-go/kafka/librdkafka/librdkafka_glibc_linux.a:
error adding symbols: file format not recognized collect2: error: ld
returned 1 exit status
```
If this happens, it's because there is no compatible `librdkafka` build on for your [ARM](https://en.wikipedia.org/wiki/ARM_architecture) based Raspberry PI. No worries! 
 We can fix this by building `librdkafka` from source. On your Pi execute:
 ````shell script
git clone https://github.com/edenhill/librdkafka.git
cd librdkafka
./configure
make
sudo make install
````
Now, `cd` to the `alertDispatcher/` folder and run the go application with the `-tags dynamic` tag: 
````shell script
go run -tags dynamic .
````
You might get an additional error:
````text
# github.com/confluentinc/confluent-kafka-go/kafka
../go/src/github.com/confluentinc/confluent-kafka-go/kafka/message.go:115:16:
type [1073741824]_Ctype_struct_tmphdr_s larger than address
space ../go/src/github.com/confluentinc/confluent-kafka-go/kafka/message.go:115:16:
type [1073741824]_Ctype_struct_tmphdr_s too
large ../go/src/github.com/confluentinc/confluent-kafka-go/kafka/message.go:190:30:
constant 2147483648 overflows
int ../go/src/github.com/confluentinc/confluent-kafka-go/kafka/message.go:190:30:
array bound is too large
````
Those errors are due to the fact that the code in ``confluentinc/confluent-kafka-go/kafka/message.go`` uses the `unsafe` Go library to perform type conversion by directly manipulating values at a given 
memory address instead of using Go type casting features. This makes the code non portable and risks of breaking when running on different platforms.

Again, we can easily fix this by modifying two lines in the `message.go` file. The error indicates that the lines 115 and 190 are responsible for the exception raised.  
Line 115 reads: 
```go
tmphdr := (*[1 << 30]C.tmphdr_t)(unsafe.Pointer(fcMsg.tmphdrs))[n]
```
Let's understand why it fails. This line is part of a function named: ``newMessageFromFcMsg``.  It creates a message 
from a C typed message, and a conversion occurs for the header. It initializes an address by right shifting the
bits 30 times (`[1 << 30]`) but this is not a valid memory address on our Pi leading to the exception. We also have an
information that this value is too large in the logs, trying with a lower value, like 10 works on a Raspberry Pi 4 B.

Thus replace the line 14 by:
````go
tmphdr := (*[1 << 10]C.tmphdr_t)(unsafe.Pointer(fcMsg.tmphdrs))[n]
````

The line 190 reads:
 ```go
keyp = unsafe.Pointer(&((*[1 << 31]byte)(payload)[valueLen]))
```
Line 190 belong to the function ``messageToC`` whose signature indicates it sets the attribute of a C message based on the value from a Go message.
Again, the logs tell us we overflow when copying the payload. Let's decrease the value of the right bit shift until we don't have overflow anymore.
With 30 instead of 31 it works. Thus, replace line 190 with:
 ```go
keyp = unsafe.Pointer(&((*[1 << 30]byte)(payload)[valueLen]))
```

Now, let's build again our app with the dynamic tag:
_To execute within the `alertDispatcher/` folder._
```shell script
go build -tags dynamic .
```
And _Voilà_, it should now work. Run the app with the command: `./alertDispatcher`


_Note: If experienced C programmer or Go programmers familiar with the `unsafe` package could explain what happens in the lines 115 and 190 in more
details; I'd be happy to have your insights._ 