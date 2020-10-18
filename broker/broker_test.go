package broker

import (
	"fmt"
	"github.com/fhmq/hmq/metrics"
	"github.com/prometheus/common/expfmt"
	"net/http"
	"os"
	"runtime"
	"sync"
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	dto "github.com/prometheus/client_model/go"
)

const (
	defaultPacketPayload = "mymessage"
	defaultTopic         = "mytopic/test"
)

type BrokerTests struct {
	client mqtt.Client
}

var (
	bt BrokerTests
	b  Broker
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	startBroker()

	// Wait for broker to start-up before trying to create a connection
	for ok := true; ok; ok = b.started {
		time.Sleep(time.Second)
	}

	bt = BrokerTests{
		client: nil,
	}

	bt.connect()
}

func teardown() {
	bt.client.Disconnect(0)
}

func startBroker() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	config := DefaultConfig
	config.HTTPPort = "8080"

	b, err := NewBroker(config)
	if err != nil {
		log.Fatal(fmt.Sprintf("New Broker error: %e", err))
	}
	b.Start()
}

func (bt *BrokerTests) connect() {
	opts := mqtt.NewClientOptions().AddBroker("tcp://127.0.0.1:1883").SetClientID("broker-test")
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(fmt.Sprintf("%e", token.Error()))
	}

	bt.client = c
}

func (bt *BrokerTests) publishMessage(topic string, message string, wg *sync.WaitGroup) error {
	if token := bt.client.Publish(topic, 0, false, message); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	wg.Add(1)
	return nil
}

func (bt *BrokerTests) listenMessages(topic string, message string, wg *sync.WaitGroup) error {
	if token := bt.client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		if string(msg.Payload()) == message {
			log.Info(fmt.Sprintf("received expected message: %s", msg.Payload()))
			wg.Done()
		} else {
			log.Info(fmt.Sprintf("received message not matching criteria: %s", msg.Payload()))
		}
	}); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func fetchMetrics() (map[string]*dto.MetricFamily, error) {
	resp, err := http.Get("http://127.0.0.1:8080/metrics")
	if err != nil {
		log.Fatal("failed to fetch metrics")
	}
	defer resp.Body.Close()

	var parser expfmt.TextParser
	return parser.TextToMetricFamilies(resp.Body)
}

func TestBrokerPubSub(t *testing.T) {
	var wg sync.WaitGroup
	if err := bt.listenMessages(defaultTopic, defaultPacketPayload, &wg); err != nil {
		t.Fatal(err)
	}

	for i := 0; i <= 10; i++ {
		if err := bt.publishMessage(defaultTopic, defaultPacketPayload, &wg); err != nil {
			t.Fatal(err)
		}
	}
	wg.Wait() // Wait for all messages to be processed before proceeding

	metricFamilies, err := fetchMetrics()
	if err != nil {
		t.Fatal(fmt.Sprintf("failed parsing prometheus metrics: %e", err))
	}

	metric1 := metricFamilies[fmt.Sprintf("gin_gin_%s", metrics.MetricNumberOfMessages)]
	metric2 := metricFamilies[fmt.Sprintf("gin_gin_%s", metrics.MetricNumberOfClients)]
	if metric1 == nil || metric2 == nil {
		t.Fatal("metric should not be nil")
	}

	if *metric1.Metric[0].Gauge.Value < 0 {
		t.Fatal("metric 0, but should be > 0")
	}

	if *metric2.Metric[0].Gauge.Value != 1 {
		t.Fatal("metric 0, but should be 1")
	}
}
