package main

import (
	"flag"
	"log"

	"github.com/nats-io/stan.go"
)

var (
	addr      = flag.String("addr", "localhost:4222", "addr to stan")
	clusterID = flag.String("cluster_id", "some_cluster", "id for connection to stan cluster")
	subject   = flag.String("subject", "some_subject", "stan subject for publishing")
)

func main() {
	flag.Parse()

	connect, err := stan.Connect(*clusterID, "cache-heater-enabler",
		stan.Pings(60, 2*60),
		stan.NatsURL(*addr),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}))
	if err != nil {
		log.Fatal(err)
	}

	_, err = connect.PublishAsync(*subject, []byte("start"), func(_ string, _ error) {})
	if err != nil {
		log.Fatalln(err)
	}

	if err = connect.Close(); err != nil {
		log.Fatal(err)
	}
}
