"use client";

import { useEffect, useState } from "react";
import { DeviceStatus, fetchStatus } from "@/lib/api";

type UiState = {
  ssid: string;
  signalPercent: number;
  avgPingMs: number;
  score: number;
};

function scoreLabel(score: number) {
  if (score >= 80) return "Excellent";
  if (score >= 60) return "Good";
  if (score >= 40) return "Fair";
  return "Poor";
}

export default function DashboardWidget() {
  const [data, setData] = useState<UiState>({
    ssid: "Waiting for agent…",
    signalPercent: 0,
    avgPingMs: 0,
    score: 0,
  });
  const [autoSwitch, setAutoSwitch] = useState(true);
  const [error, setError] = useState<string | null>(null);


  useEffect(() => {
    let cancelled = false;

    async function load() {
      try {
        console.log("Fetching status...");
        const res = await fetchStatus();
        if (!res || res.length === 0) {
          if (!cancelled) {
            setError("No device reporting yet");
          }
          return;
        }
        const d = res[0]; 
        if (!cancelled) {
          setError(null);
          setData({
            ssid: d.ssid || "(unknown)",
            signalPercent: d.signal_percent ?? 0,
            avgPingMs: d.avg_ping_ms ?? 0,
            score: d.experience_score ?? 0,
          });
        }
      } catch (e: any) {
        if (!cancelled) {
          setError(e.message ?? "Failed to fetch status");
        }
      }
    }

    load();
    const id = setInterval(load, 5000);
    return () => {
      cancelled = true;
      clearInterval(id);
    };
  }, []);

  const angle = (data.signalPercent / 100) * 180 - 90; // 0..100 => -90..+90 deg

  return (
    <main className="min-h-screen bg-slate-950 flex items-center justify-center">
      <div className="w-[400px] rounded-2xl border border-slate-800 bg-slate-900/90 shadow-2xl p-4">
        <header className="flex items-center justify-between mb-3">
          <div>
            <h1 className="text-sm font-semibold text-slate-100">
              NetShield Client
            </h1>
            <p className="text-xs text-slate-400">
              Local Wi-Fi health & auto-switch
            </p>
          </div>
       
          <div className="flex gap-1">
            <span className="h-2.5 w-2.5 rounded-full bg-red-500/80" />
            <span className="h-2.5 w-2.5 rounded-full bg-yellow-400/80" />
            <span className="h-2.5 w-2.5 rounded-full bg-green-500/80" />
          </div>
        </header>

        <div className="flex gap-4">
     
          <div className="flex-1 flex flex-col items-center justify-center">
            <div className="relative w-32 h-16 overflow-hidden">
          
              <div className="absolute inset-0 rounded-t-full border-4 border-slate-800 border-b-0" />
        
              <div className="absolute inset-0 rounded-t-full border-4 border-gradient border-b-0 [border-image:linear-gradient(90deg,#22c55e,#eab308,#ef4444)_1]" />
     
              <div className="absolute bottom-0 left-1/2 h-[2px] w-[2px] -translate-x-1/2 translate-y-1/2 rounded-full bg-slate-300" />
            
              <div
                className="absolute bottom-0 left-1/2 origin-bottom h-14 w-[2px] -translate-x-1/2 bg-slate-100"
                style={{ transform: `rotate(${angle}deg)` }}
              />
            </div>
            <div className="mt-2 text-xs text-slate-400">
              Signal strength:{" "}
              <span className="text-slate-100 font-medium">
                {data.signalPercent}%
              </span>
            </div>
          </div>

       
          <div className="flex-1 flex flex-col justify-between">
            <div>
              <div className="text-xs text-slate-400 mb-1">SSID</div>
              <div className="text-sm font-medium text-slate-100 truncate">
                {data.ssid}
              </div>

              <div className="mt-3 grid grid-cols-2 gap-2 text-xs">
                <div className="rounded-lg bg-slate-800/70 px-2 py-1.5">
                  <div className="text-slate-400">Ping</div>
                  <div className="text-slate-100 font-medium">
                    {data.avgPingMs} ms
                  </div>
                </div>
                <div className="rounded-lg bg-slate-800/70 px-2 py-1.5">
                  <div className="text-slate-400">Score</div>
                  <div className="text-slate-100 font-medium">
                    {data.score} / 100
                  </div>
                </div>
              </div>

              <div className="mt-2 text-xs text-slate-400">
                Status:{" "}
                <span className="font-medium text-slate-100">
                  {scoreLabel(data.score)}
                </span>
              </div>
            </div>

    
            <div className="mt-3 flex items-center justify-between rounded-xl bg-slate-800/70 px-3 py-2">
              <div>
                <div className="text-xs font-medium text-slate-100">
                  Auto-switch
                </div>
                <div className="text-[11px] text-slate-400">
                  Automatically hop to better Wi-Fi
                </div>
              </div>
              <button
                type="button"
                onClick={() => setAutoSwitch((v) => !v)}
                className={`relative inline-flex h-6 w-11 items-center rounded-full transition ${
                  autoSwitch ? "bg-emerald-500" : "bg-slate-600"
                }`}
              >
                <span
                  className={`inline-block h-5 w-5 transform rounded-full bg-white shadow transition ${
                    autoSwitch ? "translate-x-5" : "translate-x-1"
                  }`}
                />
              </button>
            </div>
          </div>
        </div>

        {error && (
          <div className="mt-3 rounded-md border border-yellow-500/40 bg-yellow-900/30 px-2 py-1.5 text-[11px] text-yellow-100">
            {error}
          </div>
        )}
      </div>
    </main>
  );
}
