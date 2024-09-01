package web

import (
	"distributed-kv/db"
	"fmt"
	"hash/fnv"
	"io"
	"net/http"
)

type Server struct {
	db         *db.Database
	shardId    int
	shardCount int
	addrMap    map[int]string
}

// Returns a newly created server instance
func NewServer(db *db.Database, shardId int, shardCount int, addrMap map[int]string) *Server {
	fmt.Println("Constructor :", shardCount)
	return &Server{db: db, shardId: shardId, shardCount: shardCount, addrMap: addrMap}
}

func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		// Handle the error in case of a problem with parsing
		fmt.Fprintf(w, "Error parsing form: %v", err)
		return
	}

	key := r.Form.Get("key")
	value, err := s.db.GetKey(key)

	shardId := s.getShardId(key)

	if shardId != s.shardId {
		redirect(w, "http://"+s.addrMap[shardId]+r.RequestURI)
		return
	}

	fmt.Fprintf(w, "Value = %q, Error: %v, ShardId: %d", value, err, shardId)
}

func (s *Server) SetHandler(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		// Handle the error in case of a problem with parsing
		fmt.Fprintf(w, "Error parsing form: %v", err)
		return
	}

	key := r.Form.Get("key")
	value := r.Form.Get("value")

	shardId := s.getShardId(key)

	if shardId != s.shardId {
		redirect(w, "http://"+s.addrMap[shardId]+r.RequestURI)
		return
	}

	err := s.db.SetKey(key, value)

	fmt.Fprintf(w, "Value = %q, Error: %v, ShardId: %d", value, err, shardId)
}

func redirect(w http.ResponseWriter, url string) {
	fmt.Fprintf(w, "Directing request to url: %q", url)

	resp, err := http.Get(url)

	if err != nil {
		fmt.Fprintf(w, "Error getting value: %v", err)
	}

	if resp == nil || resp.Body == nil {
		http.Error(w, "Received invalid response", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	io.Copy(w, resp.Body)
}

func (s *Server) getShardId(key string) int {
	fmt.Println("Printing :", s.shardCount)
	h := fnv.New64()
	h.Write([]byte(key))
	return int(h.Sum64() % uint64(s.shardCount))
}
