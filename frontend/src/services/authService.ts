export type User = { id: string; email: string; name?: string; avatarUrl?: string };

export async function login(email: string, password: string): Promise<User> {
  const res = await fetch("/api/users/login", {
    method: "POST",
    credentials: "include", // include HTTP-only cookie support
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ email, password }),
  });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(`Login failed: ${text}`);
  }
  // Assume backend sets a `Set-Cookie` header for JWT on success.
  return fetchMe();
}

export async function signup(
  email: string,
  password: string,
  name?: string,
  timezone?: string,
  avatarUrl?: string
): Promise<User> {
  const res = await fetch("/api/users/register", {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ email, password, name, timezone, avatarUrl }),
  });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(`Signup failed: ${text}`);
  }
  // Backend registers user, might return user object
  return fetchMe();
}

export async function fetchMe(): Promise<User> {
  const res = await fetch("/api/users/me", {
    method: "GET",
    credentials: "include",
  });
  if (!res.ok) {
    throw new Error("Not authenticated");
  }
  return (await res.json()) as User;
}

export async function logout(): Promise<void> {
  // Backend should clear cookie on this endpoint
  const res = await fetch("/api/users/logout", {
    method: "POST",
    credentials: "include",
  });
  if (!res.ok) {
    throw new Error("Logout failed");
  }
}

export async function getUserFromCookie(): Promise<User | null> {
  const res = await fetch("/api/auth/me", { credentials: "include" });
  if (!res.ok) return null;
  return (await res.json()) as User;
}
