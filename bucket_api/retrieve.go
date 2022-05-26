package bucket_api

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"

	//"net/http"
	"encoding/gob"
	"os"
	"os/signal"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"

	// "github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/service/S3"
	// "github.com/aws/aws-sdk-go/aws/session"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	//"github.com/segmentio/kafka-go"
)

type KafkaMsg struct {
	Topic      string
	Partition  int
	Key, Value []byte
	Offset     int
	Brokers    []string
	Topics     []string
	ClientId   string
}
type retrieve struct {
	pipeReader *io.PipeReader
	pipeWriter *io.PipeWriter
}

func DecodeAndUpload(reader *bufio.Reader, context context.Context) error {
	// Calling Pipe method
	ret := retrieve{}
	ret.pipeReader, ret.pipeWriter = io.Pipe()

	go func() {
		defer ret.pipeWriter.Close()
		//remember to cancel the context also

		decoder := gob.NewDecoder(reader)
		interruptChan := make(chan os.Signal, 1)
		signal.Notify(interruptChan, syscall.SIGINT, syscall.SIGTERM)

		for {
			select {
			case <-interruptChan:
				log.Info("DecodeAndSend got an interrupt")
				break
			case <-context.Done():
				log.Info("DecodeAndSend got an interrupt")
				break
			default:
				msg := KafkaMsg{}
				err := decoder.Decode(&msg)
				if err != nil {
					log.Info("DecodeAndSend Got an error while decoding msg, %v", err)
					break
				}
				// ret.pipeWriter.Write([]byte(msg))
				// log.Info("DecodeAndSend: Msg written to pipe, msg length %v", len(msg.Value))

				buffer := new(bytes.Buffer)
				buffer.ReadFrom(ret.pipeReader)
				res := buffer.Bytes()

				// Upload Files
				err = uploadFile(res)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
		decoder = nil
	}()
	return nil
}

func DownlaodVideoFromKafka() {

	context := context.Background()

	kafkaConfig := KafkaMsg{}
	kafkaBrokerUrl := ""
	brokers := strings.Split(kafkaBrokerUrl, ",")
	kafkaConfig.Brokers = brokers
	kafkaConfig.Topics = []string{}

	// kafkaconfig := kafka.ReaderConfig{
	// 	Brokers:         brokers,
	// 	GroupID:         kafkaClientId,
	// 	Topic:           kafkaTopic,
	// 	MinBytes:        10e3,            // 10KB
	// 	MaxBytes:        10e6,            // 10MB
	// 	MaxWait:         1 * time.Second, // Maximum amount of time to wait for new data to come when fetching batches of messages from kafka.
	// 	ReadLagInterval: -1,
	// }
	consumer, err1 := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "host1:9092,host2:9092",
		"group.id":          "foo",
		"auto.offset.reset": "smallest"})

	if err1 != nil {
		log.Fatal("Error in creating a consumer")
	}

	//reader := kafka.NewReader(kafkaconfig)
	//defer reader.Close()
	pipeReader, pipeWriter := io.Pipe()
	defer pipeWriter.Close()
	vidReader := bufio.NewReader(pipeReader)
	encoder := gob.NewEncoder(pipeWriter)

	go func() {
		if err := DecodeAndUpload(vidReader, context); err != nil {
			log.Infof("Could not download and send video, %v", err)
		}
	}()

	// interruptChan := make(chan os.Signal, 1)
	// signal.Notify(interruptChan, syscall.SIGINT, syscall.SIGTERM)
	var err error
	err = consumer.SubscribeTopics(kafkaConfig.Topics, nil)

	// for {
	// 	select {
	// 	case <-interruptChan:
	// 		log.Info("DownlaodVideoFromKafka got an interrupt")
	// 		err = nil
	// 		break

	// 	case <-context.Done():
	// 		log.Info("DownlaodVideoFromKafka got a context closure")
	// 		err = nil
	// 		break
	// 	default:
	// msg, err := reader.ReadMessage(context)
	// if err != nil {
	// 	//log.Error().Msgf("error while receiving message: %s", err.Error())
	// 	log.Errorf("error while receiving message: %s", err.Error())
	// 	continue
	// } else if msg != nil {
	// 	log.Infof("message at topic/partition/offset %v/%v/%v: %s\n", msg.Topic, msg.Partition, msg.Offset, string(msg.value))

	// 	if err = encoder.Encode(msg.value); err != nil {
	// 		err = fmt.Errorf("Error while encoding with gob: %v", err)
	// 		break
	// 	}
	// } else {
	// 	log.Infof("Messages are not coming")
	// }
	// 	}
	// }

	// if err := reader.Close(); err != nil {
	// 	log.Fatal("failed to close reader:", err)
	// }

	go func() {
		var run bool
		for run == true {
			ev := consumer.Poll(0)
			switch e := ev.(type) {
			case *kafka.Message:
				fmt.Printf("%% Message on %s:\n%s\n",
					e.TopicPartition, string(e.Value))
				if err = encoder.Encode(e.Value); err != nil {
					err = fmt.Errorf("Error while encoding with gob: %v", err)
					break
				}
			case kafka.PartitionEOF:
				fmt.Printf("%% Reached %v\n", e)
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
				run = false
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}()

	consumer.Close()
}
