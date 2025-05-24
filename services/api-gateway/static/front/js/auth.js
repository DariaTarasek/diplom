document.getElementById('loginForm').addEventListener('submit', async function(e) {
    e.preventDefault();
    const phone = document.getElementById('login').value.replace(/\D/g, '');
    const password = document.getElementById('password').value;

    const response = await fetch("http://192.168.1.207:8080/api/login", {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ phone, password }),
        credentials: 'include' 
    });

    if (response.ok) {
        window.location.href = "/patient_account.html";
    } else {
        alert("Ошибка входа");
    }
});
