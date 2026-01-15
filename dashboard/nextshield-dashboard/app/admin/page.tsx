"use client";

import React, { useCallback, useEffect, useMemo, useState, memo } from "react";
import { supabase } from "../../lib/client";
import { DomainFilter } from "../../components/DomainFilter";
import { DataTable } from "../../components/DataTable";
import { DeviceStatus } from "../../lib/api";

export default function Admin() {
  const [filter, setFilter] = useState("all");
  const [query, setQuery] = useState("");
  const [checking, setChecking] = useState(true);
  const [authed, setAuthed] = useState(false);
  const [devices, setDevices] = useState<DeviceStatus[]>([]);
  const [networkId, setNetworkId] = useState<string | null>(null);
  const [time, setTime] = useState("");

  /* ---------------- AUTH CHECK ---------------- */
  useEffect(() => {
    supabase.auth.getUser().then(({ data }) => {
      setAuthed(!!data.user);
      setChecking(false);
    });
  }, []);

  /* ---------------- CLOCK (isolated) ---------------- */
  useEffect(() => {
    const update = () => setTime(new Date().toLocaleTimeString());
    update();
    const id = setInterval(update, 1000);
    return () => clearInterval(id);
  }, []);

  /* ---------------- FETCH DEVICES ---------------- */
  const fetchDevices = useCallback(async (ssid: string) => {
    const res = await fetch(`/api/getDevices?ssid=${ssid}`);
    const json = await res.json();
    setDevices(json.devices || []);
  }, []);

  useEffect(() => {
    if (networkId) fetchDevices(networkId);
  }, [networkId, fetchDevices]);

  /* ---------------- FILTER DATA ---------------- */
  const filteredData = useMemo(() => {
    return devices
      .filter((d) =>
        filter === "all" ? true : d.domain?.toLowerCase() === filter
      )
      .filter((d) => {
        if (!query) return true;
        const q = query.toLowerCase();
        return (
          d.device_id?.toLowerCase().includes(q) ||
          d.user_id?.toLowerCase().includes(q) ||
          d.ssid?.toLowerCase().includes(q)
        );
      });
  }, [devices, filter, query]);

  /* ---------------- STABLE HANDLERS ---------------- */
  const handleSignOut = useCallback(async () => {
    await supabase.auth.signOut();
    setAuthed(false);
  }, []);

  const clearNetwork = useCallback(() => {
    setNetworkId(null);
    setDevices([]);
  }, []);

  /* ---------------- GUARDS ---------------- */
  if (checking) return <Loading text="Checking access…" />;
  if (!authed) return <Unauthorized />;

  /* ---------------- UI ---------------- */
  return (
    <main className="min-h-screen bg-[#020617] text-slate-100 px-6 py-10">
      <div className="max-w-7xl mx-auto space-y-6">
        <Header time={time} onSignOut={handleSignOut} />

        <div className="flex items-center gap-3">
          <NetworkIdDialog onSubmit={setNetworkId} />
          {networkId && (
            <NetworkBanner networkId={networkId} onClear={clearNetwork} />
          )}
        </div>

        <div className="flex justify-between items-center gap-3">
          <DomainFilter value={filter} onChange={setFilter} />
          <input
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Search device, user or SSID"
            className="w-72 rounded-md border border-slate-700 bg-slate-900 px-3 py-1 text-xs text-slate-100 outline-none focus:ring-2 focus:ring-emerald-500/50"
          />
        </div>

        <DataTable rows={filteredData} />
      </div>
    </main>
  );
}

/* ------------------------------------------------------------------ */
/* -------------------------- COMPONENTS ------------------------------ */
/* ------------------------------------------------------------------ */

const Header = memo(function Header({
  time,
  onSignOut,
}: {
  time: string;
  onSignOut: () => void;
}) {
  return (
    <header className="flex justify-between items-center">
      <div>
        <h1 className="text-3xl font-bold bg-gradient-to-r from-emerald-400 to-cyan-400 bg-clip-text text-transparent">
          WifiShield Network Intelligence Console
        </h1>
        <p className="text-slate-400 text-sm">
          Real-time telemetry across connected devices
        </p>
      </div>

      <div className="text-right text-xs text-slate-400 space-y-1">
        <p>{time}</p>
        <button
          onClick={onSignOut}
          className="px-3 py-1 rounded-md border border-slate-700 hover:bg-slate-800"
        >
          Sign out
        </button>
      </div>
    </header>
  );
});

const Loading = ({ text }: { text: string }) => (
  <main className="min-h-screen bg-[#020617] flex items-center justify-center text-slate-300 text-sm">
    {text}
  </main>
);

const Unauthorized = () => (
  <main className="min-h-screen bg-[#020617] flex items-center justify-center text-slate-300 text-sm">
    <div className="text-center space-y-2">
      <p>Unauthorized. Please sign in.</p>
      <a
        href="/login"
        className="inline-flex px-3 py-1 rounded-md bg-emerald-500 text-slate-900 text-xs font-semibold"
      >
        Go to login
      </a>
    </div>
  </main>
);

/* ---------------- NETWORK DIALOG ---------------- */

const NetworkIdDialog = memo(function NetworkIdDialog({
  onSubmit,
}: {
  onSubmit: (id: string) => void;
}) {
  const [open, setOpen] = useState(false);
  const [value, setValue] = useState("");

  const submit = () => {
    if (!value.trim()) return;
    onSubmit(value.trim());
    setOpen(false);
  };

  return (
    <>
      <button
        onClick={() => setOpen(true)}
        className="rounded-md border border-slate-700 bg-slate-900 px-3 py-1 text-xs hover:border-emerald-500"
      >
        + Set Network ID
      </button>

      {open && (
        <div className="fixed inset-0 z-50 bg-black/60 flex items-center justify-center">
          <div className="w-96 rounded-lg border border-slate-700 bg-slate-900 p-4">
            <h3 className="mb-2 text-sm font-semibold">Enter Network ID</h3>

            <input
              value={value}
              onChange={(e) => setValue(e.target.value)}
              placeholder="e.g. esperance"
              className="w-full rounded-md border border-slate-700 bg-slate-900 px-3 py-2 text-xs outline-none focus:ring-2 focus:ring-emerald-500/50"
            />

            <div className="mt-4 flex justify-end gap-2">
              <button
                onClick={() => setOpen(false)}
                className="text-xs text-slate-400"
              >
                Cancel
              </button>
              <button
                onClick={submit}
                className="rounded-md bg-emerald-600 px-3 py-1 text-xs text-black"
              >
                Send
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  );
});

/* ---------------- NETWORK BANNER ---------------- */

const NetworkBanner = memo(function NetworkBanner({
  networkId,
  onClear,
}: {
  networkId: string;
  onClear: () => void;
}) {
  return (
    <div className="flex items-center gap-2 rounded-md border border-emerald-600/40 bg-emerald-500/10 px-3 py-1 text-xs text-emerald-400">
      <span>
        Network: <b>{networkId}</b>
      </span>
      <button onClick={onClear} className="hover:text-red-400">
        ✕
      </button>
    </div>
  );
});


