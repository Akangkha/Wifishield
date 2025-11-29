package wifi

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)


type WifiProfile struct {
	RawName   string 
	CleanName string 
}

// WifiStatus represents the current connection info.
type WifiStatus struct {
	InterfaceName string
	SSID          string
	ProfileName   string
	Signal        int 
}


type Manager interface {
	ListProfiles() ([]WifiProfile, error)
	GetCurrentStatus() (*WifiStatus, error)
	Connect(profile WifiProfile) error
}


type WindowsManager struct{}

func (w WindowsManager) runNetsh(args ...string) (string, error) {
	cmd := exec.Command("netsh", args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("netsh %v failed: %v | stderr: %s", args, err, stderr.String())
	}
	return out.String(), nil
}


func (w WindowsManager) ListProfiles() ([]WifiProfile, error) {
	out, err := w.runNetsh("wlan", "show", "profiles")
	if err != nil {
		return nil, err
	}
	return ParseProfiles(out), nil
}


func (w WindowsManager) GetCurrentStatus() (*WifiStatus, error) {
	out, err := w.runNetsh("wlan", "show", "interfaces")
	if err != nil {
		return nil, err
	}
	status := ParseCurrentStatus(out)
	if status == nil {
		return nil, fmt.Errorf("no active Wi-Fi interface found")
	}
	return status, nil
}


func (w WindowsManager) Connect(profile WifiProfile) error {
	
	args := []string{"wlan", "connect", "name=" + profile.RawName}
	_, err := w.runNetsh(args...)
	return err
}


func FindProfileByCleanName(profiles []WifiProfile, name string) *WifiProfile {
	for _, p := range profiles {
		if strings.TrimSpace(p.CleanName) == strings.TrimSpace(name) {
			cp := p
			return &cp
		}
	}
	return nil
}
