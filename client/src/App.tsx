import { useEffect, useState } from "react";
import AddTodo from "./addtodo";
import TodoLists from "./todoLists";
type todo = {
  id: string;
  title: string;
  completed: boolean;
};
function App() {
  // const [todos, setTodos] = useState<todo[]>([])
  const [todos, setTodos] = useState<todo[]>(() => {
    // Load todos from localStorage on initial render
    const savedTodos = localStorage.getItem("todos");
    return savedTodos ? JSON.parse(savedTodos) : [];
  })

  useEffect(() => {
    localStorage.setItem("todos", JSON.stringify(todos))
  }, [todos])

  const handleAddTodo = (newTodo: string) => {
    const newTodoItem = {
      id: crypto.randomUUID(),
      title: newTodo,
      completed: false
    };
    setTodos([...todos, newTodoItem]);
  };

  function toggleTodo(id: string, completed: boolean) {
    setTodos((todos) => {
      return todos.map((todo) => {
        if (todo.id === id) {
          return { ...todo, completed };
        }
        return todo;
      });
    });
  }

  function deleteTodo(id: string) {
    setTodos((todos) => {
      return todos.filter((todo) => todo.id !== id);
    });
  }

  return (
    <div className="flex flex-col justify-center items-center space-y-8 mt-16">
      <AddTodo onAddTodo={handleAddTodo} />
      <TodoLists todos={todos} onToggle={toggleTodo} onDelete={deleteTodo} />
    </div>
  )
}

export default App;
