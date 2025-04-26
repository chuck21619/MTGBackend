document.getElementById('login-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;
  
    const response = await fetch('https://goanddocker.onrender.com/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, password })
    });
  
    const messageEl = document.getElementById('message');
    if (response.ok) {
      const data = await response.json();
      messageEl.style.color = 'green';
      messageEl.textContent = 'Login successful!';
      // localStorage.setItem('token', data.token); // if needed later
    } else {
      const error = await response.text();
      messageEl.style.color = 'red';
      messageEl.textContent = `Login failed: ${error}`;
    }
    console.log(messageEl.textContent);
  });
  