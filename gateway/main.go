package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"tollCalculator.com/aggregator/client"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func main() {
	listenAddress := flag.String("listenAddr", ":6000", "the listen address of the HTTP server")
	flag.Parse()
	invoiceHandler := &InvoiceHandler{
		client: client.NewHTTPClient("http://localhost:3000"),
	}
	http.HandleFunc("/invoice", makeAPIFunc(invoiceHandler.handleGetInvoice))
	logrus.Infof("gateway HTTP server running on port %s", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}

type InvoiceHandler struct {
	client client.Client
}

func (h *InvoiceHandler) handleGetInvoice(w http.ResponseWriter, r *http.Request) error {
	values, ok := r.URL.Query()["obu"]
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing obu ID"})
		return nil
	}
	obuID, err := strconv.Atoi(values[0])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid obu ID"})
		return nil
	}
	inv, err := h.client.GetInvoice(context.Background(), obuID)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, inv)

}

func writeJSON(w http.ResponseWriter, code int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func makeAPIFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(start time.Time) {
			logrus.WithFields(logrus.Fields{
				"took": time.Since(start),
				"uri":  r.RequestURI,
			}).Info("Gateway")
		}(time.Now())
		if err := fn(w, r); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

	}
}
