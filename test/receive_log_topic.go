package main

import (
    "log"
    "os"
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

    // 声明交换机
    err = ch.ExchangeDeclare(
        "logs_topic",        //name
        "topic",             //type
        true,
        false,
        false,
        false,
        nil,
    )
    failOnError(err, "Failed to declare an exchange")

    // 声明队列
    q, err := ch.QueueDeclare(
        "topic_queue",    //name
        false,
        false,
        true,
        false,
        nil,
    )
    log.Printf("%s: %s", q, err)
    failOnError(err, "Faild to declare a queue")

    // 绑定队列交换机
    for _, s := range os.Args[1:]{
        log.Printf("Binding queue %s to exchange %s with routing key %s", q.Name, "logs_topic", s)
        err = ch.QueueBind(
            q.Name,
            s,
            "logs_topic",
            false,
            nil,
        )
        failOnError(err, "Failed to bind a queue")
    }

    // 订阅消息
    msgs, err := ch.Consume(
        q.Name,
        "",
        true,
        false,
        false,
        false,
        nil,
    )

    failOnError(err, "Failed to register a consumer")
    forever := make(chan bool)
    go func(){
        for d := range msgs{
            log.Printf("[x] %s", d.Body)
        }
    }()
    log.Printf("[*] Waiting for logs. To exit press CTRL+C")
    <-forever
}

