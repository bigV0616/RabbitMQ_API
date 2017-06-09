package main

import (
    "log"
    "github.com/streadway/amqp"
)

func failOnError(err error, msg string){
    if err != nil{
        log.Fatalf("%s:%s", msg, err)
    }
}

func main(){
    conn, err := amqp.Dial("amqp://devops:123456@114.215.102.224:5672/bbd_task")
    failOnError(err, "Faild to connect to RabbitMQ")
    defer conn.Close()

    ch, err := conn.Channel()
    failOnError(err, "Faild to open a channel")
    defer ch.Close()

    q, err := ch.QueueDeclare(
        "hello", //name
        false, //durable
        false, //delete when unused
        false, // exclusive
        false, //no-wait
        nil,   // arguments
    )
    failOnError(err, "Faild to declare a queue")

    body := "hello"
    err = ch.Publish(
        "",   //exchange
        q.Name, //routing key
        false,
        false,
        amqp.Publishing{
            ContentType: "text/plain",
            Body: []byte(body),
        })

    log.Printf("[x] sent %s", body)
    failOnError(err, "Faild to publish a message")


}
