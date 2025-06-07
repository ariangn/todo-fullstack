import type { DragEndEvent, DragStartEvent } from "@dnd-kit/core";
import { DndContext, DragOverlay } from "@dnd-kit/core";
import { useEffect, useState, useCallback } from "react";
import Header from "../components/Header";
import SortFilterBar from "../components/SortFilterBar";
import Column from "../components/Column";
import CategoriesPanel from "../components/CategoriesPanel";
import type { Todo } from "../services/todoService";
import {
  fetchAllTodos,
  updateTodoStatus,
  deleteTodo as apiDeleteTodo,
} from "../services/todoService";
import type { Category } from "../types";
import {
  fetchCategories,
  createCategory,
  updateCategory,
  deleteCategory,
} from "../services/categoryService";
import type { Tag } from "../services/tagService";
import { fetchTags, createTag, deleteTag } from "../services/tagService";
import { logout, type User } from "../services/authService";
import CategoryModal from "../components/CategoryModal";
import TaskModal from "../components/TaskModal";
import TagModal from "../components/TagModal";
import TaskCard from "../components/TaskCard";

export default function DashboardPage({ user, onLogout }: { user: User; onLogout: () => void }) {
  const [todos, setTodos] = useState<Todo[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [tags, setTags] = useState<Tag[]>([]);
  const [sortBy, setSortBy] = useState<"dueDate" | "createdAt" | "updatedAt">("dueDate");
  const [filterCats, setFilterCats] = useState<string[]>([]);
  const [filterTags, setFilterTags] = useState<string[]>([]);
  const [modal, setModal] = useState<
    | null
    | { type: "addTask"; status: Todo["status"] }
    | { type: "editTask"; todo: Todo }
    | { type: "addCategory" }
    | { type: "editCategory"; categoryId: string }
    | { type: "addTag" }
    | { type: "editTag"; tagId: string }
  >(null);
  const [activeTodo, setActiveTodo] = useState<Todo | null>(null);

  const loadAllData = useCallback(async () => {
    let allTodos = await fetchAllTodos();
    if (!Array.isArray(allTodos)) allTodos = [];
    let filtered = allTodos;

    if (filterCats.length) {
      filtered = filtered.filter((t) => filterCats.includes(t.categoryId || ""));
    }
    if (filterTags.length) {
      filtered = filtered.filter((t) => t.tags.some((tag) => filterTags.includes(tag)));
    }

    filtered.sort((a, b) => {
      const aVal = a[sortBy] ? new Date(a[sortBy]!).getTime() : 0;
      const bVal = b[sortBy] ? new Date(b[sortBy]!).getTime() : 0;
      return aVal - bVal;
    });

    setTodos(filtered);
    setCategories(await fetchCategories());
    setTags(await fetchTags());
  }, [sortBy, filterCats, filterTags]);

  useEffect(() => {
    void loadAllData();
  }, [loadAllData]);

  const handleDragStart = (event: DragStartEvent) => {
    const dragged = todos.find((t) => t.id === event.active.id);
    setActiveTodo(dragged || null);
  };

  const handleDragEnd = async (event: DragEndEvent) => {
    const { active, over } = event;
    setActiveTodo(null);
    if (!over || active.id === over.id) return;

    const newStatus = over.id as Todo["status"];
    const draggedTodo = todos.find((t) => t.id === active.id);
    if (!draggedTodo || draggedTodo.status === newStatus) return;

    try {
      setTodos((prev) =>
        prev.map((t) => (t.id === active.id ? { ...t, status: newStatus } : t))
      );
      await updateTodoStatus(active.id as string, newStatus);
    } catch (err) {
      console.error("Failed to update status", err);
    }
  };

  const todosByStatus: Record<Todo["status"], Todo[]> = {
    TODO: todos.filter((t) => t.status === "TODO"),
    IN_PROGRESS: todos.filter((t) => t.status === "IN_PROGRESS"),
    COMPLETED: todos.filter((t) => t.status === "COMPLETED"),
  };

  const performLogout = async () => {
    await logout();
    onLogout();
  };

  const closeModal = () => setModal(null);

  return (
    <div className="h-screen flex flex-col">
      <Header
        userName={user.name || user.email}
        avatarUrl={user.avatarUrl}
        onLogout={performLogout}
      />

      <div className="flex-1 flex overflow-hidden">
        <div className="flex-1 flex flex-col p-6 overflow-hidden">
          <SortFilterBar onSortChange={setSortBy} onFilterChange={(cats, tgs) => { setFilterCats(cats); setFilterTags(tgs); }} />

          <DndContext onDragStart={handleDragStart} onDragEnd={handleDragEnd}>
            <div className="flex-1 flex space-x-4 overflow-x-auto">
              <Column
                title="To Do"
                status="TODO"
                todos={todosByStatus.TODO}
                onAddClick={() => setModal({ type: "addTask", status: "TODO" })}
              />
              <Column
                title="In Progress"
                status="IN_PROGRESS"
                todos={todosByStatus.IN_PROGRESS}
                onAddClick={() => setModal({ type: "addTask", status: "IN_PROGRESS" })}
              />
              <Column
                title="Completed"
                status="COMPLETED"
                todos={todosByStatus.COMPLETED}
                onAddClick={() => setModal({ type: "addTask", status: "COMPLETED" })}
              />
            </div>
            <DragOverlay>{activeTodo && <TaskCard todo={activeTodo} />}</DragOverlay>
          </DndContext>
        </div>

        <div className="w-1/3 border-l border-gray-200 p-4 overflow-y-auto">
          <CategoriesPanel
            categories={categories}
            onEdit={(id) => setModal({ type: "editCategory", categoryId: id })}
            onDelete={async (id) => { await deleteCategory(id); void loadAllData(); }}
            refreshCategories={() => void loadAllData()}
          />
          <div className="mt-6">
            <button onClick={() => setModal({ type: "addCategory" })} className="flex items-center text-primary hover:text-primary-dark">
              + Add Category
            </button>
          </div>
        </div>
      </div>

      {modal?.type === "addCategory" && (
        <CategoryModal
          mode="create"
          onSave={async (name, color) => {
            await createCategory(name, color);
            void loadAllData();
          }}
          onClose={closeModal}
        />
      )}

      {modal?.type === "editCategory" && modal.categoryId && (
        <CategoryModal
          mode="edit"
          categoryId={modal.categoryId}
          onSave={async (name, color) => {
            await updateCategory(modal.categoryId, name, color);
            void loadAllData();
            closeModal();
          }}
          onClose={closeModal}
        />
      )}

      {modal?.type === "addTask" && modal.status && (
        <TaskModal
          mode="create"
          status={modal.status}
          categories={categories}
          tags={(tags ?? []).map((t) => t.name)}
          onSave={async (data) => {
            await fetch("/api/todos", {
              method: "POST",
              credentials: "include",
              headers: { "Content-Type": "application/json" },
              body: JSON.stringify(data),
            });
            void loadAllData();
            closeModal();
          }}
          onClose={closeModal}
        />
      )}

      {modal?.type === "editTask" && modal.todo && (
        <TaskModal
          mode="edit"
          todo={modal.todo}
          categories={categories}
          tags={tags.map((t) => t.name)}
          onSave={async (data) => {
            await fetch(`/api/todos/${modal.todo.id}`, {
              method: "PUT",
              credentials: "include",
              headers: { "Content-Type": "application/json" },
              body: JSON.stringify(data),
            });
            void loadAllData();
            closeModal();
          }}
          onDelete={async () => {
            await apiDeleteTodo(modal.todo.id);
            void loadAllData();
            closeModal();
          }}
          onClose={closeModal}
        />
      )}

      {modal?.type === "addTag" && (
        <TagModal
          mode="create"
          onSave={async (name) => {
            await createTag(name);
            void loadAllData();
            closeModal();
          }}
          onClose={closeModal}
        />
      )}

      {modal?.type === "editTag" && modal.tagId && (
        <TagModal
          mode="edit"
          tagId={modal.tagId}
          onSave={async () => {
            void loadAllData();
            closeModal();
          }}
          onDelete={async () => {
            await deleteTag(modal.tagId);
            void loadAllData();
            closeModal();
          }}
          onClose={closeModal}
        />
      )}
    </div>
  );
}
