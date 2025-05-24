document.addEventListener("DOMContentLoaded", function () {
  const form = document.getElementById('registrationForm');
  const email = form.email;

  const emailError = document.getElementById('emailError');
  const specError = document.getElementById('specError');

  const firstName = form.firstName
  const secondName = form.secondName;
  const firstNameError = document.getElementById('firstNameError');
  const secondNameError = document.getElementById('secondNameError');

  const roleSelect = document.getElementById('role');
  const doctorFields = document.getElementById('doctorFields');

  const container = document.getElementById("specializations-container");

  // Загрузка специализаций с сервера
  fetch('/api/specialties')
    .then(response => response.json())
    .then(data => {
      data.forEach(spec => {
        const label = document.createElement('label');
        label.style.display = 'block';

        const checkbox = document.createElement('input');
        checkbox.type = 'checkbox';
        checkbox.name = 'specializations';
        checkbox.value = spec.id;

        label.appendChild(checkbox);
        label.appendChild(document.createTextNode(` ${spec.name}`));
        container.appendChild(label);
      });
    })
    .catch(error => {
      console.error('Ошибка загрузки специализаций:', error);
    });

  // Проверка роли
  roleSelect.addEventListener('change', function () {
    if (parseInt(this.value) === 3) {
      doctorFields.style.display = 'block';
    } else {
      doctorFields.style.display = 'none';
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
            return
        } else {
            phoneInput.classList.remove('is-invalid');
            phoneError.textContent = '';
           return
        }
    };

    phoneInput.addEventListener('input', validatePhone);

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


  container.addEventListener('change', () => {
    const selectedSpecs = container.querySelectorAll('input[name="specializations"]:checked');
    if (selectedSpecs.length === 0) {
        specError.textContent = 'Выберите хотя бы одну специализацию';
    } else {
        specError.textContent = '';
    }
    });

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

  form.addEventListener('submit', function (event) {
    event.preventDefault();

    const isFirstNameValid = firstName.value.trim().length > 0;
    const isSecondNameValid = secondName.value.trim().length > 0;

    validateEmail();
    validatePhone();

    if (
      emailError.textContent || 
      phoneError.textContent || 
      !isFirstNameValid ||
      !isSecondNameValid
    ) {
      return;
    }

  const selectedSpecs = Array.from(
    document.querySelectorAll('input[name="specializations"]:checked')
  ).map(cb => cb.value);

  
  if (parseInt(roleSelect.value) === 3 && selectedSpecs.length === 0) {
    alert('Выберите хотя бы одну специализацию.');
    return;
  }

    const formData = {
      secondName: form.secondName.value,
      firstName: form.firstName.value,
      surname: form.surname.value,
      gender: form.gender.value,
      phone: form.phone.value.replace(/\D/g, ''),
      email: form.email.value,
      role: parseInt(form.role.value),
    };
  console.log(formData.role)

    if (formData.role === 3) {
        formData.education = form.education.value;
        formData.experience = parseInt(form.experience.value);
        formData.specializations = selectedSpecs.map(id => parseInt(id, 10));
    } 

    fetch('/api/employee-register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(formData)
    })
      .then(response => {
        if (response.ok) {
            alert('Регистрация прошла успешно!');
            if (formData.role === 3) {
                window.location.href = "/admins_doctor_list.html"
            }
            else {
            window.location.href = "/admins_admin_list.html"
            }
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