package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"netshield/agent/internal/monitor"
	"netshield/agent/internal/wifi"
	"os"
	"os/signal"
	"time"
)

func main() {
	wm := wifi.WindowsManager{}

	cfg := monitor.Config{
		MinSignalPercent: 60,
		MaxAvgPingMs:     120,
		PingHost:         "8.8.8.8",
		CheckInterval:    10 * time.Second,
		PreferredProfiles: []string{
			"esperance",
			"KIIT-WIFI-DU",
			"OPPO A9 2020",
		},
	}

	m := &monitor.Monitor{
		Wifi:   wm,
		Config: cfg,
	}

	go startLocalAPI(m)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := m.Start(ctx); err != nil && err != context.Canceled {
		log.Println("monitor stopped with error:", err)
	}
}

func startLocalAPI(m *monitor.Monitor) {
	mux := http.NewServeMux()

	mux.HandleFunc("/current", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		snap := m.GetSnapshot()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(snap)
	})

	addr := "127.0.0.1:9090"
	log.Println("[agent] local API on http://" + addr + "/current")
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Println("[agent] local api error:", err)
	}
}
