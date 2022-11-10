package main

import (
	"encoding/json"
	// "fmt"

	"context"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-redis/redis/v9"
)

type dataType []string

var ctx = context.Background()

func main() {
	// Connect to redis database
	rdb := redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		Password:     "", // no password set
		DB:           0,  // use default DB
		ReadTimeout:  300 * time.Second,
		WriteTimeout: 300 * time.Second,
	})

	// Get sample symbols for loading the data
	var symbols dataType
	file, err := ioutil.ReadFile("currency.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(file, &symbols)
	if err != nil {
		log.Fatal(err)
	}

	start := time.Now()

	// Generate & push random configuration
	pipe := rdb.Pipeline()
	for j := 0; j < 5000000; j++ {
		id := strconv.Itoa(j)
		price := "P_LT_" + strconv.Itoa(rand.Intn(1000-100)+100)
		rsi := "RSI_15_14_LT_" + strconv.Itoa(rand.Intn(30-2)+2)
		value := id + ":" + price + ":" + rsi
		pipe.LPush(ctx, "CONFIGS", value)
	}
	cmds, err := pipe.Exec(ctx)
	if err != nil {
		panic(err)
	}
	log.Printf("Added %d", len(cmds))

	elapsed := time.Since(start)
	log.Printf("Looping took %s", elapsed)
}
