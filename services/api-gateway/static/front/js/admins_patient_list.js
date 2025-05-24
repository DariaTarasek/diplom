const { createApp, ref, reactive, onMounted, watch, computed } = Vue;

createApp({
  setup() {
    const patients = ref([]);
    const search = ref('');
    const form = reactive({
      id: null,
      second_name: '',
      first_name: '',
      birthDate: '',
      gender: '',
      surname: '',
      phone: '',
      email: ''
    });

    const admin = ref({
    second_name: '',
    first_name: ''
    })
    const fullName = computed(() => {
  return [admin.value.first_name, admin.value.second_name].filter(Boolean).join(' ');
    });
    const isPopoverVisible = ref(false);

    let modalEdit = null;
    let modalLogin = null;
    let modalChoice = null;
    const selectedPatient = ref(null);


    const emailError = ref('');
    const phoneError = ref('');
    const selectedPatientId = ref(null);
    const isSwitchingModals = ref(false);

     const firstNameError = ref('');
    const secondNameError = ref('');

    const phoneVerified = ref(false);
    const verifiedPhoneNumber = ref('');
    const smsCode = ref('');
    const codeMessage = ref('');
    const codeSectionVisible = ref(false);
    const resendTimer = ref(null);
    const resendButtonDisabled = ref(true);
    const resendButtonText = ref('Отправить код повторно');

    const requestButtonDisabled = ref(false);
    const phoneSectionDisabled = ref(false);
    const codeSended = ref(false);

    const userRole = ref(null);

        function togglePopover() {
            isPopoverVisible.value = !isPopoverVisible.value;
        }

    function handleClickOutside(event) {
            const popover = document.getElementById('admin-profile');
            if (popover && !popover.contains(event.target)) {
                isPopoverVisible.value = false;
            }
        }

    function formatDate(dateStr) {
        if (!dateStr) return '';
        const parts = dateStr.split('.');
        if (parts.length !== 3) return '';
        const [day, month, year] = parts;
        return `${day.padStart(2, '0')}.${month.padStart(2, '0')}.${year}`;
        }

       function formatPhone(phone) {
            const digits = phone.replace(/\D/g, '');
            if (digits.length !== 11 || (!digits.startsWith('7') && !digits.startsWith('8'))) return phone;

            const code = digits.slice(1, 4);
            const part1 = digits.slice(4, 7);
            const part2 = digits.slice(7, 9);
            const part3 = digits.slice(9, 11);

            return `+7 (${code}) ${part1}-${part2}-${part3}`;
            }


    async function fetchUserRole() {
        try {
            const res = await fetch('http://192.168.1.207:8080/api/user-role'); 
            if (res.ok) {
            const data = await res.json();
            userRole.value = data.role;
            } else {
            console.error('Ошибка получения роли пользователя');
            userRole.value = 'admin';
            }
        } catch (err) {
            console.error('Ошибка запроса роли:', err);
            userRole.value = 'admin'; 
        }
        }


    async function loadPatients() {
      const res = await fetch(`http://192.168.1.207:8080/api/patients?search=${encodeURIComponent(search.value)}`);
      patients.value = await res.json();
    }

    function onRowClick(p) {
        selectedPatient.value = p;
        selectedPatientId.value = p.id;

        const modalElement = document.getElementById('actionsModal');
        if (modalElement) {
            modalChoice = new bootstrap.Modal(modalElement);
            modalChoice.show();
        }
    }

    function updateRequestButtonState() {
        const digits = form.phone.replace(/\D/g, '');
  const isValid = digits.length === 11 && digits.startsWith('7');
  requestButtonDisabled.value = !isValid;
        }


    function validatePhone(phone) {
  const digits = phone.replace(/\D/g, '');
  return digits.length === 11 && digits.startsWith('7');
}

    async function requestCode() {
    const phoneInput = document.getElementById('phoneLogin');
    if (phoneInput?.maskRef) {
        phoneInput.maskRef.updateValue(); // обновляем внутреннее состояние маски
        form.phone = phoneInput.maskRef.value; // обновляем поле
    }
    const phone = form.phone.replace(/\D/g, '');
    console.log('phone dlya confirm:' + phone)

    phoneError.value = ''; // сбрасываем ошибку, если всё ок

    // Блокируем поле и кнопку
    resendButtonDisabled.value = true;
    requestButtonDisabled.value = true;
    phoneSectionDisabled.value = true;
    codeSended.value = true;

    const res = await fetch('http://192.168.1.207:8080/api/request-code', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ phone })
    });

    if (res.ok) {
        codeSectionVisible.value = true;
        verifiedPhoneNumber.value = phone;
        phoneVerified.value = false;
        codeMessage.value = 'Код отправлен на номер.';
        startResendTimer();
    } else {
        codeMessage.value = 'Ошибка при отправке кода.';
        // Разблокируем поле в случае ошибки
        if (phoneInput) phoneInput.disabled = false;
    }
}


    async function verifyCode() {
    if (!smsCode.value || !verifiedPhoneNumber.value) {
        codeMessage.value = 'Введите код';
        return;
    }

    const res = await fetch('http://192.168.1.207:8080/api/verify-code', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ phone: verifiedPhoneNumber.value, code: smsCode.value })
    });

    if (res.ok) {
        phoneVerified.value = true;
        codeSectionVisible.value = false;
        codeMessage.value = 'Телефон подтвержден!';
    } else {
        codeMessage.value = 'Неверный код.';
    }
    }

    function startResendTimer() {
        clearInterval(resendTimer.value);
        let seconds = 60;
        resendButtonDisabled.value = true;
        resendButtonText.value = `Отправить код повторно (${seconds})`;

        resendTimer.value = setInterval(() => {
            seconds--;
            resendButtonText.value = `Отправить код повторно (${seconds})`;

            if (seconds <= 0) {
            clearInterval(resendTimer.value);
            resendButtonDisabled.value = false;
            resendButtonText.value = 'Отправить код повторно';
            }
        }, 1000);
        }

        function resetPhoneConfirmation() {
            clearInterval(resendTimer.value);
            phoneVerified.value = false;
            verifiedPhoneNumber.value = '';
            smsCode.value = '';
            codeMessage.value = '';
            codeSectionVisible.value = false;
            resendButtonDisabled.value = true;
            resendButtonText.value = 'Отправить код повторно';

            requestButtonDisabled.value = false;
            phoneSectionDisabled.value = false;
            codeSended.value = false;
    }



    function openEdit() {
        isSwitchingModals.value = true;
        modalChoice?.hide();

        Object.assign(form, selectedPatient.value); // копируем данные пациента

         form.gender = selectedPatient.value.gender === 'м' ? 'male' : 'female';
         console.log(selectedPatient.value.birthDate);

         //  Преобразование даты из ДД.ММ.ГГГГ в ГГГГ-MM-ДД
        if (selectedPatient.value.birthDate) {
            const parts = selectedPatient.value.birthDate.split('.');
            if (parts.length === 3) {
                const [day, month, year] = parts;
                form.birthDate = `${year}-${month.padStart(2, '0')}-${day.padStart(2, '0')}`;
            } else {
                form.birthDate = '';
            }
        }


        if (!modalEdit) {
            modalEdit = new bootstrap.Modal(document.getElementById('editModal'));
        }
        modalEdit.show();
    }

    function isSelected(p) {
        return selectedPatientId.value === p.id;
    }



    function openLoginChange() {
        isSwitchingModals.value = true;
      modalChoice.hide();

      if (!modalLogin) {
        modalLogin = new bootstrap.Modal(document.getElementById('changeLoginModal'));
      }
        resetPhoneConfirmation();

      modalLogin.show();

      // Маска логина
      setTimeout(() => {
        const phoneInput = document.getElementById('phoneLogin');
        if (phoneInput) {
          if (phoneInput.maskRef?.destroy) {
            phoneInput.maskRef.destroy();
          }

          let digits = (form.phone || '').replace(/\D/g, '');
          digits = digits.slice(1);

          const mask = IMask(phoneInput, {
            mask: '+7 (000) 000-00-00',
            lazy: false,
            overwrite: true
          });

          phoneInput.maskRef = mask;

          mask.on('accept', () => {
            const digits = mask.value.replace(/\D/g, '');
            const valid = digits.length === 11 && digits.startsWith('7');
            phoneError.value = digits && !valid ? 'Неверный формат телефона' : '';
            form.phone = mask.value;
            updateRequestButtonState();
          });

          mask.unmaskedValue = digits;
          form.phone = mask.value;
          updateRequestButtonState();
        }
      }, 300);
    }

    function onModalHidden() {
        if (isSwitchingModals.value) {
            isSwitchingModals.value = false;
            return;
        }
        selectedPatientId.value = null;
        selectedPatient.value = null;
    }


    async function savePatient() {
      const emailValid = !emailError.value;

      if (!emailValid) {
        return;
      }

      await fetch(`http://192.168.1.207:8080/api/patients/${form.id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(form)
      });

      modalEdit.hide();
      await loadPatients();
    }

    async function saveLogin() {
        const phoneValid = !phoneError.value && form.phone && phoneVerified.value;

        if (!phoneValid) {
            codeMessage.value = 'Номер не подтвержден.';
            return;
        }

        await fetch(`http://192.168.1.207:8080/api/patients/${form.id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ phone: form.phone })
        });

      modalLogin.hide();
      selectedPatientId.value = null;
        selectedPatient.value = null;
      await loadPatients();
    }

        async function fetchAdminData() {
  try {
    const res = await fetch('http://192.168.1.207:8080/api/admin-data');
    const data = await res.json();

    admin.value.first_name = data.first_name || '';
    admin.value.second_name = data.second_name || '';

  } catch (err) {
    console.error('Ошибка при загрузке данных:', err);
  }
}

    async function deletePatient() {
  if (!selectedPatient.value) return;

  const confirmed = confirm(`Вы уверены, что хотите удалить пользователя: ${selectedPatient.value.first_name} ${selectedPatient.value.second_name}?`);
  if (!confirmed) return;

  try {
    const res = await fetch(`http://192.168.1.207:8080/api/patients/${selectedPatient.value.id}`, {
      method: 'DELETE'
    });

    if (res.ok) {
      modalChoice?.hide();
      selectedPatient.value = null;
      selectedPatientId.value = null;
      await loadPatients();
      alert('Пользователь удалён.');
    } else {
      alert('Ошибка при удалении пользователя.');
    }
  } catch (e) {
    console.error(e);
    alert('Произошла ошибка.');
  }
}


    watch(() => form.email, (val) => {
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
      emailError.value = val && !emailRegex.test(val) ? 'Неверный формат email' : '';
    });

    watch(() => form.first_name, (val) => {
        if (val.trim().length === 0) firstNameError.value = 'Имя не может быть пустым'
        else firstNameError.value = '';
    });

    watch(() => form.second_name, (val) => {
        if (val.trim().length === 0) secondNameError.value = 'Фамилия не может быть пустой'
        else secondNameError.value = '';
    });

    onMounted(async () => {
  const editEl = document.getElementById('editModal');
  const loginEl = document.getElementById('changeLoginModal');
  const actionsEl = document.getElementById('actionsModal');

  if (editEl) {
    editEl.addEventListener('hidden.bs.modal', onModalHidden);
  }
  if (loginEl) {
    loginEl.addEventListener('hidden.bs.modal', onModalHidden);
  }
  if (actionsEl) {
    actionsEl.addEventListener('hidden.bs.modal', onModalHidden);
  }

  loadPatients();
  await fetchUserRole();
  document.addEventListener('click', handleClickOutside);
    await fetchAdminData();
});

    return {
      patients,
      search,
      form,
      phoneError,
      emailError,
      formatPhone,
      loadPatients,
      onRowClick,
      openEdit,
      openLoginChange,
      savePatient,
      saveLogin,
      selectedPatient,
      selectedPatientId,
      isSelected,
      onModalHidden,
      formatDate,
      phoneVerified,
      verifiedPhoneNumber,
      smsCode,
      codeMessage,
      codeSectionVisible,
      resendButtonDisabled,
      resendButtonText,
      requestCode,
      verifyCode,
      resetPhoneConfirmation,
      requestButtonDisabled,
      phoneSectionDisabled,
      codeSended,
      userRole,
      fetchUserRole,
      deletePatient, 
      firstNameError,
      secondNameError,
      isPopoverVisible,
      fullName
    };
  }
}).mount('#app');
