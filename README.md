<div align="left">
  <img src="https://github.com/user-attachments/assets/08af2f54-b9ce-4c0a-b701-c024ce031c95" width="80" height="80" align="left" style="margin-right: 20px" />
  <h1 style="padding-top: 20px;">
    WifiShield 
    <span style="font-size: 0.5em; font-weight: normal;">Intelligent Wi-Fi Failover & Real-Time Network Health</span>
  </h1>
</div>
<br clear="left"/>
 <img width="1918" height="855" alt="image" src="https://github.com/user-attachments/assets/33faa55e-034b-4c93-b2c6-fed78ffe837c" />


> **“Because losing your network mid-exam / interview / call should never be the reason you fail.”**

WifiShield is a **desktop network resilience system** that runs on your machine, continuously monitors Wi-Fi quality, and automatically switches to better networks when your current one degrades.

It ships as **three executables**:

- **WifiShield Agent** – lightweight Go-based background process (CLI).
- **WifiShield Desktop** – Electron-based GUI widget with system tray integration.
- **WifiShield Admin Console** – NextJs based admin console for the admin to monitor the network failover of the participants.

Together, they give you **live Wi-Fi health, a beautiful desktop widget, and automatic failover** – built with **Golang, Electron, Next.js, Supabase & PostgreSQL**.

---

## 📥 Downloads

End users don’t need Go, Node, or Docker.

1. Go to **GitHub → Releases** for this repo.
2. Download:
   - `WifiShield-Agent.exe` – CLI background agent
   - `Wifihield-Desktop-Setup.exe` – installer for the GUI widget
3. Install & run:
   - Start **WifiShield Agent** (runs in the background, console or service).
   - Install and launch **WifiShield Desktop** (system tray + widget window).

> 💡 The widget reads live metrics from the **local agent** via `http://127.0.0.1:9090/current`.  
> No internet / server required for the GUI to function.

---

## ✨ What You Get

- 🧠 **Go Agent**
  - Monitors current Wi-Fi SSID, signal strength, RSSI, latency (ping) & computes an **experience score**.
  - Talks to Windows via `netsh` and `ping`.
  - Exposes a local HTTP API: `GET http://127.0.0.1:9090/current`.
  - Designed to be lightweight & always running in the background.

- 📊 **Desktop Network Widget (Electron + React/Next.js)**
  - Clean, neon-styled **speedometer gauge** for Wi-Fi signal.
  - Shows SSID, ping, score, and connection status (“Excellent / Good / Fair / Poor”).
  - **Always-on-top** frameless window – feels like a modern desktop HUD.
  - **System tray icon** with custom app icon; can hide/show the widget.

- 🛰 **(Pluggable) Central Telemetry Pipeline**
  - Agent can stream metrics via **gRPC** to a central Go server.
  - Backend persists metrics in **PostgreSQL** (via Docker).
  - Designed for **multi-device / admin monitoring scenarios** (exam lab, campus, org).

- 🖥 **Admin / Web Console (Next.js)**
  - Landing page with **access status check** (Supabase Auth).
  - `/admin` dashboard UI with **dummy / seeded** device data for now:
    - Device ID, user, domain (exam / remote-work / telemedicine)
    - Signal %, ping, experience score, last seen.
  - Dark + neon cyber theme, designed as a **Network Intelligence Console**.
  - Ready to be wired to the real server `/status` endpoint.

- 🔐 **Authentication (Supabase)**
  - Supabase email/password login for the admin dashboard.
  - Simple login/signup page (`/login`) with access gating for `/admin`.

---

## 🧩 Architecture Overview

NetShield is intentionally split into **two processes** on the user’s machine:

````txt
+--------------------------------------------+
|       NetShield Desktop (Electron)         |
|  – React/Next.js widget                    |
|  – System tray, gauge UI                   |
|  – Fetches local status from agent         |
+------------------------▲-------------------+
                         │ HTTP (localhost)
                         │ GET /current
+------------------------┴-------------------+
|          NetShield Agent (Go)              |
|  – Runs in background (CLI/service)        |
|  – Uses netsh + ping to monitor Wi-Fi      |
|  – Computes score, auto-failover Wi-Fi     |
|  – Exposes /current JSON snapshot          |
+--------------------------------------------+
````

Optionally, the agent can stream telemetry to a central server:

```txt
    Agent(s)            Central Backend             UI
+-------------+      +---------------------+    +-------------------+
| Go Agent    | ---> | Go gRPC Server      | -> | Admin Dashboard   |
| (on laptop) |      | + PostgreSQL        |    | (Next.js Web App) |
+-------------+      +---------------------+    +-------------------+
```

---



## 🧪 Developer Setup

> This section is for developers, not end-users.

### Prereqs

* Go 1.21+
* Node.js 18+ & npm
* Docker (for server + DB, optional)
* Supabase project (for admin auth, optional)

### 1. Run the Agent from source

```bash
cd agent
go run ./cmd/shieldagent
```

It will:

* Start the monitor loop.
* Serve `http://127.0.0.1:9090/current`.

### 2. Run the Widget in Dev

```bash
cd client/widget
npm install
npm run dev:desktop
```

This runs:

* Next.js dev server
* Electron pointing at `http://localhost:3000`

### 3. Build Agent EXE (with icon)

```bash
cd agent
go build -o bin/NetShield-Agent.exe ./cmd/shieldagent
```

### 4. Build Desktop Installer

```bash
cd client/widget
npm run build:desktop
```

This produces:

```text
client/widget/dist/NetShield Setup.exe
```



## 🧑‍💻 About the Author

Built end-to-end by **Akangkha**:
<br/>
<br/>
<a href="https://www.linkedin.com/in/akangkha-sarkar">
    <img src="https://img.shields.io/badge/LinkedIn-Connect-0077B5?style=for-the-badge&logo=linkedin&logoColor=white" alt="LinkedIn" />
 </a>

* Full-stack engineer with strong systems, cloud & networking background.
* Implemented:

  * Golang Wi-Fi agent
  * Failover logic
  * grpc API backend
  * Electron widget UI
  * Next.js admin console
  * Supabase auth integration
  * Dockerized backend
  

---

## ⭐ Support

If this project resonates with you:

* 📩 Reach out for:

  * Internship/full-time SDE roles/ cloud infrastructure roles
  * System design discussions
  * Collaboration on network tooling
  * Linux/UNIX distro discussion

> NetShield is both a **real tool** and a **portfolio-grade systems project** – showing comfort with OS APIs, distributed thinking, and polished UI/UX.

