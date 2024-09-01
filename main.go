package main

import (
	"distributed-kv/config"
	"distributed-kv/db"
	"distributed-kv/web"
	"flag"
	"log"
	"net/http"
)

var (
	dbLocation = flag.String("db-location", "", "Path to boltdb")
	httpAddr   = flag.String("http-addr", "127.0.0.1:8080", "HTTP Host and Port")
	configFile = flag.String("config-file", "sharding.toml", "Config file for sharding")
	shard      = flag.String("shard", "", "Shard to use")
)

func parseFlag() {
	flag.Parse()

	if *dbLocation == "" {
		log.Fatalf("db-location is required")
	}

	if *shard == "" {
		log.Fatalf("shard is required")
	}
}

func main() {
	parseFlag()

	c, err := config.ParseConfigFile(*configFile)

	if err != nil {
		log.Fatalf("Parse config file failed = ", err)
	}

	var shardCount int
	var shardId = -1

	var addrMap = make(map[int]string)

	shardCount = len(c.Shards)

	for _, s := range c.Shards {
		if s.Name == *shard {
			shardId = s.Id
		}

		addrMap[s.Id] = s.Addr
	}

	if shardId == -1 {
		log.Fatalf("shard %s not found", *shard)
	}

	log.Printf("shardCount: %d, shardIdx: %d", shardCount, shardId)

	dbWrapper, close, error := db.NewDatabase(*dbLocation)

	if error != nil {
		log.Fatalf("New Database not initiated, error = %v", error)
	}

	defer close()

	server := web.NewServer(dbWrapper, shardId, shardCount, addrMap)

	http.HandleFunc("/set", server.SetHandler)

	http.HandleFunc("/get", server.GetHandler)

	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}
