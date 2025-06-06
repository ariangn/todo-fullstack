import type { DragEndEvent } from "@dnd-kit/core";
import { DndContext } from "@dnd-kit/core";
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

type SortKey = "dueDate" | "createdAt" | "updatedAt";

type ModalType =
  | null
  | { type: "addTask"; status: Todo["status"] }
  | { type: "editTask"; todo: Todo }
  | { type: "addCategory" }
  | { type: "editCategory"; categoryId: string }
  | { type: "addTag" }
  | { type: "editTag"; tagId: string };

interface DashboardPageProps {
  user: User;
  onLogout: () => void;
}

export default function DashboardPage({
  user,
  onLogout,
}: DashboardPageProps) {
  const [todos, setTodos] = useState<Todo[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [tags, setTags] = useState<Tag[]>([]);
  const [sortBy, setSortBy] = useState<SortKey>("dueDate");
  const [filterCats, setFilterCats] = useState<string[]>([]);
  const [filterTags, setFilterTags] = useState<string[]>([]);
  const [modal, setModal] = useState<ModalType>(null);

  const loadAllData = useCallback(async () => {
    let allTodos = await fetchAllTodos();
    if (!Array.isArray(allTodos)) allTodos = [];
    let filtered = allTodos;

    if (filterCats.length) {
      filtered = filtered.filter((t) => filterCats.includes(t.categoryId || ""));
    }
    if (filterTags.length) {
      filtered = filtered.filter((t) =>
        t.tags.some((tag) => filterTags.includes(tag))
      );
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

  const handleDragEnd = async (event: DragEndEvent) => {
    const { active, over } = event;
    if (over && active.id !== over.id) {
      const newStatus = over.id as Todo["status"];
      await updateTodoStatus(active.id as string, newStatus);
      void loadAllData();
    }
  };

  const handleSortChange = (key: SortKey) => {
    setSortBy(key);
  };
  const handleFilterChange = (cats: string[], tgs: string[]) => {
    setFilterCats(cats);
    setFilterTags(tgs);
  };

  const handleAddCategory = () => {
    setModal({ type: "addCategory" });
  };
  const handleEditCategory = (categoryId: string) => {
    setModal({ type: "editCategory", categoryId });
  };
  const handleDeleteCategory = async (categoryId: string) => {
    await deleteCategory(categoryId);
    void loadAllData();
  };

  const handleAddTask = (status: Todo["status"]) => {
    setModal({ type: "addTask", status });
  };
  const handleDeleteTask = async (id: string) => {
    await apiDeleteTodo(id);
    void loadAllData();
  };

  const closeModal = () => {
    setModal(null);
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

  return (
    <div className="h-screen flex flex-col">
      <Header
        userName={user.name || user.email}
        avatarUrl={user.avatarUrl}
        onLogout={performLogout}
      />

      <div className="flex-1 flex overflow-hidden">
        <div className="flex-1 flex flex-col p-6 overflow-hidden">
          <SortFilterBar
            onSortChange={handleSortChange}
            onFilterChange={handleFilterChange}
          />

          <DndContext onDragEnd={handleDragEnd}>
            <div className="flex-1 flex space-x-4 overflow-x-auto">
              <Column
                title="To Do"
                status="TODO"
                todos={todosByStatus.TODO}
                onAddClick={() => handleAddTask("TODO")}
              />
              <Column
                title="In Progress"
                status="IN_PROGRESS"
                todos={todosByStatus.IN_PROGRESS}
                onAddClick={() => handleAddTask("IN_PROGRESS")}
              />
              <Column
                title="Completed"
                status="COMPLETED"
                todos={todosByStatus.COMPLETED}
                onAddClick={() => handleAddTask("COMPLETED")}
              />
            </div>
          </DndContext>
        </div>

        <div className="w-1/3 border-l border-gray-200 p-4 overflow-y-auto">
          <CategoriesPanel
            categories={categories}
            onEdit={handleEditCategory}
            onDelete={handleDeleteCategory}
            refreshCategories={() => void loadAllData()}
          />
          <div className="mt-6">
            <button
              onClick={handleAddCategory}
              className="flex items-center text-primary hover:text-primary-dark"
            >
              + Add Category
            </button>
          </div>
        </div>
      </div>

      {modal?.type === "addCategory" && (
        <CategoryModal
          mode="create"
          onSave={async (name: string, color: string) => {
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
          onSave={async (name: string, color: string) => {
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
            await handleDeleteTask(modal.todo.id);
            closeModal();
          }}
          onClose={closeModal}
        />
      )}

      {modal?.type === "addTag" && (
        <TagModal
          mode="create"
          onSave={async (name: string) => {
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
