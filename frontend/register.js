document.getElementById("registerForm").addEventListener("submit", async (e) => {
    e.preventDefault();
  
    const username = document.getElementById("username").value;
    const email = document.getElementById("email").value;
    const password = document.getElementById("password").value;
  
    const response = await fetch("/register", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({ username, email, password })
    });
  
    const messageEl = document.getElementById("registerMessage");
    if (response.ok) {
      messageEl.textContent = "Registration successful!";
    } else {
      const errorText = await response.text();
      messageEl.textContent = `Error: ${errorText}`;
    }
  });
  