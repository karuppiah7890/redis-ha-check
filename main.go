package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/TogaiHQ/infra-tools/redis-ha-check/pkg/config"
	"github.com/redis/go-redis/v9"
)

func main() {
	// redis-ha-check config.yaml
	if len(os.Args) != 2 {
		printUsage()
		if len(os.Args) >= 2 {
			os.Exit(1)
		}
		return
	}

	configFilePath := os.Args[1]
	c, err := config.LoadConfig(configFilePath)
	if err != nil {
		log.Fatalf("error loading config file: %v", err)
	}

	rdb := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    c.SentinelConfig.MasterName,
		SentinelAddrs: getSentinelAddresses(c),
	})

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	go handleSignals(signals)

	ctx := context.Background()

	for {
		pingResult, err := rdb.Ping(ctx).Result()
		printRedisPingStatus(pingResult, err)
	}
}

func handleSignals(signals chan os.Signal) {
	<-signals
	os.Exit(0)
}

func printUsage() {
	fmt.Println("usage: redis-ha-check <path-to-config-yaml-file>")
}

func getSentinelAddresses(c *config.Config) []string {
	sentinelAddresses := make([]string, 0, len(c.SentinelConfig.Sentinels))

	for _, sentinel := range c.SentinelConfig.Sentinels {
		sentinelAddresses = append(sentinelAddresses, fmt.Sprintf("%s:%v", sentinel.Host, sentinel.Port))
	}

	return sentinelAddresses
}

func printRedisPingStatus(pingResult string, err error) {
	if err != nil {
		log.Printf("error occurred while pinging redis: %v\n\n", err)
		return
	}

	if pingResult == "PONG" {
		log.Printf("ping okay\n\n")
		return
	}

	log.Printf("ping not okay. ping result: %v\n\n", pingResult)
}
