package main

import(
    "log"
    "fmt"
    "os"
    "strings"
    "github.com/streadway/amqp"
)

func failOnError(err error, msg string){
    if err != nil {
        log.Fatalf("%s: %s", msg, err)
    }
}


func main(){
    conn, err := amqp.Dial("amqp://devops:123456@114.215.102.224:5672/bbd_task")
    failOnError(err, "Faild to connect to RabbitMQ")
    defer conn.Close()

    ch, err := conn.Channel()
    failOnError(err, "Faild to open a channel")
    defer conn.Close()

    q, err := ch.QueueDeclare(
        "task_queue", //name
        true,         //durable
        false,       //delete when unused
        false,       //exclusive
        false,      //no-wait
        nil,        //arguments
    )
    failOnError(err, "Faild decclare a queue")

    body := bodyFrom(os.Args)
    err = ch.Publish(
        "", //exchange
        q.Name,   //routing key
        false,    // mandatory
        false,
        amqp.Publishing{
            DeliveryMode: amqp.Persistent,
            ContentType: "text/plain",
            Body:        []byte(body),
        })

    failOnError(err, "Faild to publish a message")
    log.Printf("[x] Sent %s", body)

}

func bodyFrom(args []string) string{
    var s string
    if (len(args)<2) || os.Args[1] == ""{
        s = "hello"
    }else{
        s = strings.Join(args[1:], " ")
        fmt.Printf("strings Join:[%s] ", s)
    }
    return s
}
