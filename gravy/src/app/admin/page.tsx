"use client";
import LockIcon from "@/components/lock.svg";
import { useForm, SubmitHandler } from "react-hook-form";
import { toast } from "react-hot-toast";
import axios from "axios";
import { API } from "@/config";

export default function AdminPage() {
  const { register, handleSubmit, formState } = useForm<{
    username: string;
    password: string;
  }>();
  const onSubmit: SubmitHandler<{
    username: string;
    password: string;
  }> = (data) => {
    return toast.promise(
      axios
        .post(API.BASE_URL + API.LIBRARIAN_LOGIN_URL, {
          username: data.username,
          password: data.password,
        })
        .then((res) => {
          console.log(res.data);
        }),
      {
        loading: "Signing In...",
        success: "Logged In",
        error: "Incorrect Credentials",
      },
      {
        position: "bottom-center",
      }
    );
  };

  return (
    <div className="w-full h-screen">
      <div className="flex items-center justify-center flex-col h-full">
        <LockIcon className="w-16 h-16 fill-amber-950" />
        <h1 className="text-3xl">Librarian Access</h1>

        <form
          className="lg:w-2/6 md:w-1/2 rounded-lg p-8 flex flex-col w-full"
          onSubmit={handleSubmit(onSubmit)}
        >
          <div className="relative mb-4">
            <label
              htmlFor="username"
              className="leading-7 text-sm text-gray-600"
            >
              Username
            </label>
            <input
              {...register("username", { required: true })}
              type="text"
              id="username"
              name="username"
              className="w-full bg-white rounded border border-gray-300 focus:border-indigo-500 focus:ring-2 focus:ring-indigo-200 text-base outline-none text-gray-700 py-1 px-3 leading-8 transition-colors duration-200 ease-in-out"
            />
          </div>
          <div className="relative mb-4">
            <label
              htmlFor="password"
              className="leading-7 text-sm text-gray-600"
            >
              Password
            </label>
            <input
              {...register("password", { required: true })}
              type="password"
              id="password"
              name="password"
              className="w-full bg-white rounded border border-gray-300 focus:border-indigo-500 focus:ring-2 focus:ring-indigo-200 text-base outline-none text-gray-700 py-1 px-3 leading-8 transition-colors duration-200 ease-in-out"
            />
          </div>
          <button
            type="submit"
            disabled={formState.isSubmitting}
            className="text-white bg-indigo-500 border-0 py-2 px-8 focus:outline-none hover:bg-indigo-600 cursor-pointer transition-all rounded text-lg disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Login
          </button>
        </form>
      </div>
    </div>
  );
}
