document.addEventListener("DOMContentLoaded", function () {
  const form = document.getElementById("loginForm");
  const emailInput = document.getElementById("login");
  const emailError = document.getElementById("emailError");

  form.addEventListener("submit", async function (e) {
    e.preventDefault(); 
    const emailValue = emailInput.value.trim();
    const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]{2,}$/;

    if (!emailPattern.test(emailValue)) {
        emailError.textContent = 'Некорректный email'
        return;
      } else {
        emailError.textContent = '';
      }
    
    const password = document.getElementById("password").value;

       const response = await fetch("http://192.168.1.207:8080/api/login", {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ emailValue, password }),
        credentials: 'same-origin' 
    });

    if (response.ok) {
        window.location.href = "#";
    } else {
        alert("Ошибка входа");
    }
  });

  
});
