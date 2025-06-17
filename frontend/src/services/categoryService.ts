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

export async function updateCategory(
  id: string,
  name: string,
  color: string,
  description?: string
): Promise<Category> {
  const res = await fetch(`/api/categories/${id}`, {
    method: "PUT",
    credentials: "include",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ name, color, description }),
  });

  // Read raw text (even if empty)
  const text = await res.text();

  // If we got an error status, throw with the server’s message (or empty)
  if (!res.ok) {
    throw new Error(`Update failed (${res.status}): ${text || res.statusText}`);
  }

  // If body is empty, that’s unexpected—but we can return a minimal stub or throw
  if (!text) {
    // Option A: throw so you know something’s off
    throw new Error("Update succeeded but server returned no data");
    // —or—
    // Option B: return a placeholder Category object if you can reconstruct it locally
  }

  // Otherwise parse and return the updated category
  return JSON.parse(text) as Category;
}

export async function deleteCategory(id: string): Promise<void> {
  const res = await fetch(`/api/categories/${id}`, {
    method: "DELETE",
    credentials: "include",
  });
  if (!res.ok) throw new Error("Failed to delete category");
}
