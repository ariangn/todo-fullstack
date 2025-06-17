// src/components/Column.tsx
import React from "react";
import { useDroppable, useDraggable } from "@dnd-kit/core";
import type { UniqueIdentifier } from "@dnd-kit/core";
import TaskCard from "./TaskCard";
import type { Todo } from "../services/todoService";
import { deleteTodo } from "../services/todoService";

interface ColumnProps {
  title: string;
  status: "TODO" | "IN_PROGRESS" | "COMPLETED";
  todos: Todo[];
  setModal: React.Dispatch<React.SetStateAction<any>>;
  loadAllData: () => Promise<void>;
}

export default function Column({
  title,
  status,
  todos,
  setModal,
  loadAllData,
}: ColumnProps) {
  const { isOver, setNodeRef: setDroppableRef } = useDroppable({ id: status });

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
          <DraggableTaskCard
            key={todo.id}
            todo={todo}
            setModal={setModal}
            loadAllData={loadAllData}
          />
        ))}
      </div>
    </div>
  );
}

function DraggableTaskCard({
  todo,
  setModal,
  loadAllData,
}: {
  todo: Todo;
  setModal: React.Dispatch<
    React.SetStateAction<
      | null
      | { type: "editTask"; todo: Todo }
      | { type: "addTask"; status: Todo["status"] }
    >
  >;
  loadAllData: () => Promise<void>;
}) {
  const {
    attributes,
    listeners,
    setNodeRef,
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
    cursor: isDragging ? "grabbing" : "grab",
  };

  return (
    <div
      ref={setNodeRef}
      style={style}
      {...listeners}
      {...attributes}
      className="mb-2"
    >
      <TaskCard
        todo={todo}
        onEdit={() => setModal({ type: "editTask", todo })}
        onDelete={async () => {
          try {
            await deleteTodo(todo.id);
            await loadAllData();
          } catch (err) {
            console.error("Delete failed", err);
          }
        }}
      />
    </div>
  );
}
