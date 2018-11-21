# OPC DA in Go
Read process and automation data in Go from an OPC server for monitoring and data analysis purposes (OPC DA protocol).

## Installation

* You need an the OPC components installed ("OPC.Automation.1")
* Install go-ole and check out pull request #155 that fixed the 32-bit environment:
  1. ```go get go-ole ```
  2. Go to $ENV:GOPATH/src/github.com/go-ole/go-ole
  3. Get pull request: ```git fetch origin pull/155/head:pr155```
  4. Check out pull request: ```git checkout pr155```
* Start Gray Simulator v1.8 (OPC Simulation Server; this is optional but necessary for testing); can be obtained [here](http://www.gray-box.net/download_graysim.php).
* Compile Go projects for 32-bit with ```$ENV:GOARCH=386```

## Example 

```
package main

import (
	"fmt"
	"github.com/konimarti/opc"
)

func main() {
	client := opc.NewConnection(
		"Graybox.Simulator",
		[]string{"localhost"},
		[]string{"numeric.sin.int64", "numeric.saw.float"},
	)

	// read single tag: value, quality, timestamp
	fmt.Println(client.ReadItem("numeric.sin.int64"))

	// read all added tags
	fmt.Println(client.Read())

	client.Close()
}``` 

with the following output:

```
{-34 192 2018-11-21 20:59:10 +0000 UTC}
map[numeric.sin.int64:-34 numeric.saw.float:88.9]
```

## OPCFLUX

* Application to write OPC data directly to InfluxDB.

## OPCAPI

* Application to expose OPC tags with a restful API.


