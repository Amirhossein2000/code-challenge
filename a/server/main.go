package main

import (
	"encoding/json"
	"fmt"
	log "github.com/apsdehal/go-logger"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Order struct {
	ID    int
	Price int
	Title string
}

var serverAddress = ":7070"
var logger *log.Logger

func init() {
	var err error
	logger, err = log.New("go-server", 1, os.Stdout)
	if err != nil {
		panic(fmt.Sprintf("create new logger err: %s", err.Error()))
	}

	sa := os.Getenv("SERVER_ADDRESS")
	if sa != "" {
		if len(strings.Split(sa, ":")) == 2 {
			serverAddress = sa
		} else {
			logger.Warningf("invalid SERVER_ADDRESS %s, serving on default", sa)
		}
	}
}

func main() {
	http.HandleFunc("/api/order", orderHandler)
	err := http.ListenAndServe(serverAddress, nil)
	if err != nil {
		logger.Errorf("listen err: %s", err.Error())
	}
}

func orderHandler(w http.ResponseWriter, r *http.Request) {
	order := Order{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Errorf("read body err: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &order)
	if err != nil {
		logger.Errorf("unmarshal body err: %s", err.Error())
		writeResponse(
			w,
			map[string]string{
				"result": "request struct is wrong",
			},
			http.StatusBadRequest)
		return
	}

	err = getRedisClient().Publish(redisChan, body).Err()
	if err != nil {
		logger.Errorf("publish into redis err: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	logger.Infof("order: %s has published into redis successfully", body)

	writeResponse(w,
		map[string]string{
			"result": "order has published successfully",
		},
		http.StatusCreated)
}

func writeResponse(w http.ResponseWriter, resp interface{}, statusCode int) {
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		logger.Errorf("marshal response err: %s", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(jsonResp)
	if err != nil {
		logger.Errorf("write response err: %s", err.Error())
		return
	}
}
