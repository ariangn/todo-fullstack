export type Tag = {
  id: string;
  name: string;
};

export async function fetchTags(): Promise<Tag[]> {
  const res = await fetch("/api/tags", { credentials: "include" });
  if (!res.ok) throw new Error("Failed to fetch tags");
  return (await res.json()) as Tag[];
}

export async function createTag(name: string): Promise<Tag> {
  const res = await fetch("/api/tags", {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ name }),
  });
  if (!res.ok) throw new Error("Failed to create tag");
  return (await res.json()) as Tag;
}

export async function deleteTag(id: string): Promise<void> {
  const res = await fetch(`/api/tags/${id}`, {
    method: "DELETE",
    credentials: "include",
  });
  if (!res.ok) throw new Error("Failed to delete tag");
}
