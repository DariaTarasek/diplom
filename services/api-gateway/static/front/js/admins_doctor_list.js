const { createApp, ref, reactive, watch, onMounted, computed } = Vue;


createApp({
  setup() {
    const search = ref('');
    const filters = reactive({ specialty: '' });
    const staff = ref([]);
    const specialties = ref([]);
    const form = reactive({
      id: '',
      second_name: '',
      first_name: '',
      surname: '',
      email: '',
      specialty: [],
      gender: '',
      phone: '',
      experience: '',
      education: ''
    });

    const admin = ref({
    second_name: '',
    first_name: ''
    })
    const fullName = computed(() => {
  return [admin.value.first_name, admin.value.second_name].filter(Boolean).join(' ');
    });
    const isPopoverVisible = ref(false);

    const emailError = ref('');
    const phoneError = ref('');
    const specError = ref('');
    const selectedDoctor = ref(null);
    const selectedDoctorId = ref(null);

    const firstNameError = ref('');
    const secondNameError = ref('');

    const isSwitchingModals = ref(false);

    let modalEdit = null;
    let modalLogin = null;
    let modalChoice = null;


    function togglePopover() {
            isPopoverVisible.value = !isPopoverVisible.value;
        }

    function handleClickOutside(event) {
            const popover = document.getElementById('admin-profile');
            if (popover && !popover.contains(event.target)) {
                isPopoverVisible.value = false;
            }
        }

    async function loadSpecialties() {
      specialties.value = await (await fetch('http://192.168.1.207:8080/api/admin/specialties')).json();
    }

    async function loadStaff() {
      const params = new URLSearchParams();
      if (search.value) params.append('search', search.value);
      if (filters.specialty) params.append('specialty', filters.specialty);
    

      const res = await fetch(`http://192.168.1.207:8080/api/staff-doctors?${params}`);
      staff.value = await res.json();
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

    function getSpecialtyName(id) {
        return specialties.value.find(s => s.id === id)?.name || '-';
    }

    function getSpecialtyNames(ids) {
        if (!Array.isArray(ids)) ids = [ids];
        return specialties.value
            .filter(s => ids.includes(s.id))
            .map(s => s.name)
            .join(', ') || '-';
}

    loadStaff();

       function onRowClick(s) {
            selectedDoctor.value = s;
            selectedDoctorId.value = s.id;

            const modalElement = document.getElementById('actionsModal');
            if (modalElement) {
                modalChoice = new bootstrap.Modal(modalElement);
                modalChoice.show();
            }
        }

  function openEdit(s) {
        isSwitchingModals.value = true;
            modalChoice?.hide();
            Object.assign(form, s);   
            form.gender = s.gender === 'м' ? 'male' : 'female';
            form.specialty = Array.isArray(s.specialty)
            ? s.specialty
            : String(s.specialty).split(',').map(Number);
            emailError.value = '';
            phoneError.value = '';
            specError.value = '';
            if (!modalEdit) {
                modalEdit = new bootstrap.Modal(document.getElementById('editModal'));
            }
            modalEdit.show();

            setTimeout(() => {
                const phoneInput = document.getElementById('phoneLogin');
                if (phoneInput) {
                // Удаляем старую маску, если была
                if (phoneInput.maskRef?.destroy) {
                    phoneInput.maskRef.destroy();
                }

                // Нормализация номера
                let digits = (form.phone || '').replace(/\D/g, '');

                digits = digits.slice(1);

                // Создаем маску
                const mask = IMask(phoneInput, {
                    mask: '+7 (000) 000-00-00',
                    lazy: false,
                    overwrite: true
                });

                // Привязываем маску к элементу
                phoneInput.maskRef = mask;

                mask.on('accept', () => {
                    const digits = mask.value.replace(/\D/g, ''); // Удаляем все кроме цифр
                    const valid = digits.length === 11 && digits.startsWith('7');
                    phoneError.value = digits && !valid ? 'Неверный формат телефона' : '';
                    form.phone = mask.value;
                });

                // Устанавливаем значение без потерь
                mask.unmaskedValue = digits;
                form.phone = mask.value;
                }
            }, 300);
        }

    function saveStaff() {
      const phoneValid = !phoneError.value && form.phone;

      if (!phoneValid) {
        return;
      }

      if (specError.value) {
        return;
  }

    secondNameError.value = '';
    firstNameError.value = '';

    const isSecondNameEmpty = !form.second_name || form.second_name.trim() === '';
    const isFirstNameEmpty = !form.first_name || form.first_name.trim() === '';

    if (isSecondNameEmpty) {
        secondNameError.value = 'Фамилия не может быть пустой';
    }
    if (isFirstNameEmpty) {
        firstNameError.value = 'Имя не может быть пустым';
    }


    if (isSecondNameEmpty || isFirstNameEmpty) {
        return;
    }

      const payload = {
        id: form.id,
        first_name: form.first_name,
        second_name: form.second_name,
        surname: form.surname,
        phone: form.phone,
        email: form.email,
        specialty: Array.isArray(form.specialty) ? form.specialty : [form.specialty],
        experience: form.experience,
        education:form.education
      };

      fetch('http://192.168.1.207:8080/api/save-staff', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      })
        .then(response => response.json())
        .then(data => {
          if (data.success) {
            loadStaff();
            if (modalEdit) {
                modalEdit.hide();
                selectedDoctorId.value = null;
                selectedDoctor.value = null;
            }
          } else {
            alert('Ошибка при сохранении');
          }
        });
    }

     function isSelected(s) {
        return selectedDoctorId.value === s.id;
    }

     function openLoginChange() {
        isSwitchingModals.value = true;
        modalChoice.hide();

        if (!modalLogin) {
            modalLogin = new bootstrap.Modal(document.getElementById('changeLoginModal'));
        }

      modalLogin.show();
    }

     function onModalHidden() {
        if (isSwitchingModals.value) {
            isSwitchingModals.value = false;
            return;
        }
        selectedDoctorId.value = null;
        selectedDoctor.value = null;
    }

     async function saveLogin() {
        const emailValid = !emailError.value && form.email;

        if (!emailValid || form.email.length === 0) {
            emailError.value = 'Некорректный email';
            return;
        }

        await fetch(`http://192.168.1.207:8080/api/doctors/${form.id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ phone: form.phone })
        });

      modalLogin.hide();
      selectedDoctorId.value = null;
      selectedDoctor.value = null;
      await loadStaff();
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

      async function deleteDoctor() {
            if (!selectedDoctor.value) return;

            const confirmed = confirm(`Вы уверены, что хотите удалить пользователя: ${selectedDoctor.value.first_name} ${selectedDoctor.value.second_name}?`);
            if (!confirmed) return;

            try {
                const res = await fetch(`http://192.168.1.207:8080/api/doctors/${selectedDoctor.value.id}`, {
                method: 'DELETE'
                });

                if (res.ok) {
                modalChoice?.hide();
                await loadStaff();
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

     watch(() => form.specialty, (val) => {
      if (val.length === 0) {
        specError.value = "Выберите хотя бы одну специализацию";
        } else {
            specError.value = '';
        }
    });

    watch(() => filters.specialty, () => {
        loadStaff();
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
            loadSpecialties();
            loadStaff();
            document.addEventListener('click', handleClickOutside);
        await fetchAdminData();
    });


    return {
      search,
      filters,
      staff,
      specialties,
      form,
      emailError,
      phoneError,
      specError,
      firstNameError,
      secondNameError,
      loadStaff,
      getSpecialtyName,
      getSpecialtyNames,
      openEdit,
      saveStaff,
      onRowClick,
      openLoginChange,
      saveLogin,
      selectedDoctor,
      selectedDoctorId,
      isSelected,
      onModalHidden,
      deleteDoctor, 
      formatPhone,
        isPopoverVisible,
      fullName
    };
  }
}).mount('#app');
