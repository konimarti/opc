# OPC DA in Go

[![License](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://github.com/konimarti/opc/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/konimarti/observer?status.svg)](https://godoc.org/github.com/konimarti/opc)
[![goreportcard](https://goreportcard.com/badge/github.com/konimarti/observer)](https://goreportcard.com/report/github.com/konimarti/opc)

Read and write process and automation data in Go from an OPC server for monitoring and data analysis purposes (OPC DA protocol).

```go get github.com/konimarti/opc```

## Usage
```go
client, _ := opc.NewConnection(
	"Graybox.Simulator", 		// ProgId
	[]string{"localhost"}, 		// Nodes
	[]string{"numeric.sin.float"}, 	// Tags
)
defer client.Close()
client.ReadItem("numeric.sin.float")
```

```go
browser, _ := opc.CreateBrowser(
	"Graybox.Simulator", 		// ProgId
	[]string{"localhost"}, 		// Nodes	
)
opc.PrettyPrint(browser)
```


## Installation

* ```go get github.com/konimarti/opc```

### Troubleshooting

* OPC DA Automation Wrapper 2.02 should be installed on your system (```OPCDAAuto.dll``` or ```gbda_aut.dll```); the automation wrapper is usually shipped as part of the OPC Core Components of your OPC Server.
* You can get the Graybox DA Automation Wrapper [here](http://gray-box.net/download_daawrapper.php?lang=en). Follow the [installation instruction](http://gray-box.net/daawrapper.php) for this wrapper. 
* Depending on whether your OPC server and automation wrapper are 32-bit or 64-bit, set the Go architecture correspondingly:
  - For 64-bit OPC servers and wrappers: DLL should be in ```C:\Windows\System32```, use ```$ENV:GOARCH="amd64"```
  - For 32-bit OPC servers and wrappers: DLL should be in ```C:\Windows\SysWOW64```, use ```$ENV:GOARCH="386"```
* Make sure to have correct DCOM settings on your local and remote computers: ```Dcomcnfg.exe```

### Debugging

* Add ```opc.Debug()``` before the ```opc.NewConnection``` call to print more debug-related information.

### Testing

* Start Graybox Simulator v1.8. This is a free OPC simulation server and require for testing this package. It can be downloaded [here](http://www.gray-box.net/download_graysim.php).
* If you use the Graybox Simulator, set $GOARCH environment variable to "386", i.e. enter ```$ENV:GOARCH=386``` in Powershell.
* Test code with ```go test -v```

## Example 

```go
package main

import (
	"fmt"
	"github.com/konimarti/opc"
)

func main() {
	client, _ := opc.NewConnection(
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
{91 192 2019-06-21 15:23:08 +0000 UTC}
map[numeric.sin.int64:{91 192 2019-06-21 15:23:08 +0000 UTC} numeric.saw.float:{-36.42 192 2019-06-21 15:23:08 +0000 UTC
}]
```

## Applications 

### opc-cli

* ```opc-cli``` is a command-line interface to work with OPC servers: list available OPC servers, browse OPC tags on server, and read/write OPC tags.
* Install it with ```go install github.com/konimarti/opc/cmds/opc-cli```

  - List OPC servers on a specific node: 
    ```
    $ opc-cli.exe list localhost
	Found 3 server(s) on 'localhost':
	Graybox.Simulator.1
	INAT TcpIpH1 OPC Server
	Prosys.OPC.Simulation
    ```

  - Browse OPC tags (in sub-branch):
    ```
    $ opc-cli.exe browse localhost Graybox.Simulator.1 textual
	textual
	   - textual.color
	   - textual.number
	   - textual.random
	   - textual.weekday
    ```

  - Write to OPC tag:
    ```
    $ opc-cli.exe write localhost Graybox.Simulator.1 options.sinfreq 0.01
    ```

  - Read OPC tags:
    ```
    $ opc-cli.exe read localhost Graybox.Simulator.1 options.sinfreq numeric.sin.float
	map[options.sinfreq:{0.05 192 2019-06-21 15:26:02 +0000 UTC} numeric.sin.float:{22.916641 192 2019-06-21 15:26:02 +0000 UTC}]
    ```


### OPCAPI

* Application to expose OPC tags with a JSON REST API.

  - Install the app: ```go install github.com/konimarti/opc/cmds/opcapi```

  - Create config file:
    ```
    [config]
    allow_write = false
    allow_add = true
    allow_remove = true
  
    [opc]
    server = "Graybox.Simulator"
    nodes = [ "localhost" ]
    tags = [ "numeric.sin.float", "numeric.saw.float" ]
  
    ```

  - Run app: 
    ```
    $ opcapi.exe -conf api.conf -addr ":4444"
    ```

  - Access API:
    - Get tags: 
      ```
      $ curl.exe -X GET localhost:4444/tags
      {"numeric.saw.float":-21.41,"numeric.sin.float":62.303356}
      ```
    - Add tag: 
      ```
      $ curl.exe -X POST -d '["numeric.triangle.float"]' localhost:4444/tag
      {"result": "created"}
      ```
    - Remove tag: 
      ```
      $ curl.exe -X DELETE localhost:4444/tag/numeric.triangle.float
      {"result": "removed"}
      ```

### OPCFLUX

* Application to write OPC data directly to InfluxDB.

  - Install the app: ```go install github.com/konimarti/opc/cmds/opcflux```

  - Create InfluxDB database "test"

  - Create config file:
    Put OPC tags in []. This is required for the expression evaluation. Any calculation can be performed that can evaluated.
    ```
    ---
    server: "Graybox.Simulator"
    nodes: ["localhost", "127.0.0.1"]
    monitoring: ""
    influx:
     addr: "http://localhost:8086"
     database: test
     precision: s
    measurements: 
     numeric:
       - tags: {type: sin}
         fields: {float: "[numeric.sin.float]", int: "[numeric.sin.int32]"}
       - tags: {type: saw}
         fields: {float: "[numeric.saw.float]", int: "[numeric.saw.int32]"}
       - tags: {type: calculation}
         fields: {float: "[numeric.triangle.float] / [numeric.triangle.int32]"}
     textual:
       - tags: {type: color}
         fields: {text: "[textual.color]", brown: "[textual.color] == 'Brown'"}        
       - tags: {type: weekday}
         fields: {text: "[textual.weekday]"}        
    ```

  - Run app: ```opcflux.exe -conf influx.yml -rate 1s```

## Credits

This software package has been developed for and is in production at [Kalkfabrik Netstal](http://www.kfn.ch/en).
