import React, { Dispatch, SetStateAction, useState } from "react";
import { useMutation, useQueries, useQuery, useQueryClient } from "react-query";
import { Todo } from "./App";

const backendURL = import.meta.env.VITE_BACKEND_URL;

const today = () => {
  let date = new Date();
  let tDate = new Date(0);

  tDate.setDate(date.getDate());
  tDate.setMonth(date.getMonth());
  tDate.setFullYear(date.getFullYear());
  return tDate;
};

export default function TodoForm(props: {
  open: boolean;
  setOpen: Dispatch<SetStateAction<boolean>>;
}) {
  const queryClient = useQueryClient();

  const createTodoMutation = useMutation({
    mutationKey: ["todoscreate"],
    mutationFn: async (todo: Partial<Todo>) => {
      let resp = await fetch(`${backendURL}/todos`, {
        method: "POST",
        body: JSON.stringify({
          title: todo.title,
          body: todo.body,
          due_on: todo.due_on,
        }),
      });

      let data = await resp.json();

      if (resp.status != 201) {
        throw new FormError(data.error);
      }

      return data.data;
    },

    onSuccess: async () => {
      await queryClient.refetchQueries({ queryKey: ["todos"] });
    },
  });

  if (!props.open) {
    return null;
  }

  return (
    <>
      <div className="fixed inset-0 bg-black/50"></div>
      <form
        action=""
        className="p-4 rounded-md fixed top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 bg-white w-full max-w-md space-y-4"
        onSubmit={(ev) => {
          ev.preventDefault();
          let formData = new FormData(ev.currentTarget);

          let dueOn = formData.get("due_on");

          let todo: Partial<Todo> = {
            title: formData.get("title")?.toString() ?? "",
            due_on: dueOn
              ? new Date(dueOn.toString()).toISOString()
              : undefined,
            body: formData.get("body")?.toString() ?? "",
          };

          createTodoMutation.mutate(todo, {
            onSuccess: async () => {
              props.setOpen(false);
            },
          });
        }}
      >
        <button
          className="absolute top-2 right-2"
          type="button"
          onClick={() => {
            props.setOpen(false);
          }}
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="3"
            strokeLinecap="round"
            strokeLinejoin="round"
          >
            <path d="M18 6 6 18" />
            <path d="m6 6 12 12" />
          </svg>
        </button>
        <h1 className="text-2xl text-center ">New Todo</h1>

        <div className="grid gap-2">
          <label htmlFor="title" className=" text-base">
            Title
          </label>
          <input
            type="text"
            id="title"
            name="title"
            placeholder="Get bananas"
            className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 "
          />
        </div>

        <div className="grid gap-2">
          <label htmlFor="due-date" className=" text-base">
            Due on
          </label>
          <input
            type="date"
            id="due-date"
            name="due_on"
            className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 "
          />
        </div>

        <div className="grid gap-2">
          <label htmlFor="todo-body" className=" text-base">
            Body
          </label>
          <textarea
            id="todo-body"
            name="body"
            rows={4}
            className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 "
          />
        </div>

        <button disabled={createTodoMutation.isLoading}>
          {!createTodoMutation.isLoading ? "Create Todo" : "Creating"}
        </button>
      </form>
    </>
  );
}

class FormError extends Error {
  private errors: Record<string, string>;

  constructor(errors: Record<string, string>) {
    super();
    this.errors = errors;
  }

  get getErrors() {
    return this.errors;
  }
}
