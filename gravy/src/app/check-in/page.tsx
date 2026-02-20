"use client";
import QRCode from "react-qr-code";
import generateRoomID from "@/utils/generate_room_id";
import { useEffect, useState, useCallback } from "react";
import useWebSocket from "react-use-websocket";
import AddIcon from "@/components/add.svg";
import CloseIcon from "@/components/close.svg";
import { useKioskClient } from "@/hooks/useKioskClient";
import { useRouter } from "next/navigation";
import { Modal } from "react-responsive-modal";
import { Scanner } from "@yudiel/react-qr-scanner";
import toast from "react-hot-toast";

export default function CheckInPage() {
  const [roomId, setRoomId] = useState("");
  console.log(roomId.split("?")[0])
  const [user, setUser] = useState<{ id: string; name: string } | null>(null);
  const [loading, setLoading] = useState(false);
  const kioskClient = useKioskClient();
  const [modelOpen, setModalOpen] = useState(false);
  const router = useRouter();
  const [books, setBooks] = useState<
    Array<{
      id: string;
      title: string;
      author: string;
      publisher: string;
      isbn: string;
    }>
  >([]);

  useWebSocket(
    "ws://localhost:8000/ws/" + roomId,
    {
      onMessage: async (msg) => {
        setLoading(true);
        setUser(
          (await kioskClient.post("/user", { id: JSON.parse(msg.data).id }))
            .data
        );
        setLoading(false);
      },
    },
    !!roomId
  );

  useEffect(() => {
    if (!roomId) {
      const id = generateRoomID(256);
      setRoomId(id + "?token=" + localStorage.getItem("kiosk_token"));
    }
  }, []);

  const reset = () => {
    setRoomId("");
    setUser(null);
    setLoading(false);
    setBooks([]);
    router.replace("/");
  };

  const [checkingIn, setCheckingIn] = useState(false);
  const checkInBooks = useCallback(async () => {
    return toast.promise(
      (async () => {
        await kioskClient.post("/book/check-in", {
          user_id: user?.id,
          book_ids: (books || []).map((b) => b.id),
        });
        reset();
      })(),
      {
        loading: "Checking in books...",
        success: "Books checked in successfully!",
        error: "Failed to check in books",
      }
    );
  }, [books, user, kioskClient]);

  return (
    <>
      <Modal open={modelOpen} onClose={() => setModalOpen(false)} center>
        <div className="m-2">
          <h2 className="text-xl my-4">Scan the books one by one</h2>
          <Scanner
            onScan={async (result) => {
              if (
                result[0].rawValue &&
                !(books || []).find((b) => b.id === result[0].rawValue)
              ) {
                toast.promise(
                  (async () => {
                    let res = await kioskClient.post("/book/get", {
                      book_id: result[0].rawValue,
                    });
                    setBooks((prev) => [res.data, ...prev]);
                  })(),
                  {
                    loading: "Adding book...",
                    success: "Book added!",
                    error: "Failed to add book",
                  }
                );
              }

              setModalOpen(false);
            }}
            onError={(error) => console.log(error)}
            constraints={{
              aspectRatio: 1,
              width: { ideal: 300 },
              height: { ideal: 300 },
            }}
          />
        </div>
      </Modal>
      <div className="w-full h-screen">
        <div
          className={`flex items-center flex-col h-full ${!user && "justify-center"
            }`}
        >
          <div className="max-w-4xl flex items-center flex-col">
            <h1 className={`text-4xl ${user && "opacity-0"}`}>
              Scan the QR code using the library app
            </h1>
            {roomId && !loading && !user && (
              <>
                <QRCode className="mt-4" value={roomId.split("?")[0]} />
              </>
            )}
            {loading && !user && <p className="mt-4">Checking in...</p>}
            {user && (
              <div className="w-full">
                <div className="mt-4 flex flex-col items-center">
                  <h2 className="text-2xl">Checked in as:</h2>
                  <p className="text-xl font-semibold">{user.name}</p>
                </div>
                <div className="w-full flex">
                  <button
                    onClick={() => reset()}
                    className="bg-red-700 p-2 text-white rounded-lg flex flex-col items-center mr-auto cursor-pointer"
                  >
                    Cancel
                  </button>

                  <button
                    disabled={(books || []).length === 0 || checkingIn}
                    onClick={() => {
                      setCheckingIn(true);
                      checkInBooks().finally(() => setCheckingIn(false));
                    }}
                    className="bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed p-2 text-white rounded-lg mr-4 flex flex-col items-center ml-auto cursor-pointer"
                  >
                    Check-In
                  </button>
                </div>
                <div className="w-full mt-4">
                  <button
                    onClick={() => setModalOpen(true)}
                    className="bg-blue-600 flex rounded-lg text-white p-2 mx-auto cursor-pointer"
                  >
                    <AddIcon className="w-6 h-6 fill-white" />
                    Add Books
                  </button>

                  {(books || []).map((book) => (
                    <div
                      key={book.id}
                      className="flex w-full py-2 mt-4 border-b border-gray-300"
                    >
                      <div>
                        <h3 className="text-lg font-semibold">{book.title}</h3>
                        <p>Author: {book.author}</p>
                        <p>Publisher: {book.publisher}</p>
                        <p>ISBN: {book.isbn}</p>
                      </div>
                      <CloseIcon
                        className="w-6 h-6 fill-black cursor-pointer ml-auto"
                        onClick={() => {
                          setBooks((prev) =>
                            prev.filter((b) => b.id !== book.id)
                          );
                        }}
                      />
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </>
  );
}
