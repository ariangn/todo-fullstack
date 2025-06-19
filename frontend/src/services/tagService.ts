export type Tag = {
  id: string;
  name: string;
};

const API = import.meta.env.VITE_API_URL as string;

export async function fetchTags(): Promise<Tag[]> {
  const res = await fetch(`${API}/tags`, { credentials: "include" });
  if (!res.ok) throw new Error("Failed to fetch tags");
  return (await res.json()) as Tag[];
}

export async function createTag(name: string): Promise<Tag> {
  const res = await fetch(`${API}/tags`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    credentials: "include",
    body: JSON.stringify({ name }),
  });

  if (!res.ok) {
    const msg = await res.text();
    console.error("createTag failed:", msg);
    throw new Error("Failed to create tag");
  }

  return await res.json();
}

export async function deleteTag(id: string): Promise<void> {
  const res = await fetch(`${API}/tags/${id}`, {
    method: "DELETE",
    credentials: "include",
  });
  if (!res.ok) throw new Error("Failed to delete tag");
}
