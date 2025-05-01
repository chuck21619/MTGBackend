import { useState } from "react";
import "./App.css";

export default function App() {
  const [view, setView] = useState("login"); // "login", "register", "dashboard"
  const [accessToken, setAccessToken] = useState(localStorage.getItem("access_token") || "");
  const [message, setMessage] = useState("");

  const handleLogin = async (e) => {
    e.preventDefault();
    const username = e.target.username.value;
    const password = e.target.password.value;

    const res = await fetch("/api/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, password }),
    });

    const data = await res.json();
    if (data.access_token) {
      localStorage.setItem("access_token", data.access_token);
      setAccessToken(data.access_token);
      setView("dashboard");
    } else {
      alert(data.message || "Login failed");
    }
  };

  const handleRegister = async (e) => {
    e.preventDefault();
    const username = e.target.username.value;
    const email = e.target.email.value;
    const password = e.target.password.value;

    const res = await fetch("/api/register", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, email, password }),
    });

    const data = await res.json();
    alert(data.message || "Registered");
  };

  const handleLogout = () => {
    localStorage.removeItem("access_token");
    setAccessToken("");
    setView("login");
    alert("You have been logged out.");
  };

  const handleUpdateEmail = async () => {
    const newEmail = document.getElementById("newEmailInput").value.trim();
    if (!accessToken || !newEmail) {
      alert("Please log in and enter a valid email.");
      return;
    }

    const res = await authFetch("/api/update-email", {
      method: "POST",
      body: JSON.stringify({ new_email: newEmail }),
    });

    const data = await res.json();
    alert(data.message || "No response message.");
  };

  async function authFetch(url, options = {}) {
    let token = localStorage.getItem("access_token");

    const res = await fetch(url, {
      ...options,
      headers: {
        ...(options.headers || {}),
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
      },
    });

    if (res.status === 401) {
      const refreshRes = await fetch("/api/refresh-token", {
        method: "POST",
        credentials: "include",
      });

      const refreshData = await refreshRes.json();
      if (refreshRes.ok && refreshData.access_token) {
        localStorage.setItem("access_token", refreshData.access_token);
        setAccessToken(refreshData.access_token);
        return authFetch(url, options);
      } else {
        alert("Session expired. Please log in again.");
        setView("login");
        return;
      }
    }

    return res;
  }

  return (
    <div className="App">
      {view === "login" && (
        <div>
          <h2>Login</h2>
          <form onSubmit={handleLogin}>
            <input name="username" type="text" placeholder="Username" required />
            <input name="password" type="password" placeholder="Password" required />
            <button type="submit">Login</button>
          </form>
          <a onClick={() => setView("register")} style={{ cursor: "pointer" }}>
            Don't have an account? Register
          </a>
        </div>
      )}

      {view === "register" && (
        <div>
          <h2>Register</h2>
          <form onSubmit={handleRegister}>
            <input name="username" type="text" placeholder="Username" required />
            <input name="email" type="email" placeholder="Email" required />
            <input name="password" type="password" placeholder="Password" required />
            <button type="submit">Register</button>
          </form>
          <a onClick={() => setView("login")} style={{ cursor: "pointer" }}>
            Already have an account? Login
          </a>
        </div>
      )}

      {view === "dashboard" && (
        <div>
          <h2>Welcome to your dashboard!</h2>
          <p>{message || "Your secret message will appear here."}</p>
          <button onClick={() => console.log("Fetch secret message here")}>
            Fetch Secret Message
          </button>
          <div>
            <h3>Update Email</h3>
            <input id="newEmailInput" type="email" placeholder="Enter new email" />
            <button onClick={handleUpdateEmail}>Update Email</button>
          </div>
          <button onClick={handleLogout}>Logout</button>
        </div>
      )}
    </div>
  );
}
