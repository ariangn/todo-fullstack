import { useEffect, useState } from "react";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import LoginPage from "./pages/LoginPage";
import SignupPage from "./pages/SignupPage";
import DashboardPage from "./pages/DashboardPage";
import {
  getUserFromCookie,
  login as loginService,
  signup as signupService,
  logout as logoutService,
} from "./services/authService";

// user type
type User = {
  id: string;
  email: string;
  name: string;
  avatarUrl?: string;
};

export default function App() {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // check if user is already logged in via cookie
    (async () => {
      try {
        const fetched = await getUserFromCookie();
        if (fetched) {
          setUser({
            id: fetched.id,
            email: fetched.email!,
            name: fetched.name || "",
            avatarUrl: fetched.avatarUrl,
          });
        } else {
          setUser(null);
        }
      } catch {
        setUser(null);
      } finally {
        setLoading(false);
      }
    })();
  }, []);

  // auth handlers
  async function handleLogin(email: string, password: string) {
    const loggedIn = await loginService(email, password);
    setUser({
      id: loggedIn.id,
      email: loggedIn.email!,
      name: loggedIn.name || "",
      avatarUrl: loggedIn.avatarUrl,
    });
  }

  async function handleSignup(
    email: string,
    password: string,
    name: string,
    timezone: string,
    avatarUrl?: string
  ) {
    const signedUp = await signupService(email, password, name, timezone, avatarUrl);
    setUser({
      id: signedUp.id,
      email: signedUp.email!,
      name: signedUp.name || "",
      avatarUrl: signedUp.avatarUrl,
    });
  }

  async function handleLogout() {
    await logoutService();
    setUser(null);
  }

  if (loading) return <div>Loading…</div>;

  return (
    <BrowserRouter>
      <Routes>
        <Route
          path="/login"
          element={
            user ? <Navigate to="/dashboard" replace /> : <LoginPage onLogin={handleLogin} />
          }
        />
        <Route
          path="/signup"
          element={
            user ? <Navigate to="/dashboard" replace /> : <SignupPage onSignup={handleSignup} />
          }
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
    </BrowserRouter>
  );
}
