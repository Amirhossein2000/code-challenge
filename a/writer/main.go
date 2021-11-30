package main

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/apsdehal/go-logger"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type Order struct {
	ID    int
	Price int
	Title string
}

var tableName = "orders"
var bufferSize = 4096
var logger *log.Logger

func init() {
	var err error
	logger, err = log.New("go-server", 1, os.Stdout)
	if err != nil {
		panic(fmt.Sprintf("create new logger err: %s", err.Error()))
	}

	strBuf := os.Getenv("BUFFER_SIZE")
	if strBuf != "" {
		bs, err := strconv.Atoi(strBuf)
		if err != nil {
			logger.Warningf("invalid buffer %v size: %s", bs, err.Error())
			return
		} else {
			bufferSize = bs
		}
	}

	tbName := os.Getenv("BUFFER_SIZE")
	if tbName != "" {
		tableName = tbName
	}
}

func main() {
	orderChan := make(chan *Order, bufferSize)
	ctx, cancel := context.WithCancel(context.Background())

	go startSubscribing(ctx, orderChan)
	go startWriting(orderChan)

	signChan := make(chan os.Signal, 1)
	signal.Notify(signChan, syscall.SIGINT)
	<-signChan

	cancel()
	err := getRedisClient().Close()
	if err != nil {
		logger.Errorf("close redis client err: %s", err.Error())
	}

	for {
		<-time.After(time.Millisecond * 100)

		remainingOrders := len(orderChan)
		if remainingOrders != 0 {
			logger.Infof("remaining orders: %d", remainingOrders)
			continue
		}

		close(orderChan)
		err = getDB().Close()
		if err != nil {
			logger.Errorf("close msql client err: %s", err.Error())
		}
		return
	}
}

func startSubscribing(ctx context.Context, orderChan chan<- *Order) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := getRedisClient().Subscribe(redisChan).ReceiveMessage()
			if err != nil {
				logger.Errorf("receive message err: %s", err.Error())
				continue
			}

			order := Order{}
			err = json.Unmarshal([]byte(msg.Payload), &order)
			if err != nil {
				logger.Errorf("unmarshal order err: %s", err.Error())
				continue
			}

			orderChan <- &order
		}
	}
}

func startWriting(orderChan <-chan *Order) {
	for order := range orderChan {
		_, err := getDB().Query(`INSERT INTO ? VALUES (?,?,?);`, tableName, order.ID, order.Price, order.Title)
		if err != nil {
			logger.Errorf("insert query err: %s", err.Error())
			continue
		}

		logger.Infof("order: %+v has inserted into db successfully", order)
	}
}
