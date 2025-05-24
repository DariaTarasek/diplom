document.addEventListener('DOMContentLoaded', () => {
  const form = document.getElementById('recoveryForm');
  const input = document.getElementById('login');
  const loginError = document.getElementById('loginError');

  form.addEventListener('submit', async (e) => {
    e.preventDefault();

    const rawValue = input.value.trim();

    // Проверка Email
    const isEmail = /^[^\s@]+@[^\s@]+\.[^\s@]{2,}$/.test(rawValue);

    if (!isEmail) {
        loginError.textContent = 'Введите корректный email'
      return;
    } else {
        loginError.textContent = ''
    }

    const data = {
      contact: rawValue,
      type: 'email',
    };

    try {
      const response = await fetch('/api/emaployee-password-recovery', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(data),
      });

      if (response.ok) {
        alert('На указанный email отправлен код восстановления пароля.');
        window.location.href = "/auth_doc.html"
      } else {
        alert('Ошибка на сервере. Попробуйте позже.');
      }
    } catch (err) {
      console.error(err);
      alert('Ошибка сети.');
    }
  });
});
