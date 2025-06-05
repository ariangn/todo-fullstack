import React, { useEffect, useState } from "react";
import { BrowserRouter, Routes, Route, useNavigate, Navigate } from "react-router-dom";
import LoginPage from "./pages/LoginPage";
import SignupPage from "./pages/SignupPage";
import DashboardPage from "./pages/DashboardPage";
import {
  getUserFromCookie,
  login as loginService,
  signup as signupService,
  logout as logoutService,
} from "./services/authService";

// Our frontend “User” type must have required fields (id, email, name)
type User = {
  id: string;
  email: string;
  name: string;
  avatarUrl?: string;
};

function AppRouter() {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(true);
  const [user, setUser] = useState<User | null>(null);

  useEffect(() => {
    (async () => {
      // Skip auth check on public pages
      if (
        window.location.pathname === "/login" ||
        window.location.pathname === "/signup"
      ) {
        setLoading(false);
        return;
      }

      try {
        const fetched = await getUserFromCookie();
        if (fetched !== null) {
          setUser({
            id: fetched.id,
            email: fetched.email!,
            name: fetched.name || "",
            avatarUrl: fetched.avatarUrl,
          });
        } else {
          setUser(null);
        }
        setLoading(false);

        // Redirect logged-in users away from public pages
        if (
          window.location.pathname === "/login" ||
          window.location.pathname === "/signup"
        ) {
          navigate("/dashboard", { replace: true });
        }
      } catch {
        setUser(null);
        setLoading(false);
        navigate("/login", { replace: true });
      }
    })();
  }, [navigate]);

  if (loading) return <div>Loading…</div>;

  // LoginPage’s onLogin signature is (email: string, password: string) => void
  async function handleLogin(email: string, password: string) {
    const loggedIn = await loginService(email, password);
    setUser({
      id: loggedIn.id,
      email: loggedIn.email!,
      name: loggedIn.name || "",
      avatarUrl: loggedIn.avatarUrl,
    });
    navigate("/dashboard", { replace: true });
  }

  // SignupPage’s onSignup signature is (email: string, password: string, name: string, timezone: string, avatarUrl?: string) => void
  async function handleSignup(
    email: string,
    password: string,
    name: string,
    timezone: string,
    avatarUrl?: string
  ) {
    const signedUp = await signupService(
      email,
      password,
      name,
      timezone,
      avatarUrl
    );
    setUser({
      id: signedUp.id,
      email: signedUp.email!,
      name: signedUp.name || "",
      avatarUrl: signedUp.avatarUrl,
    });
    navigate("/dashboard", { replace: true });
  }

  async function handleLogout() {
    await logoutService();
    setUser(null);
    navigate("/login", { replace: true });
  }

  return (
    <Routes>
      <Route
        path="/login"
        element={<LoginPage onLogin={handleLogin} />}
      />
      <Route
        path="/signup"
        element={<SignupPage onSignup={handleSignup} />}
      />
      <Route
        path="/dashboard"
        element={
          user ? (
            <DashboardPage user={user} onLogout={handleLogout} />
          ) : (
            <Navigate to="/login" replace />
          )
        }
      />
      <Route
        path="*"
        element={<Navigate to={user ? "/dashboard" : "/login"} replace />}
      />
    </Routes>
  );
}

export default function App() {
  return (
    <BrowserRouter>
      <AppRouter />
    </BrowserRouter>
  );
}
