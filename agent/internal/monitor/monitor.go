package monitor

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
	"netshield/agent/internal/wifi"
)

// policy configuration for the monitor.
type Config struct {
	MinSignalPercent int           
	MaxAvgPingMs     int           
	PingHost         string       
	CheckInterval    time.Duration 
	PreferredProfiles []string
	OnMetric func(status *wifi.WifiStatus, ping *wifi.SimplePingResult)
}


type Monitor struct {
	Wifi   wifi.Manager
	Config Config
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


func (m *Monitor) checkOnce() error {
	status, err := m.Wifi.GetCurrentStatus()
	if err != nil {
		return fmt.Errorf("get current status: %w", err)
	}
	fmt.Println("[monitor] current:", wifi.DebugStatus(status))
	pingRes, err := pingHost(m.Config.PingHost)
	if err != nil {
		fmt.Println("[monitor] ping failed:", err)
	} else {
		fmt.Printf("[monitor] ping avg: %dms\n", pingRes.AvgMs)
	}

	if m.Config.OnMetric != nil && pingRes != nil {
		m.Config.OnMetric(status, pingRes)
	}

	
	badSignal := status.Signal > 0 && status.Signal < m.Config.MinSignalPercent
	badPing := pingRes != nil && pingRes.AvgMs > m.Config.MaxAvgPingMs

	if !badSignal && !badPing {
		fmt.Println("[monitor] connection healthy, no switch")
		return nil
	}

	fmt.Println("[monitor] connection poor, trying failover...")

	return m.tryFailover(status)
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
