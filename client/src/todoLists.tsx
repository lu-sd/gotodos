import { todo } from "./App";

interface TodoListsProps {
  onToggle: (id: string, completed: boolean) => void; // Function to notify parent
  onDelete: (id: string) => void; // Function to notify parent
  todos: todo[];
}
export default function TodoLists({ todos, onToggle, onDelete }: TodoListsProps) {
  return (
    <div>
      <h1 className="text-start text-xl m-5 font-bold ">To Do List</h1>
      <ul className=" flex flex-col pl-3">
        {todos.map((todo) => {
          return (
            <li className="flex gap-2 " key={todo.id}>
              <input
                type="checkbox"
                checked={todo.completed}
                onChange={(e) => onToggle(todo.id, e.target.checked)}
              />
              <label htmlFor="item">{todo.title}</label>
              <button
                className="border-2 border-red-500 rounded border-spacing-3"
                onClick={() => {
                  if (
                    window.confirm(
                      'Are you sure you want to delete the task?',
                    )
                  ) onDelete(todo.id)
                }
                }
              >
                Delete
              </button>
            </li>
          );
        })}
      </ul>
    </div>
  )
}
