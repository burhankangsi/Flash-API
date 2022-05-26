package content

import (

	//"net/http"

	"io"
	"os"
	"time"

	"github.com/Shopify/sarama"
	"github.com/burhankangsi/LetsYouTube/bucket_api"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

var writer *kafka.Writer

func ConfigureProducer() (sarama.AsyncProducer, error) {

	var asyncProducer sarama.AsyncProducer
	var err error

	config := sarama.NewConfig()
	config.Metadata.RefreshFrequency = time.Duration(20) * time.Second
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Compression = sarama.CompressionNone

	kafkaMsg := &bucket_api.KafkaMsg{}
	kafkaMsg.Brokers = []string{"efd", "fgr", "acd"}
	kafkaMsg.Topic = ""

	asyncProducer, err = sarama.NewAsyncProducer(kafkaMsg.Brokers, config)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	defer func() {
		if err := asyncProducer.Close(); err != nil {
			log.Error("Failed to shutdown producer cleanly", err)
		}
	}()

	log.Info("Connected to brokers %v", kafkaMsg.Brokers)
	log.Info("Created producer for topic %v", kafkaMsg.Topic)
	return asyncProducer, nil
}

func UploadToTopic(prod sarama.AsyncProducer, video string) error {

	file, err := os.Open(video)
	if err != nil {
		log.Errorf("Could'nt open video as file. %v", err)
		return err
	}
	defer file.Close()
	defer prod.Close()

	buf := make([]byte, 1316)
	currOffset := 0
	topic := ""

	log.Infof("Started pushing video to kafka")
	for {
		reader, err := file.Read(buf)
		if err == io.EOF {
			log.Infof("End of file %v. Offset %v", video, currOffset)
			err = nil
			break
		}

		if err = pushMsg([]byte(topic), buf[:reader], topic, prod); err != nil {
			log.Info("Error sending msg to topic. %v", err)
			break
		}
		currOffset++
	}
	return nil
}

func pushMsg(key []byte, value []byte, topic string, prod sarama.AsyncProducer) error {

	var sarama_msg sarama.ProducerMessage
	if key != nil {
		sarama_msg.Key = sarama.ByteEncoder(key)
	}
	sarama_msg.Value = sarama.ByteEncoder(value)
	sarama_msg.Topic = topic

	prod.Input() <- &sarama_msg
	log.Debugf("Successfully sent msg to kafka. %v", string(value))

	return nil
}
