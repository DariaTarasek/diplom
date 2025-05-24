document.addEventListener('DOMContentLoaded', () => {
  const form = document.getElementById('recoveryForm');
  const input = document.getElementById('login');
  const loginError = document.getElementById('loginError');

  form.addEventListener('submit', async (e) => {
    e.preventDefault();

    const rawValue = input.value.trim();

    // Проверка номера (нач. с 7, 11 символов)
    const digits = phoneInput.value.replace(/\D/g, '');
    const isPhone = digits.length === 11 && /^[7]/.test(digits);

    if (!isPhone) {
        loginError.textContent = 'Введите корректный номер телефона'
      return;
    } else {
        loginError.textContent = ''
    }

    const data = {
      contact: rawValue,
      type: 'phone',
    };

    try {
      const response = await fetch('/api/password-recovery', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(data),
      });

      if (response.ok) {
        alert('На указанный номер отправлен код восстановления пароля.');
        window.location.href = "/index.html"
      } else {
        alert('Ошибка на сервере. Попробуйте позже.');
      }
    } catch (err) {
      console.error(err);
      alert('Ошибка сети.');
    }
  });
});
