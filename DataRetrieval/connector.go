package main

import (
  "fmt"
  "os"
  "os/signal"
  "net/http"
  "io"
  "net/url"
  //forecast "github.com/mlbright/darksky/v2"
  "log"
	"io/ioutil"

  "github.com/Shopify/sarama"
	)

var brokers = []string{"kafka:9092"} // Change to localhost depending on OS, windows refer to this string format--> "localhost:9092"

func newProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true // The config.Producer.Return.Successes  must return true for the producer to generate values onto the Topic
	producer, err := sarama.NewSyncProducer(brokers, config)

	return producer, err
}

func prepareMessage(topic, message string) *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: -3,
		Value:     sarama.StringEncoder(message),
	}

	return msg
}

func Download(filepath string, url string) error {

    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    out, err := os.Create(filepath)
    if err != nil {
        return err
    }
    defer out.Close()
    _, err = io.Copy(out, resp.Body)
    return err
}

func IsValidUrl(str string) bool {
   u, err := url.Parse(str)
   return err == nil && u.Scheme != "" && u.Host != ""
}

// const message = {sessID: authData.sessID, userID: authData._id, action: 'ModelExecution', value: req.query.value, timeStamp: Date.now()}


func main() {

  resp, err := http.Get("https://api.darksky.net/forecast/68a391b503f11aa6fa13d405bfefdaba/43.6595,79.3433")
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	//log.Println(string(body))
  /*test := "68a391b503f11aa6fa13d405bfefdaba"
  key := string(test)
  key = strings.TrimSpace(key)

  lat := "43.6595"
  long := "-79.3433"

  file, err := forecast.Get(key, lat, long, "now", forecast.CA, forecast.English)
  if err != nil {
      log.Fatal(err)
  }
  fmt.Printf("%s:\n", file)
  //fmt.Printf("humidity: %.2f\n", file.Currently.Humidity)
  //fmt.Printf("temperature: %.2f Celsius\n", file.Currently.Temperature)
  //fmt.Printf("wind speed: %.2f\n", file.Currently.WindSpeed)
  */



	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true // The config.Producer.Return.Successes  must return true for the producer to generate values onto the Topic

  // consumer config declarations

  config.Consumer.Return.Errors = true

  // consumer inititalisation block
  // fileUrl := "https://engineering.arm.gov/~jhelmus/pyart_example_data/Level2_KATX_20130717_1950.ar2v"

  masterConsumer, errConsumer := sarama.NewConsumer(brokers, config) //the NewConsumer allows for the brokers to be addedwhennew topics are created
	if errConsumer != nil {
		panic(errConsumer)
	}

  masterProducer, errProducer := sarama.NewSyncProducer(brokers, config)
	if errProducer != nil {
		panic(errProducer)
	}

  defer func() {
		if errConsumer := masterConsumer.Close(); errConsumer != nil {
			panic(errConsumer)
		}
	}()

  //offset inititalisation






  defer func() {
		if errProducer := masterProducer.Close(); errProducer != nil {
			panic(errProducer)
		}
	}()

  topicProducer := "DataRetrieval"

  //message := "message string"
  for i := 0; i < 10; i++ {
		if errProducer != nil {
			panic(errProducer)  // Checks for the err returned by NewSyncProducer
		}
	}

  topic := "apigateway"
	consumer, err := masterConsumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}

  // producer func() activated here

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	count := 0
	finished := make(chan struct{}) // channel  created
	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				fmt.Println(err)
			case msg := <-consumer.Messages():
				fmt.Println("Consumer Initialised")
				count++
				fmt.Println("The messages : %v", string(msg.Value))


        link := &sarama.ProducerMessage{
                         Topic: topicProducer,
                         Value: sarama.StringEncoder(string(body)),
                   }
        partition, offset, err := masterProducer.SendMessage(link)
        fmt.Println("Producer produced")
        if err != nil {
           panic(err)
         }

      //  if err := Download("Level2_KATX_20130717_1950.ar2v", fileUrl); err != nil {
      //      panic(err)
      //  }
        fmt.Println("The link (%s) has been sent with partition(%d)/offset(%d)",link.Value,partition,offset)

			case <-signals:
				fmt.Println("Interrupt is detected")
				finished <- struct{}{}
			}
		}
	}()

  <-finished // channel finished
  fmt.Println("Processed", count, "messages")

}
