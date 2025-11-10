package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
)

// KeyValueStore represents the in-memory key-value store.
type KeyValueStore struct {
	mu    sync.RWMutex
	data  map[string]string
	peers []string
}

// NewKeyValueStore creates a new key-value store.
func NewKeyValueStore(peers []string) *KeyValueStore {
	return &KeyValueStore{
		data:  make(map[string]string),
		peers: peers,
	}
}

// Put adds or updates a key-value pair.
func (s *KeyValueStore) Put(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
	log.Printf("PUT: %s = %s\n", key, value)
}

// Get retrieves a value by key.
func (s *KeyValueStore) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	return val, ok
}

// Delete removes a key-value pair.
func (s *KeyValueStore) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
	log.Printf("DELETE: %s\n", key)
}

// replicate sends a PUT or DELETE request to all peer instances.
func (s *KeyValueStore) replicate(method, url string, body interface{}) {
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	_, err := client.Do(req)
	if err != nil {
		log.Printf("Replication failed to %s: %v\n", url, err)
	}
}

// PutWithReplication updates the local store and replicates to peers.
func (s *KeyValueStore) PutWithReplication(key, value string) {
	s.Put(key, value)
	for _, peer := range s.peers {
		go s.replicate(http.MethodPut, fmt.Sprintf("http://%s/put", peer), map[string]string{"key": key, "value": value})
	}
}

// DeleteWithReplication deletes locally and replicates to peers.
func (s *KeyValueStore) DeleteWithReplication(key string) {
	s.Delete(key)
	for _, peer := range s.peers {
		go s.replicate(http.MethodDelete, fmt.Sprintf("http://%s/delete?key=%s", peer, key), nil)
	}
}

// HTTP Handlers
func (s *KeyValueStore) putHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	s.PutWithReplication(req.Key, req.Value)
	w.WriteHeader(http.StatusOK)
}

func (s *KeyValueStore) getHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Key is required", http.StatusBadRequest)
		return
	}
	val, ok := s.Get(key)
	if !ok {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{key: val})
}

func (s *KeyValueStore) deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Key is required", http.StatusBadRequest)
		return
	}
	s.DeleteWithReplication(key)
	w.WriteHeader(http.StatusOK)
}

func main() {
	port := flag.String("port", ":8080", "Port to listen on")
	peerList := flag.String("peers", "", "Comma-separated list of peer addresses")
	flag.Parse()

	var peers []string
	if *peerList != "" {
		peers = append(peers, *peerList)
	}

	store := NewKeyValueStore(peers)

	http.HandleFunc("/put", store.putHandler)
	http.HandleFunc("/get", store.getHandler)
	http.HandleFunc("/delete", store.deleteHandler)

	log.Printf("Server started on %s with peers: %v\n", *port, peers)
	log.Fatal(http.ListenAndServe(*port, nil))
}

// go run main.go -port=:8080 -peers=localhost:8081,localhost:8082
// go run main.go -port=:8081 -peers=localhost:8080,localhost:8082
// go run main.go -port=:8082 -peers=localhost:8080,localhost:8081

// curl -X PUT -H "Content-Type: application/json" -d '{"key":"foo","value":"bar"}' http://localhost:8080/put
// curl http://localhost:8080/get?key=foo
// curl -X DELETE http://localhost:8080/delete?key=foo
