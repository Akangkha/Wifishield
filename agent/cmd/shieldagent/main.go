package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	agentclient "netshield/agent/internal/client"
	"netshield/agent/internal/monitor"
	"netshield/agent/internal/wifi"
	agentpb "netshield/agent/proto"

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
			"vivo",
		},
	}

	m := &monitor.Monitor{
		Wifi:                wm,
		Config:              cfg,
		SwitchAutomatically: true,
	}

	serverAddr := os.Getenv("NETSHIELD_SERVER_ADDR")
	if serverAddr == "" {
		serverAddr = "localhost:50051"
	}

	var client *agentclient.Client

	c, err := agentclient.New(serverAddr)
	if err != nil {
		log.Println("[agent] running in standalone mode (no server)")
	} else {
		client = c
		defer client.Close()
		log.Println("[agent] connected to server:", serverAddr)
	}

	/* ---------------- METRIC HANDLER ---------------- */

	m.OnMetric = func(metric *agentpb.NetworkMetric) {
		log.Printf(
			"[agent] metric: signal=%d ping=%d",
			metric.SignalPercent,
			metric.AvgPingMs,
		)


		if client == nil {
			return
		}

		if err := client.ReportMetric(metric); err != nil {
			log.Println("[agent] failed to report metric:", err)
		}
	}

	go startLocalAPI(m)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := m.Start(ctx); err != nil && err != context.Canceled {
		log.Println("[agent] monitor stopped with error:", err)
	}
}

func startLocalAPI(m *monitor.Monitor) {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	mux.HandleFunc("/modeSwitch", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		switchMode := m.SwitchAutomatically
		m.SwitchAutomatically = !switchMode
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]bool{"SwitchAutomatically": m.SwitchAutomatically})
	})
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

	addr := ":9090"
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Println("[agent] local api error:", err)
	}
}
