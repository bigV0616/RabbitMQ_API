//Dec: For RabbitMQ support send and receive func

package lib

import (
	"bytes"
	"conf"
	// "fmt"
	"github.com/streadway/amqp"
	"log"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func getUri() string { // Join String of amqp
	argStr := bytes.Buffer{}
	argStr.WriteString("amqp://")
	argStr.WriteString(conf.Mq_username)
	argStr.WriteString(":")
	argStr.WriteString(conf.Mq_password)
	argStr.WriteString("@")
	argStr.WriteString(conf.Mq_ip)
	argStr.WriteString(":")
	argStr.WriteString(conf.Mq_port)
	argStr.WriteString("/")
	argStr.WriteString(conf.Mq_vhost)
	return argStr.String()
}

type mqApi struct {
	exName     string           // Exchange Name
	qName      string           // Queue Name
	exType     string           // Exchange Type
	routingKey string           // Routing Key
	msg        string           // send a message
	conn       *amqp.Connection // connection of connect to RabbitMQ
	ch         *amqp.Channel    // channel
	q          amqp.Queue       // queue
}

func setConnection(exName, qName, exType, routingKey, msg string) *mqApi {
	conn, err := amqp.Dial(getUri())
	FailOnError(err, "Failed to connect to RabbitMQ")
	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	return &mqApi{conn: conn, ch: ch, exName: exName, qName: qName, exType: exType, routingKey: routingKey, msg: msg}
}

// Declare a exchange
// Default the exchange of durable is true
func (m mqApi) exchangeDeclare() error {
	err := m.ch.ExchangeDeclare(
		m.exName,
		m.exType,
		true,  //durable
		false, //auto-deleted
		false, //internal
		false, //no-wait
		nil,   //arguments
	)
	return err
}

func (m *mqApi) queueDeclare() error {
	q, err := m.ch.QueueDeclare(
		m.qName,
		true,  //durable
		false, //delete when unused
		false, //exclusive
		false, //no-wait
		nil,   //arguments
	)
	m.q = q
	log.Printf("Queue Declare is %s, err: %s", q.Name, err)
	return err
}

func (m *mqApi) queueBind() error {
	err := m.ch.QueueBind(
		m.q.Name,
		m.routingKey,
		m.exName,
		false,
		nil,
	)
	log.Printf("Binding queue %s to exchange %s with routing key %s", m.q.Name, m.exName, m.routingKey)
	return err
}

func (m mqApi) publish() error {
	err := m.ch.Publish(
		m.exName,
		m.routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(m.msg),
		})
	defer m.ch.Close()
	defer m.conn.Close()
	return err
}
func Dail(exName, qName, exType, routingKey, msg string) error {
	mApi := setConnection(exName, qName, exType, routingKey, msg)
	return dailMQ(mApi)
}

func dailMQ(m *mqApi) error {
	err := m.exchangeDeclare()
	if err != nil {
		FailOnError(err, "Failed to Declare exchange")
		return err
	}

	err = m.queueDeclare()
	if err != nil {
		FailOnError(err, "Failed to Declare queue")
		return err
	}

	err = m.queueBind()
	if err != nil {
		FailOnError(err, "Failed to bind queue")
		return err
	}

	err = m.publish()
	if err != nil {
		FailOnError(err, "Failed to publish a message")
	}
	return err
}

func NewConsumer(exName, qName, exType, routingKey string) <-chan amqp.Delivery {
	mCum := setConnection(exName, qName, exType, routingKey, "")
	return consumeMsg(mCum)
}

func consumeMsg(m *mqApi) <-chan amqp.Delivery {
	msgs, err := m.ch.Consume(
		m.qName,
		"",
		false,
		false,
		false,
		false, nil,
	)
	if err != nil {
		FailOnError(err, "Failed to receive message")
	}
	return msgs
}

func main() {
	// This is test to send a message to MQ
	// err := lib.Dail("EX-GO", "Q-GO", "topic", "Q-GO", "hello")
	// if err != nil {
	// 	fmt.Printf("出错了！%s", err)
	// }
	// fmt.Printf("sucess publish mesage")

	// This is test to receive a message from MQ Queue
	msgs := lib.NewConsumer("EX-GO", "Q-GO", "topic", "Q-GO")
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			fmt.Printf("[x] %s", d.Body)
			d.Ack(false)
		}
	}()
	<-forever
}
