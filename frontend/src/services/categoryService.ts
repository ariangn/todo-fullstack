export type Category = {
  id: string;
  name: string;
  color: string;
};

export async function fetchCategories(): Promise<Category[]> {
  const res = await fetch("/api/categories", { credentials: "include" });
  if (!res.ok) throw new Error("Failed to fetch categories");
  return (await res.json()) as Category[];
}

export async function createCategory(name: string, color: string, description?: string): Promise<Category> {
  const res = await fetch("/api/categories", {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ name, color, description }),
  });
  if (!res.ok) throw new Error("Failed to create category");
  return (await res.json()) as Category;
}

export async function updateCategory(id: string, name: string, color: string, description?: string): Promise<Category> {
  const res = await fetch(`/api/categories/${id}`, {
    method: "PUT",
    credentials: "include",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ name, color, description }),
  });
  if (!res.ok) throw new Error("Failed to update category");
  return (await res.json()) as Category;
}

export async function deleteCategory(id: string): Promise<void> {
  const res = await fetch(`/api/categories/${id}`, {
    method: "DELETE",
    credentials: "include",
  });
  if (!res.ok) throw new Error("Failed to delete category");
}
