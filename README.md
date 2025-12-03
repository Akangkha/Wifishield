🛡️ NetShield — Real-Time Network Experience & Intelligent Wi-Fi Failover
Seamless Connectivity • Zero Interruptions • Smart Network Intelligence

A cross-platform desktop solution that ensures uninterrupted connectivity by dynamically switching Wi-Fi networks based on real-time performance metrics — built with Golang, gRPC streaming, PostgreSQL, Next.js and Electron.

⭐ Overview

NetShield is an intelligent network monitoring and auto-failover system designed to prevent users from losing connectivity during critical moments such as online examinations, remote work, video conferencing, and telemedicine.

Powered by a Go-based local agent, the system continuously evaluates Wi-Fi quality (Signal %, RSSI & latency) and automatically switches to the most stable network when thresholds degrade. Real-time metrics are streamed to a central analytics server and visualized in a Next.js + Electron desktop dashboard and an Admin Monitoring Console.

🚀 Key Highlights

Automatic Wi-Fi failover, reducing connectivity drops by 80% during peak usage.

Real-time telemetry pipeline using gRPC streaming with monitoring resolution < 100 ms.

Live desktop widget with signal gauge & experience score using Electron.

Centralized insight server with PostgreSQL storing time-series metrics.

Admin Console with filters (Exam / Remote-Work / Telemedicine) to monitor all devices.

Designed as a distributed system, production-ready with Docker.
