package main

import (
	"encoding/json"
	"flag"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/konimarti/opc"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"time"
)

var (
	config = flag.String("conf", "./cmds/opcmqtt/mqtt.yml", "yaml config file for tag transport descriptions")
)

type MqttBroker struct {
	Addr     string
	Username string
	Password string
	Topic    string
}

type Conf struct {
	Server      string     `yaml:"server"`
	Nodes       []string   `yaml:"nodes"`
	RefreshRate string     `yaml:"refreshRate"`
	Mqtt        MqttBroker `yaml:"mqtt"`
	Tags        []string   `yaml:"tags"`
}

// getConfig parses configuration file
func getConfig(config string) *Conf {
	log.Println("config file: ", config)

	content, err := os.ReadFile(config)
	if err != nil {
		log.Fatalf("error reading config file %s", config)
	}

	conf := Conf{}
	err = yaml.Unmarshal(content, &conf)
	if err != nil {
		log.Fatalf("error yaml unmarshalling: %v", err)
	}

	return &conf
}

func main() {
	flag.Parse()

	conf := getConfig(*config)

	// connect opc server
	opc.Debug()
	connOpc, err := opc.NewConnection(conf.Server, conf.Nodes, conf.Tags)
	if err != nil {
		log.Fatalf("opc connection error: %v", err)
	}

	// connect mqtt broker
	opts := mqtt.NewClientOptions().AddBroker(conf.Mqtt.Addr)
	opts.SetClientID("opc tags")

	connMqtt := mqtt.NewClient(opts)
	if token := connMqtt.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("mqtt connect error: %v", token.Error())
	}

	//set refresh rate
	refreshRate, err := time.ParseDuration(conf.RefreshRate)
	if err != nil {
		log.Fatalf("error setting refresh rate")
	}

	timeC := make(chan time.Time, 10)

	// Do ticker task
	go transport(timeC, connOpc, connMqtt, conf)

	// start ticker
	ticker := time.NewTicker(refreshRate)
	for tick := range ticker.C {
		timeC <- tick
	}
}

func adapter(data map[string]opc.Item) map[string]interface{} {
	output := make(map[string]interface{})
	for k, item := range data {
		if item.Good() {
			output[k] = item.Value
		} else {
			log.Printf("k=%v, item=%v not good tag\n", k, item)
		}
	}
	return output
}

func transport(timeC chan time.Time, connOpc opc.Connection, connMqtt mqtt.Client, conf *Conf) {
	for range timeC {
		data := adapter(connOpc.Read())
		if len(data) == 0 {
			log.Printf("empty data")
			continue
		}

		b, err := json.Marshal(data)
		if err != nil {
			log.Printf("error marshalling: %v, data: %v\n", err, data)
			continue
		}

		if token := connMqtt.Publish(conf.Mqtt.Topic, 0, false, b); token.Wait() && token.Error() != nil {
			log.Printf("mqtt publish error: %v", token.Error())
			continue
		}
		log.Printf("send data=%+v success\n", data)
	}
}
