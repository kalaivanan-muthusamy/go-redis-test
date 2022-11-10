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
	"runtime"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
)

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
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	numcpu := runtime.NumCPU()
	// runtime.GOMAXPROCS(16)
	log.Printf("numcpu %d", numcpu)

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

	// READ THE VALUES FROM REDIS
	fetchingStartTime := time.Now()
	values, err := rdb.LRange(ctx, "CONFIGS", 0, 1000000).Result()
	if err != nil {
		log.Fatal(err)
	}
	fetchingTime := time.Since(fetchingStartTime)
	fmt.Println("fetched records: ", len(values))
	log.Printf("Fetching done in %s", fetchingTime)


	// COMPARE THE RESULT
	// compareStartTime := time.Now()
	// wg := new(sync.WaitGroup)
	// wg.Add(6000000)
	// for i := 0; i < 6000000; i++ {
	// 	go parseConfigStr("16:P_LT_944:RSI_15_14_LT_29", wg)
	// 	// config := parseConfigStr(values[i])
	// 	// match(config)
	// }
	// wg.Wait()
	// compareTime := time.Since(compareStartTime)
	// log.Printf("Comparison done in %s", compareTime)
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

func parseConfigStr(value string, wg *sync.WaitGroup) ConfigStruct {
	// value: 16:P_LT_944:RSI_15_14_LT_29
	defer wg.Done()
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
