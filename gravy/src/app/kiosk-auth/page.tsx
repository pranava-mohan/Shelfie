"use client";
import { useRouter } from "next/navigation";
import { useState } from "react";

export default function KioskAuthPage() {
  const [token, setToken] = useState("");
  const router = useRouter();

  return (
    <>
      <input
        type="text"
        className="bg-yellow-600"
        onChange={(e) => setToken(e.target.value)}
      />
      <button
        className="bg-blue-700"
        onClick={() => {
          localStorage.setItem("kiosk_token", token);
          router.push("/check-in");
        }}
      >
        Authorize
      </button>
    </>
  );
}
