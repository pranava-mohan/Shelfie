"use client";

import BookIcon from "@/components/book.svg";
import ReturnIcon from "@/components/return.svg";
import Link from "next/link";

export default function Home() {
  return (
    <div className="w-full h-screen">
      <div className="flex items-center justify-center flex-col h-full">
        <h1 className="text-5xl">Welcome to the library</h1>

        <div className="flex mt-4">
          <Link
            href="/check-in"
            className="bg-blue-900 p-4 text-white rounded-2xl mr-4 flex flex-col items-center"
          >
            <BookIcon className="w-24 h-24" />
            <h2 className="text-2xl">Check-In Book</h2>
            <p>Click here to scan the books which you have</p>
          </Link>

          <Link
            href="/return"
            className="bg-red-800 p-4 text-white rounded-2xl flex flex-col items-center"
          >
            <ReturnIcon className="w-24 h-24" />
            <h2 className="text-2xl">Return Book</h2>
            <p>Click here to return the books you took</p>
          </Link>
        </div>
      </div>
    </div>
  );
}
