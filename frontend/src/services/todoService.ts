export type Todo = {
  id: string;
  title: string;
  body?: string;
  status: "TODO" | "IN_PROGRESS" | "COMPLETED";
  dueDate?: string;
  tags: string[];
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
  tags: string[];
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
    tags: string[];
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

export async function updateTodoStatus(
  id: string,
  status: Todo["status"]
): Promise<Todo> {
  const res = await fetch(`/api/todos/${id}/status`, {
    method: "PUT",
    credentials: "include",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ status }),
  });
  if (!res.ok) throw new Error("Failed to update status");
  return (await res.json()) as Todo;
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
