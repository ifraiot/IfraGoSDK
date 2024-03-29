package ifrasdk

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mainflux/senml"
)

type Measurement struct {
	Name  string
	Value float64
}

type ifra struct {
	Topic        string
	Username     string
	Password     string
	Measurements []Measurement
	MQTTClient   mqtt.Client
}

type Ifra interface {
	AddMeasurement(name string, value float64)
	Send() error
	ToJson() string
	Disconnect()
}

const IFRA_MQTT_BROKER_HOST = "nebular.mqtt.ifra.io"
const IFRA_MQTT_BROKER_PORT = 1883

func NewIFRA(topic, username, password string) (Ifra,error) {
	fmt.Println("start connect: ", topic, username, password)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", IFRA_MQTT_BROKER_HOST, IFRA_MQTT_BROKER_PORT))
	opts.SetClientID("go_mqtt_client")
	opts.SetUsername(username)
	opts.SetPassword(password)
	// opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = MQTTConnectHandler
	// opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(5 * time.Second) {
	}
	if err := token.Error(); err != nil {
		fmt.Println(err)
		return nil,err
	}

	return &ifra {
		Topic:      topic,
		Username:   username,
		Password:   password,
		MQTTClient: client,
		
	},nil
}

var MQTTConnectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("IFRA MQTT: Connected")
}

func (i *ifra) AddMeasurement(name string, value float64) {
	i.Measurements = append(i.Measurements, Measurement{
		Name:  name,
		Value: value,
	})
}

func (i *ifra) Disconnect() {
	i.MQTTClient.Disconnect(250)
}

func (i *ifra) Send() error{
	//Convert data to SenML format
	// var record senml.Pack

	var pack = senml.Pack{}
	for _, ms := range i.Measurements {

		var value = ms.Value
		pack.Records = append(pack.Records, senml.Record{
			Name:  ms.Name,
			Value: &value,
		})
	}

	enc, err := senml.Encode(pack, senml.JSON)
	if err != nil {
		fmt.Println(err)
		return err
	}

	token := i.MQTTClient.Publish(i.Topic, 0, false, string(enc))
	for !token.WaitTimeout(5 * time.Second) {
	}
	if err := token.Error(); err != nil {
		fmt.Println(err)
		return err
	}

	//fmt.Println(string(enc))

	//Clear measurement data
	i.Measurements = []Measurement{}
	return nil
}

func (i *ifra) ToJson() string {
	return ""
}
