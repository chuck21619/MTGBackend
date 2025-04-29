const loginView = document.getElementById("loginView");
const registerView = document.getElementById("registerView");
const dashboardView = document.getElementById("dashboardView");
const messageLabel = document.getElementById("messageLabel");
const fetchMessageButton = document.getElementById("fetchMessageButton");
const logoutButton = document.getElementById("logoutButton");

document.getElementById("showRegister").addEventListener("click", () => {
    loginView.style.display = "none";
    registerView.style.display = "block";
});

document.getElementById("showLogin").addEventListener("click", () => {
    registerView.style.display = "none";
    loginView.style.display = "block";
});

document.getElementById("loginForm").addEventListener("submit", async (e) => {
    e.preventDefault();
    const username = document.getElementById("loginUsername").value;
    const password = document.getElementById("loginPassword").value;

    const res = await fetch("/api/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password })
    });

    const data = await res.json();
    console.log(data);

    if (data.access_token) {
        localStorage.setItem("access_token", data.access_token);
        loginView.style.display = "none";
        dashboardView.style.display = "block";
    } else {
        alert(data.message || "Login response received.");
    }
});

document.getElementById("registerForm").addEventListener("submit", async (e) => {
    e.preventDefault();
    const username = document.getElementById("registerUsername").value;
    const email = document.getElementById("registerEmail").value;
    const password = document.getElementById("registerPassword").value;

    const res = await fetch("/api/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, email, password })
    });

    const data = await res.json();
    alert(data.message || "Registration response received.");
});

// Button to fetch a message from a protected route
fetchMessageButton.addEventListener("click", async () => {
    const token = localStorage.getItem("access_token");
    if (!token) {
        alert("Please log in first.");
        return;
    }

    new_email = "nothing";

    const res = await fetch("/api/update-email", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        },
        body: JSON.stringify({ new_email: newEmail })
    });

    const data = await res.json();
    alert(data.message || "Registration response received.");
});

// Button to log the user out
logoutButton.addEventListener("click", () => {
    localStorage.removeItem("access_token");
    dashboardView.style.display = "none";
    loginView.style.display = "block";
    alert("You have been logged out.");
});
