package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type Block struct {
	Number      int `json:"height"`
	Timestamp   string `json:"date"`
	Minedby     int `json:"mined"`
	BlockReward int `json:"reward"`
	Difficulty  int    `json:"difficult"`
}

type blockHandlers struct {
	sync.Mutex
	store map[string]Block
}

func (h *blockHandlers) blocks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.get(w, r)
		return
	case "POST":
		h.post(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("not acceptable method"))
		return
	}
}

func (h *blockHandlers) get(w http.ResponseWriter, r *http.Request) {
	blocks := make([]Block, len(h.store))

	h.Lock()
	i := 0
	for _, block := range h.store {
		blocks[i] = block
		i++
	}
	h.Unlock()

	jsonBytes, err := json.Marshal(blocks)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *blockHandlers) post(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	contenttype := r.Header.Get("content-type")
	if contenttype != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("Please add content-type 'application/json' your request but '%s' invalid", contenttype)))
		return
	}

	var block Block
	err = json.Unmarshal(requestBody, &block)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	block.Timestamp = fmt.Sprintf("%d", time.Now().UTC().UnixNano())
	h.Lock()
	h.store[block.Timestamp] = block
	defer h.Unlock()
}

func newblockHandlers() *blockHandlers {
	return &blockHandlers{
		store: map[string]Block{},
	}
}

func main() {
	blockHandlers := newblockHandlers()
	http.HandleFunc("/blocks", blockHandlers.blocks)
	err := http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		panic(err)
	}
}
