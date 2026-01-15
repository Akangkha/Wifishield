package db

import (
	"context"
	"time"

	agentpb "netshield/agent/proto"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	Pool *pgxpool.Pool
}

func New(ctx context.Context, dsn string) (*Store, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return &Store{Pool: pool}, nil
}

func (s *Store) Close() {
	s.Pool.Close()
}

// SaveMetric inserts into metrics_raw and upserts into device_status.
func (s *Store) SaveMetric(ctx context.Context, m *agentpb.NetworkMetric) error {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	ts := time.Unix(m.TimestampUnix, 0)

	// Insert raw metric
	_, err = tx.Exec(ctx, `
		INSERT INTO metrics_raw (
			device_id, user_id, domain, ts,
			ssid, interface_name,
			signal_percent, avg_ping_ms, experience_score
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`,
		m.DeviceId, m.UserId, m.Domain, ts,
		m.Ssid, m.InterfaceName,
		m.SignalPercent, m.AvgPingMs, m.ExperienceScore,
	)
	if err != nil {
		return err
	}

	// Upsert device_status
	_, err = tx.Exec(ctx, `
		INSERT INTO device_status (
			device_id, user_id, domain, last_seen,
			ssid, interface_name, signal_percent, avg_ping_ms, experience_score
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		ON CONFLICT (device_id) DO UPDATE
		SET
			user_id          = EXCLUDED.user_id,
			domain           = EXCLUDED.domain,
			last_seen        = EXCLUDED.last_seen,
			ssid             = EXCLUDED.ssid,
			interface_name   = EXCLUDED.interface_name,
			signal_percent   = EXCLUDED.signal_percent,
			avg_ping_ms      = EXCLUDED.avg_ping_ms,
			experience_score = EXCLUDED.experience_score
	`,
		m.DeviceId, m.UserId, m.Domain, ts,
		m.Ssid, m.InterfaceName,
		m.SignalPercent, m.AvgPingMs, m.ExperienceScore,
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// DeviceStatusRow maps to JSON for /status.
type DeviceStatusRow struct {
	DeviceID        string    `json:"device_id"`
	UserID          string    `json:"user_id"`
	Domain          string    `json:"domain"`
	LastSeen        time.Time `json:"last_seen"`
	SSID            string    `json:"ssid"`
	InterfaceName   string    `json:"interface_name"`
	SignalPercent   int32     `json:"signal_percent"`
	AvgPingMs       int32     `json:"avg_ping_ms"`
	ExperienceScore int32     `json:"experience_score"`
}

// GetAllDeviceStatus returns one row per device.
func (s *Store) GetAllDeviceStatus(ctx context.Context) ([]DeviceStatusRow, error) {
	rows, err := s.Pool.Query(ctx, `
		SELECT device_id, user_id, domain, last_seen,
		       ssid, interface_name, signal_percent, avg_ping_ms, experience_score
		FROM device_status
		ORDER BY last_seen DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []DeviceStatusRow
	for rows.Next() {
		var r DeviceStatusRow
		if err := rows.Scan(
			&r.DeviceID, &r.UserID, &r.Domain, &r.LastSeen,
			&r.SSID, &r.InterfaceName, &r.SignalPercent, &r.AvgPingMs, &r.ExperienceScore,
		); err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	return result, rows.Err()
}

// DeleteOldMetrics deletes raw metrics older than the given retention.
func (s *Store) DeleteOldMetrics(ctx context.Context, olderThan time.Duration) error {
	cutoff := time.Now().Add(-olderThan)
	_, err := s.Pool.Exec(ctx, `DELETE FROM metrics_raw WHERE ts < $1`, cutoff)
	return err
}

func (s *Store) GetDevicesByLink(
	ctx context.Context,
	ssid string, //network identifier
	activeWithin time.Duration,
) ([]DeviceStatusRow, error) {

	rows, err := s.Pool.Query(ctx, `
		SELECT
			device_id,
			user_id,
			domain,
			ssid,
			signal_percent,
			avg_ping_ms,
			experience_score,
			last_seen
		FROM device_status
		WHERE ssid = $1
		ORDER BY last_seen DESC
	`, ssid);

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []DeviceStatusRow
	for rows.Next() {
		var d DeviceStatusRow
		if err := rows.Scan(
			&d.DeviceID,
			&d.UserID,
			&d.Domain,
			&d.SSID,
			&d.SignalPercent,
			&d.AvgPingMs,
			&d.ExperienceScore,
			&d.LastSeen,
		); err != nil {
			return nil, err
		}
		devices = append(devices, d)
	}

	return devices, nil
}
