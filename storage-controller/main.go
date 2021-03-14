package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
)

var dataCh chan *ReplicationData
var dataChLen int = 5

func main() {
	dataCh = make(chan *ReplicationData, dataChLen)

	stopCh := setupSignalHandler()
	var wg sync.WaitGroup
	startReplicationWorkers(dataCh, stopCh, &wg)

	// test endpoint
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello!\n")
	})

	// /get?key=somekey
	http.HandleFunc("/get", handleGet)

	// /put?key=somkey&val=someval
	http.HandleFunc("/put", handlePut)

	// /_replicate?key=somekey&val=someval
	http.HandleFunc("/_replicate", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello!\n")
	})

	// /_fetch?key=somekey
	http.HandleFunc("/_fetch", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello!\n")
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		<-stopCh
		wg.Wait()
		log.Fatal(err)
	}

	<-stopCh
	wg.Wait()
}

type ResponseBody struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

var keyParam string = "key"
var valParam string = "val"

func handleGet(w http.ResponseWriter, r *http.Request) {
	parsed := parseParams(r.URL)
	k, err := parsed(keyParam)
	if err != nil {
		writeResponse(w, &ResponseBody{Success: false, Message: err.Error()})
		logErr(err)
		return
	}

	v, err := Read(k)
	if err != nil {
		writeResponse(w, &ResponseBody{Success: false, Message: err.Error()})
		logErr(err)
		return
	}

	writeResponse(w, &ResponseBody{Success: true, Data: v})
}

func handlePut(w http.ResponseWriter, r *http.Request) {
	parsed := parseParams(r.URL)
	k, err := parsed(keyParam)
	if err != nil {
		writeResponse(w, &ResponseBody{Success: false, Message: err.Error()})
		logErr(err)
		return
	}
	v, err := parsed(valParam)
	if err != nil {
		writeResponse(w, &ResponseBody{Success: false, Message: err.Error()})
		logErr(err)
		return
	}

	ok := Put(k, v)
	writeResponse(w, &ResponseBody{Success: ok})
}

func parseParams(url *url.URL) func(string) (string, error) {
	q := url.Query()
	return func(key string) (string, error) {
		val := q.Get(key)
		if val == "" {
			return "", fmt.Errorf("%s param not supplied", key)
		}
		return val, nil
	}
}

func writeResponse(w http.ResponseWriter, data *ResponseBody) {
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logErr(err)
	}
}

func logErr(err error) {
	fmt.Println(fmt.Errorf("error: %v", err))
}
