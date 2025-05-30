document.getElementById('loginForm').addEventListener('submit', async function(e) {
    e.preventDefault();
    const login = document.getElementById('login').value.replace(/\D/g, '');
    const password = document.getElementById('password').value;

    const response = await fetch("/api/login", {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ login, password }),
        credentials: 'include' 
    });

    if (response.ok) {
        window.location.href = "/patient_account.html";
    } else {
        alert("Ошибка входа");
    }
});
