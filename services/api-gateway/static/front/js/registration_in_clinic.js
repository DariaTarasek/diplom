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

  fetch('/api/admin-data')
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
    if (digits === '') {
      phoneInput.classList.remove('is-invalid');
      phoneError.textContent = '';
      return true;
    }
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


    const formData = {
      secondName: form.elements['secondName'].value,
      firstName: form.elements['firstName'].value,
      surname: form.elements['surname'].value,
      gender: form.elements['gender'].value,
      birthDate: form.elements['birthDate'].value,
      phone: form.elements['phone'].value,
      email: form.elements['email'].value
    };

    formData.phone = formData.phone.replace(/\D/g, '').trim()

    fetch('/api/register-in-clinic', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(formData)
    })
        .then(response => {
          if (response.ok) {
            alert('Регистрация прошла успешно!');
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
