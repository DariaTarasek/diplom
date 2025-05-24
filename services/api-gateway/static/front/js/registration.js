document.addEventListener("DOMContentLoaded", function () {
  const form = document.getElementById('registrationForm');
  const password = form.password;
  const passwordCheck = form.password_check;
  const email = form.email;

  const passwordError = document.getElementById('passwordError');
  const passwordCheckError = document.getElementById('passwordCheckError');
  const emailError = document.getElementById('emailError');

  const firstName = form.firstName;
  const secondName = form.secondName;
  const firstNameError = document.getElementById('firstNameError');
  const secondNameError = document.getElementById('secondNameError');

  // Проверка даты 
  const dateInput = document.getElementById("birthDate");
  if (dateInput) {
    const today = new Date();
    const yyyy = today.getFullYear();
    const mm = String(today.getMonth() + 1).padStart(2, "0");
    const dd = String(today.getDate()).padStart(2, "0");

    const minDate = `${yyyy - 110}-${mm}-${dd}`;
    const maxDate = `${yyyy - 18}-${mm}-${String(today.getDate() - 1).padStart(2, "0")}`;
    dateInput.min = minDate;
    dateInput.max = maxDate;
  }

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

  // Проверка email
  function validateEmail() {
    const value = email.value.trim();
    if (value.length === 0) {
      emailError.textContent = '';
      return;
    }
    const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]{2,}$/;
    if (!emailPattern.test(value)) {
        email.classList.add('is-invalid');
        emailError.textContent = 'Некорректный email';
    } else {
        email.classList.remove('is-invalid');
        emailError.textContent = '';
    }
  }

  password.addEventListener('input', () => {
    validatePassword();
    validatePasswordMatch();
  });

  passwordCheck.addEventListener('input', validatePasswordMatch);
  email.addEventListener('input', validateEmail);

  firstName.addEventListener('input', () => {
    if (firstName.value.trim().length === 0) {
        firstNameError.textContent = 'Имя не может быть пустым';
         firstName.classList.add('is-invalid');
     } else {
        firstNameError.textContent = '';
        firstName.classList.remove('is-invalid');
        }
  });

  secondName.addEventListener('input', () => {
    if (secondName.value.trim().length === 0) {
        secondNameError.textContent = 'Фамилия не может быть пустой';
         secondName.classList.add('is-invalid');
     } else {
        secondNameError.textContent = '';
        secondName.classList.remove('is-invalid');
        }
  });

  

// Валидация номера телефона
    const phoneInput = document.getElementById('phone');
    const phoneError = document.getElementById('phoneError');

    
    const validatePhone = () => {
        const digits = phoneInput.value.replace(/\D/g, '');

        if (digits.length !== 11) {
            phoneInput.classList.add('is-invalid');
            phoneError.textContent = 'Введите полный 11-значный номер телефона.';
            return false;
        } else {
            phoneInput.classList.remove('is-invalid');
            phoneError.textContent = '';
           return true;
        }
    };

    phoneInput.addEventListener('input', validatePhone);


  // Подтверждение номера телефона 
  const requestCodeBtn = document.getElementById('requestCodeBtn');
  const verifyCodeBtn = document.getElementById('verifyCodeBtn');
  const codeSection = document.getElementById('codeSection');
  const codeInput = document.getElementById('smsCode');
  const codeMessage = document.getElementById('codeMessage');
  const phoneField = document.getElementById('phone');
  const changePhone = document.getElementById('changePhone');
  const resendCodeBtn = document.getElementById('resendCodeBtn');
  const changePhoneBtn = document.getElementById('changePhoneBtn');

  let phoneVerified = false;
  let verifiedPhoneNumber = '';

 let resendTimer; // Таймер для повторной отправки кода

  // Функция для повторной отправки кода
  function startResendTimer() {
    clearInterval(resendTimer);
    let seconds = 60;
    resendCodeBtn.disabled = true;
    resendCodeBtn.textContent = `Отправить код повторно (${seconds})`;

    resendTimer = setInterval(() => {
      seconds--;
      resendCodeBtn.textContent = `Отправить код повторно (${seconds})`;

      if (seconds <= 0) {
        clearInterval(resendTimer);
        resendCodeBtn.disabled = false;
        resendCodeBtn.textContent = 'Отправить код повторно';
      }
    }, 1000);
  }

  // Запросить код
  requestCodeBtn.addEventListener('click', () => {
     const phone = phoneField.value.trim();

    if (!phone) {
        codeMessage.textContent = 'Введите номер телефона.';
        codeMessage.classList.remove('text-success');
        codeMessage.classList.add('text-danger');
        return;
    }

    if (!validatePhone()) {
        codeMessage.textContent = 'Введите корректный номер телефона.';
        codeMessage.classList.remove('text-success');
        codeMessage.classList.add('text-danger');
        return;
    }

      requestCodeBtn.disabled = true;

    fetch('/api/request-code', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ phone })
    })
    .then(res => {
        if (res.ok) {
        phoneField.readOnly = true;
        verifiedPhoneNumber = phone;
        codeSection.style.display = 'block';
        changePhone.style.display = 'block';
        codeMessage.textContent = 'Код отправлен на номер.';
        codeMessage.classList.remove('text-danger');
        codeMessage.classList.add('text-success');
        startResendTimer();
        } else {
        requestCodeBtn.disabled = false;
        codeMessage.textContent = 'Ошибка при отправке кода.';
        codeMessage.classList.remove('text-success');
        codeMessage.classList.add('text-danger');
        }
    })
    .catch(() => {
         requestCodeBtn.disabled = false;
        codeMessage.textContent = 'Ошибка сети при отправке кода.';
        codeMessage.classList.remove('text-success');
        codeMessage.classList.add('text-danger');
    });
  });

  verifyCodeBtn.addEventListener('click', () => {
    const code = codeInput.value.trim();

    if (!code || !verifiedPhoneNumber) {
      codeMessage.textContent = 'Введите код';
      codeMessage.classList.remove('text-success');
      codeMessage.classList.add('text-danger');
      return;
    }

    fetch('/api/verify-code', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ phone: verifiedPhoneNumber, code })
    })
      .then(res => {
        if (res.ok) {
          codeMessage.classList.remove('text-danger');
          codeMessage.classList.add('text-success');
          phoneVerified = true;
          codeSection.style.display = 'none';         
          changePhone.style.display = 'none'; 
          requestCodeBtn.disabled = true;        
          phoneField.disabled = true;                
          const confirmedLabel = document.createElement('div'); 
          confirmedLabel.className = 'text-success mt-2';
          confirmedLabel.textContent = 'Телефон подтвержден!';
          phoneField.parentNode.appendChild(confirmedLabel); 

        } else {
          codeMessage.textContent = 'Неверный код.';
          codeMessage.classList.remove('text-success');
          codeMessage.classList.add('text-danger');
        }
      })
      .catch(() => {
        codeMessage.textContent = 'Ошибка проверки кода.';
        codeMessage.classList.remove('text-success');
        codeMessage.classList.add('text-danger');
      });
  });

 resendCodeBtn.addEventListener('click', () => {
    const phone = phoneField.value.trim();

    fetch('/api/request-code', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ phone })
    })
      .then(res => {
        if (res.ok) {
          startResendTimer(); // таймер для повторной отправки кода
        } else {
          codeMessage.textContent = 'Ошибка при повторной отправке кода.';
          codeMessage.classList.remove('text-success');
          codeMessage.classList.add('text-danger');
        }
      })
      .catch(() => {
        codeMessage.textContent = 'Ошибка сети при повторной отправке кода.';
        codeMessage.classList.remove('text-success');
        codeMessage.classList.add('text-danger');
      });
  });

  // Кнопка для изменения номера телефона
  changePhoneBtn.addEventListener('click', () => {
    clearInterval(resendTimer); 
    phoneField.readOnly = false;
    verifiedPhoneNumber = '';
    phoneVerified = false;
    codeSection.style.display = 'none';
    changePhone.style.display = 'none';
    codeInput.value = '';
    codeMessage.textContent = '';
    requestCodeBtn.disabled = false;
  });

  //  Отправка формы 
  form.addEventListener('submit', function (event) {
    event.preventDefault();

    const isFirstNameValid = firstName.value.trim().length > 0;
    const isSecondNameValid = secondName.value.trim().length > 0;


    validatePassword();
    validatePasswordMatch();
    validateEmail();

    if (
      passwordError.textContent ||
      passwordCheckError.textContent ||
      emailError.textContent || 
      !isFirstNameValid ||
      !isSecondNameValid
    ) {
      return;
    }

    if (!phoneVerified) {
      alert('Подтвердите номер телефона перед регистрацией.');
      return;
    }

    const formData = {
      secondName: form.secondName.value,
      firstName: form.firstName.value,
      surname: form.surname.value,
      gender: form.gender.value,
      birthDate: form.birthDate.value,
      phone: verifiedPhoneNumber,
      email: form.email.value,
      password: form.password.value
    };

    fetch('/api/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(formData)
    })
      .then(response => {
        if (response.ok) {
          alert('Регистрация прошла успешно!');
          phoneField.readOnly = false;
          phoneVerified = false;
          verifiedPhoneNumber = '';
          codeSection.style.display = 'none';
          changePhone.style.display = 'none';
          codeInput.value = '';
          codeMessage.textContent = '';
          requestCodeBtn.disabled = false;
          window.location.href = "/auth.html"
        } else {
          alert('Ошибка регистрации.');
        }
      })
      .catch(error => {
        console.error('Ошибка:', error);
        alert('Произошла ошибка при отправке данных.');
      });
  });
});