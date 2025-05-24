document.addEventListener("DOMContentLoaded", function () {
  const form = document.getElementById('registrationForm');
  const email = form.email;
  const emailError = document.getElementById('emailError');

  const firstName = form.firstName;
  const secondName = form.secondName;
  const firstNameError = document.getElementById('firstNameError');
  const secondNameError = document.getElementById('secondNameError');

  const button = document.getElementById('profile-button');
  const popover = document.getElementById('profile-popover');
  const popoverUsername = document.getElementById('popover-username');

  fetch('http://192.168.1.207:8080/api/admin-data')
    .then(response => {
      if (!response.ok) throw new Error("Ошибка при получении данных");
      return response.json();
    })
    .then(data => {
      const fullName = data.second_name + " " + data.first_name || "Неизвестный пользователь";
      button.textContent = fullName;
      popoverUsername.textContent = fullName;
    })
    .catch(error => {
      console.error("Ошибка при загрузке имени пользователя:", error);
      button.textContent = "Ошибка загрузки";
      popoverUsername.textContent = "Ошибка загрузки";
    });

  
  button.addEventListener('click', () => {
    popover.classList.toggle('d-none');
  });


  document.addEventListener('click', (event) => {
    const profileArea = document.getElementById('admin-profile');
    if (!profileArea.contains(event.target)) {
      popover.classList.add('d-none');
    }
  });

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

  email.addEventListener('input', validateEmail);

  
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
  let resendTimer;

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

    fetch('http://192.168.1.207:8080/api/request-code', {
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

    fetch('http://192.168.1.207:8080/api/verify-code', {
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

    fetch('http://192.168.1.207:8080/api/request-code', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ phone })
    })
      .then(res => {
        if (res.ok) {
          startResendTimer();
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

  form.addEventListener('submit', function (event) {
    event.preventDefault();

    const isFirstNameValid = firstName.value.trim().length > 0;
    const isSecondNameValid = secondName.value.trim().length > 0;

    validateEmail();

    if (emailError.textContent) {
      return;
    }

    if (!isFirstNameValid || !isSecondNameValid) return;

    if (!phoneVerified) {
      alert('Подтвердите номер телефона перед регистрацией.');
      return;
    }



    const formData = {
      second_name: form.elements['secondName'].value,
      first_name: form.elements['firstName'].value,
      surname: form.elements['surname'].value,
      gender: form.elements['gender'].value,
      birthDate: form.elements['birthDate'].value,
      phone: verifiedPhoneNumber,
      email: form.elements['email'].value
    };

    fetch('http://192.168.1.207:8080/api/register-in-clinic', {
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
          window.location.href = "/admins_patient_list.html"
        } else {
          response.text().then(text => {
            alert('Ошибка регистрации: ' + text);
          });
        }
      })
      .catch(error => {
        console.error('Ошибка:', error);
        alert('Произошла ошибка при отправке данных.');
      });
  });
});
