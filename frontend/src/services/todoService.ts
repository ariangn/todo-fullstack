export type Todo = {
  id: string;
  title: string;
  body?: string;
  status: "TODO" | "IN_PROGRESS" | "COMPLETED";
  dueDate?: string;
  tagIds: string[];
  tags?: string[];
  categoryId?: string;
  category?: { id: string; color: string; name: string };
  createdAt: string;
  updatedAt: string;
};

export async function fetchAllTodos(): Promise<Todo[]> {
  const res = await fetch("/api/todos", { credentials: "include" });
  if (!res.ok) throw new Error("Failed to fetch todos");
  return (await res.json()) as Todo[];
}

export async function createTodo(data: {
  title: string;
  body?: string;
  dueDate?: string;
  status: Todo["status"];
  categoryId?: string;
  tagIds: string[];
}): Promise<Todo> {
  const res = await fetch("/api/todos", {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });
  if (!res.ok) throw new Error("Failed to create todo");
  return (await res.json()) as Todo;
}

export async function updateTodo(
  id: string,
  data: Partial<{
    title: string;
    body?: string;
    dueDate?: string;
    status: Todo["status"];
    categoryId?: string;
    tagIds: string[];
  }>
): Promise<Todo> {
  const res = await fetch(`/api/todos/${id}`, {
    method: "PUT",
    credentials: "include",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });
  if (!res.ok) throw new Error("Failed to update todo");
  return (await res.json()) as Todo;
}

export async function updateTodoStatus(todoId: string, newStatus: "TODO" | "IN_PROGRESS" | "COMPLETED") {
  const res = await fetch(`/api/todos/${todoId}/status`, {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
    },
    credentials: "include", // ensures session cookies are sent
    body: JSON.stringify({ status: newStatus }),
  });

  if (!res.ok) {
    const errorText = await res.text();
    throw new Error(`Failed to update status: ${errorText}`);
  }

  return await res.json(); // optional: return the updated todo
}

export async function deleteTodo(id: string): Promise<void> {
  const res = await fetch(`/api/todos/${id}`, {
    method: "DELETE",
    credentials: "include",
  });
  if (!res.ok) throw new Error("Failed to delete todo");
}

export async function duplicateTodo(id: string): Promise<Todo> {
  const res = await fetch(`/api/todos/${id}/duplicate`, {
    method: "POST",
    credentials: "include",
  });
  if (!res.ok) throw new Error("Failed to duplicate todo");
  return (await res.json()) as Todo;
}
