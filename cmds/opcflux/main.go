//set OPC_SERVER=Graybox.Simulator && set OPC_NODES=127.0.0.1,localhost && go run main.go -conf influx.yml -rate 100ms

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/knetic/govaluate"
	"github.com/konimarti/opc"
	yaml "gopkg.in/yaml.v2"
)

var (
	config = flag.String("conf", "influx.yml", "yaml config file for tag descriptions")
	rr     = flag.String("rate", "10s", "refresh rate as duration, e.g. 100ms, 5s, 10s, 2m")
)

// M stores an InfluxDB measurement
type M struct {
	Tags   map[string]string
	Fields map[string]string
}

// Database represents an InfluxDB database connection
type Database struct {
	Addr      string
	Username  string
	Password  string
	Database  string
	Precision string
}

// Conf contains config data
type Conf struct {
	Server       string
	Nodes        []string
	Monitoring   string
	Influx       Database
	Measurements map[string][]M
}

func main() {
	flag.Parse()

	//set refresh rate
	refreshRate, err := time.ParseDuration(*rr)
	if err != nil {
		log.Fatalf("error setting refresh rate")
	}
	fmt.Println("refresh rate: ", refreshRate)

	// read config
	conf := getConfig(*config)

	//app monitoring
	if conf.Monitoring != "" {
		opc.StartMonitoring(conf.Monitoring)
	}

	// extract tags
	tags := []string{}
	exprMap := make(map[string]*govaluate.EvaluableExpression)
	for _, group := range conf.Measurements {
		for _, m := range group {
			for _, f := range m.Fields {
				expr, err := govaluate.NewEvaluableExpression(f)
				if err != nil {
					fmt.Println("Could not parse", f)
					panic(err)
				}
				exprMap[f] = expr
				tags = append(tags, expr.Vars()...)
			}
		}
	}

	//setup influxdb client
	//TODO: get username and password for influx from environment variables
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: conf.Influx.Addr,
		//Username: conf.Influx.Username,
		//Password: conf.Influx.Password,
	})
	if err != nil {
		fmt.Println(err)
		panic("Error creating InfluxDB Client")
	}
	defer c.Close()
	fmt.Println("Writing to", conf.Influx.Database, "at", conf.Influx.Addr)

	if conf.Server == "" {
		conf.Server = strings.Trim(os.Getenv("OPC_SERVER"), " ")
	}
	if len(conf.Nodes) == 0 {
		conf.Nodes = strings.Split(os.Getenv("OPC_NODES"), ",")
	}

	conn := opc.NewConnection(
		conf.Server,
		conf.Nodes,
		tags,
	)
	if conn == nil {
		panic("Could not create OPC connection.")
	}

	timeC := make(chan time.Time, 10)

	// start go routine
	go writeState(timeC, c, conn, conf, exprMap)

	// start ticker
	ticker := time.NewTicker(refreshRate)
	for tick := range ticker.C {
		timeC <- tick
	}
}

// getConfig parses configuration file
func getConfig(config string) *Conf {
	log.Println("config file: ", config)

	content, err := ioutil.ReadFile(config)
	if err != nil {
		log.Fatalf("error reading config file %s", config)
	}

	conf := Conf{}
	err = yaml.Unmarshal([]byte(content), &conf)
	if err != nil {
		log.Fatalf("error yaml unmarshalling: %v", err)
	}

	// fmt.Printf("--- conf:\n%v\n\n", conf)

	return &conf
}

// writeState collects data and writes it to the influx database
func writeState(timeC chan time.Time, c client.Client, conn opc.Connection, conf *Conf, exprMap map[string]*govaluate.EvaluableExpression) {

	batchconfig := client.BatchPointsConfig{
		Database:  conf.Influx.Database,
		Precision: conf.Influx.Precision, // "s"
	}

	for t := range timeC {

		// read data
		data := conn.ReadValues()

		// create a new point batch
		bp, err := client.NewBatchPoints(batchconfig)
		if err != nil {
			fmt.Println(err)
			return
		}

		// define measurement and create data points
		for measurement, group := range conf.Measurements {
			//t := time.Now().Local()
			//data := conn.Read()
			for _, m := range group {
				tagMap := m.Tags
				fieldMap := make(map[string]interface{})

				for fieldKey, f := range m.Fields {
					ist, err := exprMap[f].Evaluate(data)
					if err != nil {
						fmt.Println(err)
						continue
					}
					fieldMap[fieldKey] = ist
				}

				// create influx data points
				pt, err := client.NewPoint(measurement, tagMap, fieldMap, t)
				if err != nil {
					fmt.Println("Error: ", err.Error())
				}

				// add data point to batch
				bp.AddPoint(pt)
			}
		}

		// write to database
		if err := c.Write(bp); err != nil {
			fmt.Println(err)
		}
	}
}
