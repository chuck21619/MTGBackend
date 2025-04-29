const loginView = document.getElementById("loginView");
const registerView = document.getElementById("registerView");
const dashboardView = document.getElementById("dashboardView");
const messageLabel = document.getElementById("messageLabel");
const fetchMessageButton = document.getElementById("fetchMessageButton");
const logoutButton = document.getElementById("logoutButton");
const newEmailInput = document.getElementById("newEmailInput");
const updateEmailButton = document.getElementById("updateEmailButton");

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
    console.log("nothing here yet");
});

updateEmailButton.addEventListener("click", async () => {
    const token = localStorage.getItem("access_token");
    const newEmail = newEmailInput.value.trim();
    if (!token || !newEmail) {
        alert("Please log in and enter a valid email.");
        return;
    }

    const res = await authFetch("/api/update-email", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        },
        body: JSON.stringify({ new_email: newEmail })
    });

    const data = await res.json();
    alert(data.message || "No response message.");
});


// Button to log the user out
logoutButton.addEventListener("click", () => {
    localStorage.removeItem("access_token");
    dashboardView.style.display = "none";
    loginView.style.display = "block";
    alert("You have been logged out.");
});

async function authFetch(url, options = {}) {
    let token = localStorage.getItem("access_token");

    const res = await fetch(url, {
        ...options,
        headers: {
            ...options.headers,
            "Authorization": `Bearer ${token}`,
            "Content-Type": "application/json",
        },
    });

    if (res.status === 401) {
        // Try refresh
        const refreshRes = await fetch("/api/refresh-token", {
            method: "POST",
            headers: { "Authorization": `Bearer ${localStorage.getItem("refresh_token")}` }
        });

        const refreshData = await refreshRes.json();
        if (refreshRes.ok && refreshData.access_token) {
            localStorage.setItem("access_token", refreshData.access_token);
            // Retry original request
            return authFetch(url, options);
        } else {
            alert("Session expired. Please log in again.");
            localStorage.removeItem("access_token");
            localStorage.removeItem("refresh_token");
            window.location.reload();
            return;
        }
    }

    return res;
}
