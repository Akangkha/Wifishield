package wifi

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)


func ParseProfiles(output string) []WifiProfile {
	var profiles []WifiProfile
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		
		if strings.Contains(line, "All User Profile") {
		
			line = strings.TrimRight(line, "\r\n")
			parts := strings.SplitN(line, ":", 2)
			if len(parts) != 2 {
				continue
			}

	
			raw := parts[1]
			if len(raw) > 0 && raw[0] == ' ' {
				raw = raw[1:]
			}

			profiles = append(profiles, WifiProfile{
				RawName:   raw,                 
				CleanName: strings.TrimSpace(raw),
			})
		}
	}
	return profiles
}


func ParseCurrentStatus(output string) *WifiStatus {
	lines := strings.Split(output, "\n")
	status := &WifiStatus{}

	for _, line := range lines {
		line = strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(line, "Name") && status.InterfaceName == "":
	
			status.InterfaceName = afterColon(line)
		case strings.HasPrefix(line, "SSID") && !strings.Contains(line, "BSSID"):
			status.SSID = afterColon(line)
		case strings.HasPrefix(line, "Profile"):
			status.ProfileName = afterColon(line)
		case strings.HasPrefix(line, "Signal"):
			
			raw := afterColon(line)
			raw = strings.TrimSpace(strings.TrimSuffix(raw, "%"))
			if v, err := strconv.Atoi(raw); err == nil {
				status.Signal = v
			}
		}
	}

	if status.InterfaceName == "" || status.SSID == "" {
		return nil
	}
	return status
}


func afterColon(line string) string {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

func DebugStatus(s *WifiStatus) string {
	if s == nil {
		return "no wifi status"
	}
	return fmt.Sprintf("Interface=%s, SSID=%s, Profile=%s, Signal=%d%%",
		s.InterfaceName, s.SSID, s.ProfileName, s.Signal)
}


type SimplePingResult struct {
	AvgMs int
}


func ParsePingOutput(out string) *SimplePingResult {
	re := regexp.MustCompile(`Average\s*=\s*(\d+)ms`)
	m := re.FindStringSubmatch(out)
	if len(m) != 2 {
		return nil
	}
	v, err := strconv.Atoi(m[1])
	if err != nil {
		return nil
	}
	return &SimplePingResult{AvgMs: v}
}
