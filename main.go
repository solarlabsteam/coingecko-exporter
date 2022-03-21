package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	gecko "github.com/superoo7/go-gecko/v3"
)

var listenAddress = flag.String("listen-address", ":9400", "The address this exporter would listen on")
var jsonOutput = flag.Bool("json", false, "Output logs as JSON")

var log zerolog.Logger

func RatesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	baseCurrency := vars["base_currency"]
	currency := r.URL.Query().Get("currency")

	gauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "coingecko_rate",
			Help: "Coingecko exchange rate",
		},
		[]string{"base_currency", "currency"},
	)

	registry := prometheus.NewRegistry()
	registry.MustRegister(gauge)

	var cg = gecko.NewClient(nil)
	result, err := cg.SimplePrice(strings.Split(currency, ","), []string{baseCurrency})
	if err != nil {
		log.Error().Err(err).Str("currency", currency).Msg("Could not get rate")
		return
	}

	for currencyKey, currencyValue := range *result {
		for baseCurrencyKey, baseCurrencyValue := range currencyValue {
			gauge.With(prometheus.Labels{
				"base_currency": baseCurrencyKey,
				"currency":      currencyKey,
			}).Set(float64(baseCurrencyValue))
		}
	}

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
	log.Info().
		Str("method", "GET").
		Str("url", fmt.Sprintf("/metrics/rates/%s?currency=%s", baseCurrency, currency)).
		Msg("Processed request")
}

func main() {
	flag.Parse()

	if *jsonOutput {
		log = zerolog.New(os.Stdout).With().Timestamp().Logger()
	} else {
		log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
	}

	r := mux.NewRouter()
	r.HandleFunc("/metrics/rates/{base_currency}", RatesHandler)
	http.Handle("/", r)

	log.Info().Str("address", *listenAddress).Msg("Listening")
	err := http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not start application")
		panic(err)
	}
}
