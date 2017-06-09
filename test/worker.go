package main
import(
    "bytes"
    "log"
    "time"
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
        "task_queue",
        true,
        false,
        false,
        false,
        nil,
    )
    failOnError(err, "Failed to declare a queue")

    err = ch.Qos(
        1,     //prefetch cont
        0,     //prefetch size
        false, //global
    )
    failOnError(err, "Failed to set QoS")

    msgs, err := ch.Consume(
        q.Name,    //queue
        "",        //consumer
        false,     //auto-ack
        false,     //exclusive
        false,     //no-local
        false,     //no-wait
        nil,       //args
    )
    failOnError(err, "Failed to register a consumer")
    forever := make(chan bool)

    go func(){
        for d := range msgs{
            log.Printf("Received a message: %s", d.Body)
            dot_count := bytes.Count(d.Body, []byte("."))
            log.Println(dot_count)
            t := time.Duration(dot_count)
            time.Sleep(t * time.Second)
            log.Printf("Done")
            d.Ack(false)
        }
    }()
    log.Printf("[*] Waiting for messages, To exit press CTRL+C")
    <- forever
}
