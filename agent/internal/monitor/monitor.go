package monitor

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"netshield/agent/internal/wifi"
	agentpb "netshield/agent/proto"
	"os/exec"
	"sync"
	"time"
)

type Config struct {
	MinSignalPercent  int
	MaxAvgPingMs      int
	PingHost          string
	CheckInterval     time.Duration
	PreferredProfiles []string
}

type Snapshot struct {
	SSID        string    `json:"ssid"`
	Profile     string    `json:"profile"`
	Signal      int       `json:"signal_percent"`
	AvgPingMs   int       `json:"avg_ping_ms"`
	Score       int       `json:"score"`
	LastUpdated time.Time `json:"last_updated"`
}

type Monitor struct {
	Wifi   wifi.Manager
	Config Config

	mu       sync.RWMutex
	snapshot Snapshot
	OnMetric func(*agentpb.NetworkMetric)
}

func (m *Monitor) Start(ctx context.Context) error {
	ticker := time.NewTicker(m.Config.CheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := m.checkOnce(); err != nil {
				fmt.Println("[monitor] error:", err)
			}
		}
	}
}

func (m *Monitor) GetSnapshot() Snapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.snapshot
}

func (m *Monitor) checkOnce() error {
	status, err := m.Wifi.GetCurrentStatus()
	if err != nil {
		return fmt.Errorf("get current status: %w", err)
	}

	pingRes, err := pingHost(m.Config.PingHost)
	if err != nil {
		fmt.Println("[monitor] ping failed:", err)
	}

	var avgPing int
	if pingRes != nil {
		avgPing = pingRes.AvgMs
	}

	score := computeScore(status.Signal, avgPing)

	m.mu.Lock()
	m.snapshot = Snapshot{
		SSID:        status.SSID,
		Profile:     status.ProfileName,
		Signal:      status.Signal,
		AvgPingMs:   avgPing,
		Score:       score,
		LastUpdated: time.Now(),
	}

	m.mu.Unlock()
	log.Print("profile:", m.snapshot.Profile)
	if m.OnMetric != nil {
		metric := &agentpb.NetworkMetric{
			DeviceId:        status.SSID,        // or hostname / generated ID
			UserId:          status.ProfileName, // optional
			Domain:          "laptop",           // optional
			TimestampUnix:   time.Now().Unix(),
			Ssid:            status.SSID,
			InterfaceName:   status.InterfaceName, // if available
			SignalPercent:   int32(status.Signal),
			AvgPingMs:       int32(avgPing),
			ExperienceScore: int32(score),
		}

		m.OnMetric(metric)
	} else {
		log.Println("[monitor] no OnMetric handler set")
	}
	badSignal := status.Signal > 0 && status.Signal < m.Config.MinSignalPercent
	badPing := avgPing > 0 && avgPing > m.Config.MaxAvgPingMs

	if !badSignal && !badPing {
		return nil
	}

	return m.tryFailover(status)
}

func computeScore(signal, ping int) int {
	if signal <= 0 {
		return 0
	}
	if ping <= 0 {
		return signal
	}
	penalty := ping / 5
	if penalty > 40 {
		penalty = 40
	}
	score := signal - penalty
	if score < 0 {
		return 0
	}
	if score > 100 {
		return 100
	}
	return score
}
func (m *Monitor) tryFailover(current *wifi.WifiStatus) error {
	profiles, err := m.Wifi.ListProfiles()
	if err != nil {
		return fmt.Errorf("list profiles: %w", err)
	}

	for _, preferredName := range m.Config.PreferredProfiles {

		if preferredName == current.ProfileName || preferredName == current.SSID {
			continue
		}

		p := wifi.FindProfileByCleanName(profiles, preferredName)
		if p == nil {
			continue
		}

		fmt.Println("[monitor] attempting switch to:", p.CleanName)
		if err := m.Wifi.Connect(*p); err != nil {
			fmt.Println("[monitor] connect failed:", err)
			continue
		}

		time.Sleep(7 * time.Second)

		newStatus, err := m.Wifi.GetCurrentStatus()
		if err != nil {
			fmt.Println("[monitor] after-switch status error:", err)
			continue
		}

		fmt.Println("[monitor] after-switch:", wifi.DebugStatus(newStatus))

		if newStatus.ProfileName == p.RawName || newStatus.SSID == p.CleanName {
			if newStatus.Signal >= m.Config.MinSignalPercent {
				fmt.Println("[monitor] failover successful 🎉")
				return nil
			}
		}
	}

	return fmt.Errorf("no suitable alternative profile found or all failed")
}
func pingHost(host string) (*wifi.SimplePingResult, error) {
	cmd := exec.Command("ping", "-n", "3", host)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ping failed: %v | stderr: %s", err, stderr.String())
	}

	res := wifi.ParsePingOutput(out.String())
	if res == nil {
		return nil, fmt.Errorf("could not parse ping output")
	}
	return res, nil
}
