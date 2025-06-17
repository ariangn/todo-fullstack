// src/components/TaskCard.tsx
import { XIcon, PencilIcon } from "lucide-react";
import type { Todo } from "../services/todoService";

interface TaskCardProps {
  todo: Todo;
  onEdit?: () => void;
  onDelete?: () => void;
}

export default function TaskCard({ todo, onEdit, onDelete }: TaskCardProps) {
  return (
    <div
      className="relative p-4 rounded shadow"
      style={{ backgroundColor: todo.category?.color || "#fff" }}
    >
      <div className="absolute top-2 right-2 flex space-x-1">
        {onEdit && (
          <button
            onPointerDown={e => e.stopPropagation()}
            onClick={(e) => {
              e.stopPropagation();
              onEdit();
            }}
            className="text-gray-800 hover:text-blue-500"
          >
            <PencilIcon className="h-4 w-4" />
          </button>
        )}
        {onDelete && (
          <button
            onPointerDown={e => e.stopPropagation()}
            onClick={(e) => {
              e.stopPropagation();
              onDelete();
            }}
            className="text-gray-800 hover:text-red-500"
          >
            <XIcon className="h-4 w-4" />
          </button>
        )}
      </div>

      <h3 className="font-semibold text-black">{todo.title}</h3>
      {todo.body && <p className="text-sm text-black mt-1">{todo.body}</p>}
      {todo.dueDate && (
        <p className="text-xs text-gray-900 mt-2">
          Due {new Date(todo.dueDate).toLocaleString()}
        </p>
      )}
      <div className="mt-2 flex flex-wrap space-x-1">
        {(todo.tags ?? []).map((tag) => (
          <span
            key={tag}
            className="text-xs bg-gray-200 text-gray-800 px-2 py-0.5 rounded"
          >
            {tag}
          </span>
        ))}
      </div>
    </div>
  );
}
