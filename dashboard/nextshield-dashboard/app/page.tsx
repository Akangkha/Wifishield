"use client";
import React from "react";
import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { supabase } from "../lib/client";

export default function HomePage() {
  const [checking, setChecking] = useState(true);
  const [authed, setAuthed] = useState(false);
  const router = useRouter();
  useEffect(() => {
    const run = async () => {
      const { data } = await supabase.auth.getUser();
      console.log(data.user);
      console.log(!!data.user);
      setAuthed(!!data.user);
      setChecking(false);
    };
    run();
  }, []);

  const onPrimary = () => {
    if (authed) router.push("/admin");
    else router.push("/login");
  };

  return (
    <main className="min-h-screen bg-[#020617] text-slate-100 flex items-center justify-center px-6">
      <div className="w-full max-w-5xl grid grid-cols-1 md:grid-cols-[2fr_1.3fr] gap-8">
        <section className="space-y-4">
          <h1 className="text-4xl md:text-5xl font-bold leading-tight bg-gradient-to-r from-emerald-400 via-cyan-400 to-sky-500 bg-clip-text text-transparent">
            NetShield
          </h1>
          <p className="text-slate-400 text-sm md:text-base max-w-xl">
            Real-time network experience monitoring and intelligent Wi-Fi
            failover for exams, remote work and telemedicine. Desktop agents
            stream live signal, latency and experience scores into a central
            analytics console.
          </p>
          <div className="flex flex-wrap gap-3 text-[11px]">
            <span className="px-3 py-1 rounded-full border border-emerald-500/60 bg-emerald-500/10 text-emerald-200">
              Go · gRPC · PostgreSQL
            </span>
            <span className="px-3 py-1 rounded-full border border-cyan-500/60 bg-cyan-500/10 text-cyan-200">
              Next.js · Tailwind · Electron
            </span>
            <span className="px-3 py-1 rounded-full border border-fuchsia-500/60 bg-fuchsia-500/10 text-fuchsia-200">
              Supabase Auth
            </span>
          </div>
          <div className="flex flex-wrap gap-3 pt-2">
            <button
              onClick={onPrimary}
              disabled={checking}
              className="px-5 py-2 rounded-md text-xs font-semibold bg-emerald-500 text-slate-900 hover:bg-emerald-400 disabled:opacity-60"
            >
              {checking
                ? "Checking access…"
                : authed
                ? "Open Admin Console"
                : "Sign in to Admin Console"}
            </button>
            <a
              href="/admin"
              className="px-4 py-2 rounded-md text-xs border border-slate-700 bg-slate-900/60 hover:bg-slate-800"
            >
              View dashboard
            </a>
          </div>
        </section>

        <section className="rounded-2xl border border-slate-800 bg-slate-900/80 p-5 shadow-[0_0_22px_#00ffbf25] flex flex-col justify-between">
          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <p className="text-xs text-slate-400">Access status</p>
              <p className="text-[11px] text-slate-500">
                {new Date().toLocaleTimeString()}
              </p>
            </div>
            <div className="mt-1">
              <p className="text-sm font-semibold">
                {checking
                  ? "Verifying session…"
                  : authed
                  ? "Admin access granted"
                  : "Restricted – authentication required"}
              </p>
              <p className="text-[11px] text-slate-400 mt-1">
                {authed
                  ? "You are signed in and can access the NetShield Network Intelligence Console."
                  : "Sign in to view live device telemetry and network health scores in the admin dashboard."}
              </p>
            </div>
          </div>

          <div className="mt-6 grid grid-cols-2 gap-3 text-[11px]">
            <div className="rounded-xl border border-slate-800 bg-slate-950/60 px-3 py-2">
              <p className="text-slate-400">Current mode</p>
              <p className="text-slate-100 font-medium mt-1">
                Network Intelligence
              </p>
            </div>
            <div className="rounded-xl border border-slate-800 bg-slate-950/60 px-3 py-2">
              <p className="text-slate-400">Agents</p>
              <p className="text-emerald-300 font-semibold mt-1">
                Live · Simulated
              </p>
            </div>
            <div className="rounded-xl border border-slate-800 bg-slate-950/60 px-3 py-2">
              <p className="text-slate-400">Creator</p>
              <a
                className="text-cyan-300 font-medium mt-1"
                href="https://www.linkedin.com/in/akangkha-sarkar"
                target="_blank"
                rel="noopener noreferrer"
              >
                @ Akangkha
              </a>
            </div>
            <div className="rounded-xl border border-slate-800 bg-slate-950/60 px-3 py-2">
              <p className="text-slate-400">Github</p>
              <a
                className="text-slate-100 font-medium mt-1"
                href="https://github.com/akangkha"
                target="_blank"
                rel="noopener noreferrer"
              >
                /Akangkha
              </a>
            </div>
          </div>
        </section>
      </div>
    </main>
  );
}
