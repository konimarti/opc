# OPC DA in Go
Read process and automation data in Go from an OPC server for monitoring and data analysis purposes (OPC DA protocol).

```go get github.com/konimarti/opc```

## Installation

* ```go get github.com/konimarti/opc```

## Prerequisites: OPC DA Automation Wrapper

* OPC DA Automation Wrapper 2.02 should be installed on your system (```OPCDAAuto.dll``` or ```gbda_aut.dll```); the automation wrapper is usually shipped as part of the OPC Core Components of your OPC Server.
* You can get the Graybox DA Automation Wrapper [here](http://gray-box.net/download_daawrapper.php?lang=en). Follow the [installation instruction](http://gray-box.net/daawrapper.php) for this wrapper. 
* Depending on whether your OPC server and automation wrapper are 32-bit or 64-bit, set the Go architecture correspondingly:
  - For 64-bit OPC servers and wrappers: DLL should be in ```C:\Windows\System32```, use ```$ENV:GOARCH="amd64"```
  - For 32-bit OPC servers and wrappers: DLL should be in ```C:\Windows\SysWOW64```, use ```$ENV:GOARCH="386"```

## Testing

* Start Graybox Simulator v1.8. This is a free OPC simulation server and require for testing this package. It can be downloaded [here](http://www.gray-box.net/download_graysim.php).
* If you use the Graybox Simulator, set $GOARCH environment variable to "386", i.e. enter ```$ENV:GOARCH=386``` in Powershell.
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
		"Graybox.Simulator", // ProgId
		[]string{"localhost"}, //  OPC servers nodes
		[]string{"numeric.sin.int64", "numeric.saw.float"}, // slice of OPC tags
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


