document.addEventListener("DOMContentLoaded", function () {
    const form = document.getElementById("loginForm");
    const emailInput = document.getElementById("login");
    const emailError = document.getElementById("emailError");

    form.addEventListener("submit", async function (e) {
        e.preventDefault();
        const emailValue = emailInput.value.trim();
        const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]{2,}$/;

        if (!emailPattern.test(emailValue)) {
            emailError.textContent = 'Некорректный email';
            return;
        } else {
            emailError.textContent = '';
        }

        const password = document.getElementById("password").value;
        const login = emailValue

        try {
            const response = await fetch("/api/login", {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ login, password }),
                credentials: 'same-origin'
            });

            if (response.ok) {
                const data = await response.json();
                const role = data.role; // предполагается, что в ответе есть поле `role`

                if (role === "doctor") {
                    window.location.href = "doctor_account.html";
                } else if (role === "admin" || role === "superadmin") {
                    window.location.href = "administrator_account.html";
                }
            } else {
                alert("Ошибка входа");
            }
        } catch (error) {
            console.error("Ошибка при выполнении запроса:", error);
            alert("Произошла ошибка, попробуйте позже.");
        }
    });
});
