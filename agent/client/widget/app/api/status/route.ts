import { NextResponse } from "next/server";

const BACKEND_URL = process.env.BACKEND_URL ?? "http://localhost:8082";

export async function GET() {
  try {
    const res = await fetch(`${BACKEND_URL}/status`, {
      cache: "no-store",
    });

    if (!res.ok) {
      console.error("Backend /status failed:", res.status);
      return NextResponse.json(
        { error: "backend status failed" },
        { status: 500 }
      );
    }

    const data = await res.json();
    return NextResponse.json(data);
  } catch (err) {
    console.error("Error calling backend /status:", err);
    return NextResponse.json({ error: "backend unreachable" }, { status: 502 });
  }
}
