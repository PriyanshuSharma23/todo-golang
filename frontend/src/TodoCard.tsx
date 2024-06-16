import { useMutation, useQueryClient } from "react-query";
import { Todo } from "./App";
import { useState } from "react";
import { error } from "console";

const backendURL = import.meta.env.VITE_BACKEND_URL;

export default function TodoCard(props: Todo) {
  const [checked, setChecked] = useState(props.is_complete);
  const queryClient = useQueryClient();

  const deleteTodoMutation = useMutation({
    mutationKey: ["todos delete"],
    mutationFn: async (id: Number) => {
      let resp = await fetch(`${backendURL}/todos/${id}`, {
        method: "DELETE",
      });

      let data = await resp.json();

      if (resp.status != 200) {
        throw new Error(data.error);
      }

      return data.data as Todo;
    },
    onSuccess: async () => {
      await queryClient.refetchQueries({ queryKey: ["todos"] });
    },
  });

  const isCompletedMutation = useMutation({
    mutationKey: ["todos", props.id],
    mutationFn: async (checked: boolean) => {
      let resp = await fetch(`${backendURL}/todos`, {
        method: "PATCH",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          id: props.id,
          is_complete: checked,
        }),
      });

      let data = await resp.json();
      if (resp.status != 200) {
        console.error(data);
        throw new Error("error fetching");
      }

      return data.data;
    },
    onSuccess: async () => {
      await queryClient.refetchQueries({ queryKey: ["todos"] });
    },
  });

  return (
    <div className="group border space-y-2 p-4 relative">
      <p className="text-xl">{props.title}</p>
      <div className="flex justify-between">
        <p>
          <b>Due Date:</b> {props.due_on}
        </p>

        <input
          type="checkbox"
          checked={checked}
          className="accent-black"
          onChange={(e) => {
            setChecked(e.target.checked);
            isCompletedMutation.mutate(e.target.checked);
          }}
        />
      </div>
      <div className="text-sm text-gray-600">{props.body}</div>

      <button
        className="text-red-400 absolute top-0 right-1 opacity-0 pointer-events-none group-hover:pointer-events-auto group-hover:opacity-100"
        onClick={() => {
          deleteTodoMutation.mutate(props.id);
        }}
        disabled={deleteTodoMutation.isLoading}
      >
        {deleteTodoMutation.isLoading ? (
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill="none"
            className="stroke-gray-300"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
          >
            <path d="M12 2v4" />
            <path d="m16.2 7.8 2.9-2.9" />
            <path d="M18 12h4" />
            <path d="m16.2 16.2 2.9 2.9" />
            <path d="M12 18v4" />
            <path d="m4.9 19.1 2.9-2.9" />
            <path d="M2 12h4" />
            <path d="m4.9 4.9 2.9 2.9" />
          </svg>
        ) : (
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
          >
            <path d="M3 6h18" />
            <path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6" />
            <path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2" />
            <line x1="10" x2="10" y1="11" y2="17" />
            <line x1="14" x2="14" y1="11" y2="17" />
          </svg>
        )}
      </button>
    </div>
  );
}
