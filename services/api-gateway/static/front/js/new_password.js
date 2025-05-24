document.addEventListener("DOMContentLoaded", function () {
  const form = document.getElementById('newPasswordForm');
  const password = form.password;
  const passwordCheck = form.passwordCheck;

  const passwordError = document.getElementById('passwordError');
  const passwordCheckError = document.getElementById('passwordCheckError');

  // Проверка силы пароля
  function validatePassword() {
    const isValid = /^(?=.*[a-zA-Z])(?=.*\d).{8,}$/.test(password.value);
    if (!isValid) {
        password.classList.add('is-invalid');
        passwordError.textContent = 'Пароль должен быть не менее 8 символов и содержать латинские буквы и цифры.';
    } else {
        password.classList.remove('is-invalid');
        passwordError.textContent = '';
    }
  }

  // Проверка совпадения паролей
  function validatePasswordMatch() {
    if (passwordCheck.value === '') {
      passwordCheckError.textContent = '';
      return;
    }
    if (password.value !== passwordCheck.value) {
        passwordCheck.classList.add('is-invalid');
        passwordCheckError.textContent = 'Пароли не совпадают';
    } else {
        passwordCheck.classList.remove('is-invalid');
        passwordCheckError.textContent = '';
    }
  }

   password.addEventListener('input', () => {
    validatePassword();
    validatePasswordMatch();
  });

  passwordCheck.addEventListener('input', validatePasswordMatch);

   //  Отправка формы 
  form.addEventListener('submit', function (event) {
    event.preventDefault();

    validatePassword();
    validatePasswordMatch();

    if (
      passwordError.textContent ||
      passwordCheckError.textContent
    ) {
      return;
    }


    const formData = {
        login: "!!!",
        password: form.password.value
    };

    fetch('http://192.168.1.207:8080/api/new-password', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(formData)
    })
      .then(response => {
        if (response.ok) {
          alert('Пароль изменен.');
          form.reset();
        } else {
          alert('Ошибка изменения пароля.');
        }
      })
      .catch(error => {
        console.error('Ошибка:', error);
        alert('Произошла ошибка при отправке данных.');
      });
  });

  });