import { useEffect, useState } from "react";
import useSWR from "swr"
import AddTodo from "./addtodo";
import TodoLists from "./todoLists";

export type todo = {
  id: string;
  title: string;
  done: boolean;
};

export const ENDPOINT = "http://localhost:3000";

const fetcher = (url: string) =>
  fetch(`${ENDPOINT}/${url}`).then((r) => r.json());

function App() {
  const [todos, setTodos] = useState<todo[]>([])
  const { data, mutate } = useSWR<todo[]>("api/todos", fetcher)
  // const [todos, setTodos] = useState<todo[]>(() => {
  // Load todos from localStorage on initial render
  //   const savedTodos = localStorage.getItem("todos");
  //   return savedTodos ? JSON.parse(savedTodos) : [];
  // })
  //
  // useEffect(() => {
  //   localStorage.setItem("todos", JSON.stringify(todos))
  // }, [todos])

  const handleAddTodo = (newTodo: string) => {
    const newTodoItem = {
      id: crypto.randomUUID(),
      title: newTodo,
      done: false
    };
    setTodos([...todos, newTodoItem]);
  };

  async function creatTodo(newTodo: string) {
    const newTodoItem = await fetch(`${ENDPOINT}/api/todos`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        id: crypto.randomUUID(),
        title: newTodo,
        done: false
      }),
    }).then((r) => r.json());
    mutate(newTodoItem)
  }

  function toggleTodo(id: string, done: boolean) {
    setTodos((todos) => {
      return todos.map((todo) => {
        if (todo.id === id) {
          return { ...todo, done };
        }
        return todo;
      });
    });
  }
  async function markTodoDone(id: string) {
    const updated = await fetch(`${ENDPOINT}/api/todos/${id}/done`, {
      method: "PATCH",
    }).then((r) => r.json());
    mutate(updated)
  }

  function deleteTodo(id: string) {
    setTodos((todos) => {
      return todos.filter((todo) => todo.id !== id);
    });
  }
  async function removeTodo(id: string) {
    const updated = await fetch(`${ENDPOINT}/api/todos/${id}`, {
      method: "DELETE",
    }).then((r) => r.json());
    mutate(updated)
  }

  return (
    <div className="flex flex-col justify-center items-center space-y-8 mt-16">
      <AddTodo onAddTodo={creatTodo} />
      {data && <TodoLists todos={data} onToggle={markTodoDone} onDelete={removeTodo} />}
    </div>
  )
}

export default App;
