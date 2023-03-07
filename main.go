package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"

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

		testKey := fmt.Sprintf("__test-redis-ha-check__%v", rand.Uint64())

		randomValue := rand.Uint64()
		setResult, err := rdb.Set(ctx, testKey, randomValue, 0).Result()
		printRedisSetStatus(setResult, err)

		getResult, err := rdb.Get(ctx, testKey).Result()
		if err != nil {
			log.Printf("error: get not working for key %s", testKey)
		}

		expectedResult := fmt.Sprintf("%v", randomValue)

		if getResult == expectedResult {
			log.Printf("got the set value, get-set okay\n\n")
		} else {
			log.Printf("error: did not get the set value, get-set not okay. expected: %v, got: %v\n\n", expectedResult, getResult)
		}

		unlinkResult, err := rdb.Unlink(ctx, testKey).Result()
		if err != nil {
			log.Printf("error: unlink not working for key %s", testKey)
		}

		if unlinkResult == 1 {
			log.Printf("unlink worked, all okay\n\n")
		} else {
			log.Printf("error: unlink did not work for key %s. you might have to manually delete the key\n\n", testKey)
		}
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

	log.Printf("error: ping not okay. ping result: %v\n\n", pingResult)
}

func printRedisSetStatus(setResult string, err error) {
	if err != nil {
		log.Printf("error occurred while setting a value in redis: %v\n\n", err)
		return
	}

	if setResult == "OK" {
		log.Printf("set okay\n\n")
		return
	}

	log.Printf("error: set not okay. ping result: %v\n\n", setResult)
}
