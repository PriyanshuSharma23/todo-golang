import { useQuery } from "react-query";
import TodoCard from "./TodoCard";
import TodoForm from "./TodoForm";
import { useState } from "react";

const backendURL = import.meta.env.VITE_BACKEND_URL;

export type Todo = {
  id: Number;
  title: string;
  is_complete: boolean;
  due_on: string;
  created_at: string;
  body: string;
};

export default function App() {
  const [open, setOpen] = useState(false);
  const todosQuery = useQuery({
    queryKey: ["todos"],
    queryFn: async () => {
      const todosResp = await fetch(`${backendURL}/todos`, {
        method: "GET",
      });

      let data = await todosResp.json();
      if (todosResp.status != 200) {
        throw new Error(data.error);
      }

      return data.data as Todo[];
    },
  });

  if (todosQuery.isLoading) {
    return "loading...";
  }

  return (
    <main>
      <div className="container grid gap-5 mx-auto p-10">
        {todosQuery.data && todosQuery.data.length === 0 && "No todos"}
        {todosQuery.data?.map((todo) => (
          <TodoCard {...todo} key={todo.id.toString()} />
        ))}
      </div>

      <button
        className="bg-black p-2 rounded-full fixed bottom-6 right-6"
        onClick={() => {
          setOpen(true);
        }}
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="32"
          height="32"
          viewBox="0 0 24 24"
          fill="none"
          stroke="white"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
        >
          <path d="M5 12h14" />
          <path d="M12 5v14" />
        </svg>
      </button>

      <TodoForm open={open} setOpen={setOpen} />
    </main>
  );
}
