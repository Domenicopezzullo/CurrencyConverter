package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"github.com/charmbracelet/huh"
)

var (
	fromCurrency string
	toCurrency   string
	amount       string
)

func main() {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("From currency").
				Options(
					huh.NewOption("USD", "USD"),
					huh.NewOption("EUR", "EUR"),
					huh.NewOption("GBP", "GBP"),
					huh.NewOption("JPY", "JPY"),
					huh.NewOption("CNY", "CNY"),
				).
				Value(&fromCurrency),
			huh.NewSelect[string]().
				Title("To currency").
				Options(
					huh.NewOption("USD", "USD"),
					huh.NewOption("EUR", "EUR"),
					huh.NewOption("GBP", "GBP"),
					huh.NewOption("JPY", "JPY"),
					huh.NewOption("CNY", "CNY"),
				).
				Value(&toCurrency),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Amount").
				Value(&amount).
				Validate(func(v string) error {
					val, err := strconv.ParseFloat(v, 64)
					if err != nil {
						log.Fatal(err)
					}

					if val <= 0 {
						return errors.New("amount must be greater than 0")
					}
					return nil
				}),
		),
	)

	if err := form.Run(); err != nil {
		log.Fatal(err)
	}

	parsedAmount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		log.Fatal(err)
	}

	if fromCurrency == toCurrency {
		fmt.Println("Same currency -_-")
		return
	}

	apiKey, ok := os.LookupEnv("CURRENCYAPIKEY")
	if !ok {
		log.Fatal("APIKEY is not set")
		os.Exit(1)
	}

	res, err := http.Get("https://v6.exchangerate-api.com/v6/" + apiKey + "/pair/" + fromCurrency + "/" + toCurrency + "/" + amount)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}

	result := data["conversion_result"]
	hardCodedAmount := result.(float64)

	printAmount := func(value float64) string {
		if value == float64(int(value)) {
			return fmt.Sprintf("%d", int(value))
		}
		return fmt.Sprintf("%.2f", value)
	}

	fmt.Printf("%s %s is worth around %s %s\n", printAmount(parsedAmount), fromCurrency, printAmount(hardCodedAmount), toCurrency)
}

