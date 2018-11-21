# OPC DA in Go
Read process and automation data in Go from an OPC server for monitoring and data analysis purposes (OPC DA protocol).

```go get github.com/konimarti/opc```

## Prerequisites

* Install go-ole with the 32-bit fix. There is a pull request #155 in go-ole with the necessary changes:
  1. ```go get github.com/go-ole/go-ole ```
  2. Go to $GOPATH/src/github.com/go-ole/go-ole
  3. Get the pull request: ```git fetch origin pull/155/head:pr155```
  4. Check out the branch: ```git checkout pr155```
* Start Gray Simulator v1.8 (OPC Simulation Server; this is optional but necessary for testing); can be obtained [here](http://www.gray-box.net/download_graysim.php).
* Set Go architecture to 32-bit with ```$ENV:GOARCH=386```
* Get the OPC package: ```go get github.com/konimarti/opc```
* Test code with ```go test -v```

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
	defer client.Close()

	// read single tag: value, quality, timestamp
	fmt.Println(client.ReadItem("numeric.sin.int64"))

	// read all added tags
	fmt.Println(client.Read())
}
``` 

with the following output:

```
{-34 192 2018-11-21 20:59:10 +0000 UTC}
map[numeric.sin.int64:-34 numeric.saw.float:88.9]
```

## OPCFLUX

* Application to write OPC data directly to InfluxDB.

## OPCAPI

* Application to expose OPC tags with a restful API.


