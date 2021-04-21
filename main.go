package main

import (
	"flag"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	gecko "github.com/superoo7/go-gecko/v3"
)

var listenAddress = flag.String("listen-address", ":9400", "The address this exporter would listen on")

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
		log.Error("Could not get rate for \"", currency, "\", got error: ", err)
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
	log.Info("GET /metrics/rates/", baseCurrency, "?currency=", currency)
}

func main() {
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/metrics/rates/{base_currency}", RatesHandler)
	http.Handle("/", r)

	log.Info("Listening on ", *listenAddress)
	err := http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		log.Fatal("Could not start application at ", *listenAddress, ", got error: ", err)
		panic(err)
	}
}
