package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// import "github.com/go-redis/redis/v9"

var ctx = context.Background()

var _ = uuid.New()
var _ = strconv.Itoa(1)
var _ sync.WaitGroup
var _, _ = fmt.Println("")

var PRICE = 650
var RSI = 28

type RSIStruct struct {
	timeframe int
	period    int
	match     string
	value     int
}

type PriceStruct struct {
	match string
	value int
}

type ConfigStruct struct {
	id    string
	rsi   RSIStruct
	price PriceStruct
}

type dataType []string

func main() {

	// CONNECT TO THE REDIS DATABASE
	// rdb := redis.NewClient(&redis.Options{
	//     Addr:     "localhost:6379",
	//     Password: "", // no password set
	//     DB:       0,  // use default DB
	// })

	// GET ALL THE SYMBOLS
	var symbols dataType
	file, err := ioutil.ReadFile("currency.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(file, &symbols)
	if err != nil {
		log.Fatal(err)
	}

	totalStartTime := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < 1; i++ {
		wg.Add(1)
		go func(i int) {
			// // FETCH THE DATA
			// fetchStartTime := time.Now()
			// values, err := rdb.LRange(ctx, symbols[36], 0, -1).Result()
			// if err != nil {
			// 	panic(err)
			// }
			// fmt.Println("fetched records: ", len(values))
			// fetchingTime := time.Since(fetchStartTime)
			// log.Printf("Fetching done in %s", fetchingTime)

			// COMPARE THE RESULT
			compareStartTime := time.Now()
			for i := 0; i < 6000000; i++ {
				config := parseConfigStr("16:P_LT_944:RSI_15_14_LT_29")
				match(config)
			}
			compareTime := time.Since(compareStartTime)
			log.Printf("Comparison done in %s", compareTime)

			wg.Done()
		}(i)
	}
	wg.Wait()
	totalTime := time.Since(totalStartTime)
	log.Printf("Total validation done in %s", totalTime)
}

func match(config ConfigStruct) bool {
	// Price match
	isPriceMatched := false
	if config.price.match == "LT" && !(PRICE < config.price.value) {
		return false
	} else if config.price.match == "GT" && !(PRICE > config.price.value) {
		return false
	}
	isPriceMatched = true

	// RSI Match
	isRSIMatched := false
	if config.rsi.match == "LT" && !(RSI < config.rsi.value) {
		return false
	} else if config.rsi.match == "GT" && !(RSI > config.rsi.value) {
		return false
	}
	isRSIMatched = true

	return isPriceMatched && isRSIMatched
}

func parseConfigStr(value string) ConfigStruct {
	// value: 16:P_LT_944:RSI_15_14_LT_29
	res := strings.Split(value, ":")
	id := res[0]
	priceVal := res[1]
	rsiVal := res[2]

	var rsi RSIStruct
	rsiValues := strings.Split(rsiVal, "_")
	rsiTimeFrame, _ := strconv.Atoi(rsiValues[1])
	rsiPeriod, _ := strconv.Atoi(rsiValues[2])
	rsiValue, _ := strconv.Atoi(rsiValues[4])
	rsi.timeframe = rsiTimeFrame
	rsi.period = rsiPeriod
	rsi.match = rsiValues[3]
	rsi.value = rsiValue

	var price PriceStruct
	priceValues := strings.Split(priceVal, "_")
	priceValue, _ := strconv.Atoi(priceValues[2])
	price.match = priceValues[1]
	price.value = priceValue

	var config ConfigStruct
	config.id = id
	config.rsi = rsi
	config.price = price

	return config
}
