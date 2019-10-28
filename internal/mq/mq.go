package mq

import (
	"context"
	"fmt"
	"github.com/omerkaya1/watcher/internal/config"
	"github.com/omerkaya1/watcher/internal/errors"
	"github.com/omerkaya1/watcher/internal/interfaces"
	"github.com/streadway/amqp"
	"log"
	"os"
	"os/signal"
	"time"
)

// EventMQProducer .
type EventMQProducer struct {
	Conn *amqp.Connection
	db   interfaces.EventStorageProcessor
	conf config.QueueConf
}

// NewEventMQProducer .
func NewEventMQProducer(conf config.QueueConf, db interfaces.EventStorageProcessor) (*EventMQProducer, error) {
	if conf.Host == "" || conf.Port == "" || conf.User == "" || conf.Password == "" || conf.QueueName == "" {
		return nil, errors.ErrBadQueueConfiguration
	}
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", conf.User, conf.Password, conf.Host, conf.Port))
	if err != nil {
		return nil, err
	}
	return &EventMQProducer{Conn: conn, conf: conf, db: db}, nil
}

// ProduceMessages .
func (eq *EventMQProducer) ProduceMessages() error {
	ch, err := eq.Conn.Channel()
	if err != nil {
		return err
	}
	q, err := ch.QueueDeclare(
		eq.conf.QueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	interval, err := time.ParseDuration(eq.conf.Interval)
	if err != nil {
		return err
	}

	// Handle interrupt signal
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)
	// Create a ticker to trigger the the scan process and do the DB query
	tickTockBoom := time.NewTicker(interval)

MQ:
	for {
		select {
		case <-stopChan:
			ch.Close()
			eq.Conn.Close()
			log.Println("Exit the programme.")
			break MQ
		case <-tickTockBoom.C:
			ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
			events, err := eq.db.GetUpcomingEvents(ctx)
			if err != nil {
				log.Printf("%s: %s", errors.ErrMQPrefix, err)
				break
			}
			if events != nil {
				for _, e := range events {
					body := fmt.Sprintf("User: %s has '%s' from %s until %s",
						e.UserName, e.EventName, e.StartTime, e.EndTime)
					err = ch.Publish(
						"",
						q.Name,
						false,
						false,
						amqp.Publishing{
							ContentType: "application/json",
							Body:        []byte(body),
						})
					if err != nil {
						log.Printf("%s: %s", errors.ErrMQPrefix, err)
					}
				}
			}
		}
	}
	return nil
}
