const loginView = document.getElementById("loginView");
const registerView = document.getElementById("registerView");

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

    const res = await fetch("/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password })
    });

    const data = await res.json();
    alert(data.error || data.message || "Login response received.");
});

document.getElementById("registerForm").addEventListener("submit", async (e) => {
    e.preventDefault();
    const username = document.getElementById("registerUsername").value;
    const email = document.getElementById("registerEmail").value;
    const password = document.getElementById("registerPassword").value;

    const res = await fetch("/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, email, password })
    });

    const data = await res.json();
    alert(data.message || data.error || "Registration response received.");
});
