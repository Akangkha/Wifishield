package monitor

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"netshield/agent/internal/wifi"
	agentpb "netshield/agent/proto"
	"os/exec"
	"strings"
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
	Wifi                wifi.Manager
	Config              Config
	SwitchAutomatically bool
	mu                  sync.RWMutex
	snapshot            Snapshot
	OnMetric            func(*agentpb.NetworkMetric)
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

	if m.SwitchAutomatically {
		return m.tryConnectVisibleNetworks()
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

func extractSSIDs(output string) []wifi.WifiProfile {
	var ssids []wifi.WifiProfile

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Match: SSID X : NAME
		if strings.HasPrefix(line, "SSID ") && strings.Contains(line, " : ") {
			parts := strings.SplitN(line, " : ", 2)
			if len(parts) == 2 {
				name := strings.TrimSpace(parts[1])
				if name != "" {
					ssids = append(ssids, wifi.WifiProfile{
						CleanName: name,
						RawName:   name,
					})
				}
			}
		}
	}

	return ssids
}

func (m *Monitor) tryConnectVisibleNetworks() error {
	output, err := exec.Command(
		"netsh", "wlan", "show", "networks", "mode=bssid",
	).Output()
	if err != nil {
		return err
	}

	visible := extractSSIDs(string(output))

	savedProfiles, err := m.Wifi.ListProfiles()
	if err != nil {
		return err
	}

	// Build lookup for saved profiles
	savedByName := make(map[string]wifi.WifiProfile)
	for _, p := range savedProfiles {
		savedByName[p.CleanName] = p
	}

	log.Println("[agent] visible SSIDs:")
	for _, v := range visible {
		log.Println(" -", v.CleanName)
	}
	for _, v := range visible {
		// Check if visible SSID has a saved profile
		p, ok := savedByName[v.CleanName]
		if !ok {
			log.Println("[agent] visible SSID has no saved profile:", v.CleanName)
			continue // visible but no credentials
		}

		// Skip current connection
		if p.CleanName == m.snapshot.Profile {
			log.Println("[agent] already connected to:", p.RawName)
			continue
		}

		log.Println("[agent] attempting connect to:", p.RawName)

		if err := m.Wifi.Connect(p); err != nil {
			log.Println("[monitor] connect failed:", err)
			continue
		}

		log.Println("[agent] connected to:", p.RawName)
		return nil
	}

	return nil
}
