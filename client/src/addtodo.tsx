import { useState } from "react"
interface AddTodoProps {
  onAddTodo: (newTodoTerm: string) => void; // Function to notify parent
}
export default function AddTodo({ onAddTodo }: AddTodoProps) {
  const [newItem, setNewItem] = useState("");

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (newItem.trim()) {
      onAddTodo(newItem)
      setNewItem("");
    }
  }
  return (
    <div>
      <h1 className="text-blue-700 font-bold mb-9">what do you want to do?</h1>
      <form onSubmit={handleSubmit} className="flex flex-col gap-3">
        {/* <label htmlFor="todo">New Item</label> */}
        <input
          type="text"
          id="todo"
          value={newItem}
          onChange={(e) => setNewItem(e.target.value)}
          className="border-2 border-sky-500  rounded-md bg-sky-200"
        />
        <button className="border-2 border-sky-500  rounded-md">Add</button>
      </form>
    </div>

  )
}
