package main

import (
	"context"
	"fmt"
	agentclient "netshield/agent/internal/client"
	"netshield/agent/internal/metrics"
	"netshield/agent/internal/monitor"
	"netshield/agent/internal/wifi"
	agentpb "netshield/agent/proto"
	"os"
	"os/signal"
	"time"
)

func main() {
	wm := wifi.WindowsManager{}
	fmt.Println("Available Wi-Fi profiles:")
	profiles, err := wm.ListProfiles()
	if err != nil {
		fmt.Println("Error listing profiles:", err)
	} else {
		for _, p := range profiles {
			fmt.Printf(" - Raw: %q | Clean: %q\n", p.RawName, p.CleanName)
		}
	}

	serverAddr := "localhost:50051"
	agentClient, err := agentclient.New(serverAddr)
	if err != nil {
		fmt.Println("Failed to connect to gRPC server:", err)
		agentClient = nil
	} else {
		defer agentClient.Close()
	}

	deviceID := "device-1"  // config / hardware ID
	userID := "user-1"      // from config / login
	domain := "remote-work" // "exam" / "telemedicine"

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

		OnMetric: func(status *wifi.WifiStatus, ping *wifi.SimplePingResult) {
			if agentClient == nil || ping == nil {
				return
			}
			score := metrics.Score(status.Signal, ping.AvgMs)

			m := &agentpb.NetworkMetric{
				DeviceId:        deviceID,
				UserId:          userID,
				Domain:          domain,
				TimestampUnix:   time.Now().Unix(),
				Ssid:            status.SSID,
				InterfaceName:   status.InterfaceName,
				SignalPercent:   int32(status.Signal),
				AvgPingMs:       int32(ping.AvgMs),
				ExperienceScore: int32(score),
			}

			if err := agentClient.ReportMetric(m); err != nil {
				fmt.Println("[agent] failed to report metric:", err)
			} else {
				fmt.Printf("[agent] reported metric: score=%d\n", score)
			}
		},
	}

	m := &monitor.Monitor{
		Wifi:   wm,
		Config: cfg,
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	fmt.Println("Starting Wi-Fi monitor… (Ctrl+C to exit)")
	if err := m.Start(ctx); err != nil && err != context.Canceled {
		fmt.Println("Monitor stopped with error:", err)
	} else {
		fmt.Println("Monitor stopped.")
	}
}
