export type DeviceStatus = {
  device_id: string;
  user_id: string;
  domain: string;
  last_seen: string;
  ssid: string;
  interface_name: string;
  signal_percent: number;
  avg_ping_ms: number;
  experience_score: number;
};

export async function fetchStatus(): Promise<DeviceStatus[]> {
  const res = await fetch("/api/status", {
    cache: "no-store",
  });

  if (!res.ok) {
    throw new Error(`Status fetch failed: ${res.status}`);
  }
  return res.json();
}
