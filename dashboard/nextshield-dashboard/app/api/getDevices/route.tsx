import { NextResponse } from "next/server";

export async function GET(request: Request) {
  const { searchParams } = new URL(request.url);
  const ssid = searchParams.get("ssid");

  if (!ssid) {
    return NextResponse.json({ error: "SSID is required" }, { status: 400 });
  }

  const res = await fetch(
    `http://localhost:8082/api/admin/links/devices?ssid=${encodeURIComponent(
      ssid
    )}`,
    {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    }
  );

  if (!res.ok) {
    return NextResponse.json(
      { error: "Backend fetch failed" },
      { status: res.status }
    );
  }

  const data = await res.json();
  return NextResponse.json(data);
}
