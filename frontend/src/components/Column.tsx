import React from "react";
import {
  useDroppable,
  useDraggable,
} from "@dnd-kit/core";
import type { UniqueIdentifier } from "@dnd-kit/core";
import TaskCard from "./TaskCard";
import type { Todo } from "../services/todoService";

interface ColumnProps {
  title: string;
  status: "TODO" | "IN_PROGRESS" | "COMPLETED";
  todos: Todo[];
  onAddClick: () => void;
}

export default function Column({
  title,
  status,
  todos,
  onAddClick,
}: ColumnProps) {
  const { isOver, setNodeRef: setDroppableRef } = useDroppable({
    id: status,
  });

  return (
    <div
      ref={setDroppableRef}
      className={`flex-1 flex flex-col bg-gray-50 rounded p-2 ${
        isOver ? "bg-gray-100" : ""
      }`}
    >
      <h2 className="text-lg font-medium mb-2">{title}</h2>
      <div className="flex-1 space-y-2 overflow-y-auto">
        {todos.map((todo) => (
          <DraggableTaskCard key={todo.id} todo={todo} />
        ))}
      </div>
      <button
        onClick={onAddClick}
        className="mt-2 flex items-center justify-center text-primary hover:text-primary-dark"
      >
        + Add Task
      </button>
    </div>
  );
}

function DraggableTaskCard({ todo }: { todo: Todo }) {
  const {
    attributes,
    listeners,
    setNodeRef: setDraggableRef,
    transform,
    isDragging,
  } = useDraggable({
    id: todo.id as UniqueIdentifier,
    data: { status: todo.status },
  });

  const style: React.CSSProperties = {
    transform: transform
      ? `translate3d(${transform.x}px, ${transform.y}px, 0)`
      : undefined,
    opacity: isDragging ? 0.5 : 1,
  };

  return (
    <div ref={setDraggableRef} style={style} {...listeners} {...attributes}>
      <TaskCard todo={todo} innerRef={setDraggableRef} draggableProps={attributes} dragHandleProps={listeners} />
    </div>
  );
}
