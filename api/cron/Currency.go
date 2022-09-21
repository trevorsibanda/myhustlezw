package cron

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/leekchan/accounting"
	"github.com/trevorsibanda/myhustlezw/api/cache"
)

var usdToZWL int

func GetExchangeRate() int {
	return usdToZWL
}

var (
	accountingZWL = accounting.DefaultAccounting("ZWL $", 1)
	accountingUSD = accounting.DefaultAccounting("USD $", 2)
)

//ZWLToUSD converts ZWL to its USD equivalent
func ZWLToUSD(zwl float64) (usd float64) {
	return zwl / 100.0
}

func Format(currency string, amount float64) string {
	if strings.ToLower(currency) == "zwl" {
		return FormatAsZWL(amount)
	}
	return FormatAsUSD(amount)
}

//USDToZWL converts USD to ZWL price
func USDToZWL(usd float64) (zwl float64) {
	//return usd * 50
	return float64(GetExchangeRate()) * usd
}

//FormatAsZWL formats a price as ZWL
func FormatAsZWL(price float64) string {
	return accountingZWL.FormatMoneyFloat64(price)
}

//FormatAsUSD formats as USD
func FormatAsUSD(price float64) string {
	return accountingUSD.FormatMoneyFloat64(price)
}

func UpdateExchangeRate() {
	log.Printf("Retrieved currency rate at %s", time.Now())
	store := cache.RedisStore()
	cmd := store.Get(context.Background(), "currency_usd_to_zwl")
	if rate, err := cmd.Int(); err != nil {
		log.Printf("Failed to get currency rate with error %s", err)
		log.Printf("Loading default values")
		LoadDefaultExchangeRate()
		store.Set(context.Background(), "currency_usd_to_zwl", usdToZWL, time.Minute*45).Result()
		return
	} else {
		usdToZWL = rate
	}
	log.Printf("Updated exchange rate to $USD1: ZWL$ %d", usdToZWL)
}

func LoadDefaultExchangeRate() {
	if rate, err := strconv.Atoi(os.Getenv("USD_TO_ZWL")); err != nil {
		panic(fmt.Sprintf("Failed to get default exchange rate with error: %s ", err))

	} else {
		if rate < usdToZWL {
			log.Printf("Ignoring config default exchange rate of %d since current %d is higher", rate, usdToZWL)
			return
		}
		usdToZWL = rate
		log.Printf("Loaded exchange rate from env var: $USD1 : ZWL$%d", usdToZWL)
	}
}
