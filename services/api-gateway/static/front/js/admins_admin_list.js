const { createApp, ref, reactive, watch, onMounted, computed } = Vue;


createApp({
  setup() {
    const search = ref('');
    const filters = reactive({ role: '' });
    const staff = ref([]);
    const roles = ref([]);
    const form = reactive({
      id: '',
      second_name: '',
      first_name: '',
      surname: '',
      email: '',
      role: '',
      gender: '',
      phone: ''
    });

    const admin = ref({
        second_name: '',
        first_name: ''
    })
    const fullName = computed(() => {
  return [admin.value.first_name, admin.value.second_name].filter(Boolean).join(' ');
    });


    const emailError = ref('');
    const phoneError = ref('');

    const userRole = ref(null);
    const selectedAdmin = ref(null);
    const selectedAdminId = ref(null);

    const firstNameError = ref('');
    const secondNameError = ref('');

    const isSwitchingModals = ref(false);
    const isPopoverVisible = ref(false);

    let modalEdit = null;
    let modalLogin = null;
    let modalChoice = null;

    function togglePopover() {
            isPopoverVisible = !isPopoverVisible;
        }

    function handleClickOutside(event) {
            const popover = document.getElementById('admin-profile');
            if (popover && !popover.contains(event.target)) {
                isPopoverVisible = false;
            }
        }


      async function loadRoles() {
      roles.value = await (await fetch('/api/roles')).json();
    }


    async function loadStaff() {
      const params = new URLSearchParams();
      if (search.value) params.append('search', search.value);
      if (filters.role) params.append('role', filters.role);

      // ?${params}
      const res = await fetch(`/api/staff-admins`);
      staff.value = await res.json();
    }

    const filteredStaff = computed(() => {
          const searchText = search.value.toLowerCase().trim();
          const selectedRole = filters.role;

          return staff.value.filter(s => {
              const fullName = `${s.second_name} ${s.first_name} ${s.surname}`.toLowerCase();
              const matchesSearch = !searchText || fullName.includes(searchText) || s.email.toLowerCase().includes(searchText);
              const matchesRole = !selectedRole || s.role === selectedRole;
              return matchesSearch && matchesRole;
          });
      });

    function formatPhone(phone) {
            const digits = phone.replace(/\D/g, '');
            if (digits.length !== 11 || (!digits.startsWith('7') && !digits.startsWith('8'))) return phone;

            const code = digits.slice(1, 4);
            const part1 = digits.slice(4, 7);
            const part2 = digits.slice(7, 9);
            const part3 = digits.slice(9, 11);

            return `+7 (${code}) ${part1}-${part2}-${part3}`;
            }

    function getRoleName(id) {
      return roles.value.find(r => r.id === id)?.name || id;
    }

    function onRowClick(s) {

         if (userRole.value !== 'superadmin') return;
        selectedAdmin.value = s;
        selectedAdminId.value = s.id;

        const modalElement = document.getElementById('actionsModal');
        if (modalElement) {
            modalChoice = new bootstrap.Modal(modalElement);
            modalChoice.show();
        }
    }

    function onRoleChange() {
      filteredStaff();
    }


  function openEdit(s) {
        if (userRole.value !== 'superadmin') return;
        isSwitchingModals.value = true;
        modalChoice?.hide();
        Object.assign(form, s);
        
         form.gender = s.gender === 'м' ? 'м' : 'ж';
        emailError.value = '';
        phoneError.value = '';
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

      const payload = {
        id: form.id,
        first_name: form.first_name,
        second_name: form.second_name,
        surname: form.surname,
        phone: form.phone.replace(/\D/g, '').trim(),
        email: form.email,
        role: form.role,
          gender: form.gender
      };

      fetch('/api/save-admin', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      })
    .then(response => {
            if (response.ok) {
                loadStaff();
                if (modalEdit) {
                    modalEdit.hide();
                    selectedAdminId.value = null;
                    selectedAdmin.value = null;
                }
            } else {
                alert('Ошибка при сохранении (код ' + response.status + ')');
            }
        })
    }

     function isSelected(s) {
        return selectedAdminId.value === s.id;
    }

     function openLoginChange() {
        if (userRole.value !== 'superadmin') return;
        isSwitchingModals.value = true;
        modalChoice.hide();

        if (!modalLogin) {
            modalLogin = new bootstrap.Modal(document.getElementById('changeLoginModal'));
        }
         if (selectedAdmin.value) {
             form.id = selectedAdmin.value.id;
             form.email = selectedAdmin.value.email;
         }

      modalLogin.show();
    }

     function onModalHidden() {
        if (isSwitchingModals.value) {
            isSwitchingModals.value = false;
            return;
        }
        selectedAdminId.value = null;
        selectedAdmin.value = null;
    }

     async function saveLogin() {
        const emailValid = !emailError.value && form.email;

        if (!emailValid) {
            codeMessage.value = 'Некорректный email';
            return;
        }

        await fetch(`/api/admins-login/${form.id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email: form.email })
        });

      modalLogin.hide();
      selectedAdminId.value = null;
      selectedAdmin.value = null;
      await loadStaff();
    }

    async function fetchAdminData() {
  try {
    const res = await fetch('/api/admin-data');
    const data = await res.json();

    admin.value.first_name = data.first_name || '';
    admin.value.second_name = data.second_name || '';
    userRole.value = data.role || '';

  } catch (err) {
    console.error('Ошибка при загрузке данных:', err);
  }
}

      async function deleteAdmin() {
            if (!selectedAdmin.value) return;

            const confirmed = confirm(`Вы уверены, что хотите удалить пользователя: ${selectedAdmin.value.first_name} ${selectedAdmin.value.second_name}?`);
            if (!confirmed) return;

            try {
                const res = await fetch(`/api/admins/${selectedAdmin.value.id}`, {
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
            loadRoles();
            loadStaff();
            document.addEventListener('click', this.handleClickOutside);
        await fetchAdminData();
    });


    return {
      search,
      filters,
      staff,
      roles,
      form,
      emailError,
      phoneError,
      loadStaff,
      getRoleName,
      openEdit,
      saveStaff,
      onRoleChange,
      userRole,
      onRowClick,
      openLoginChange,
      saveLogin,
      selectedAdmin,
      selectedAdminId,
      isSelected,
      onModalHidden,
      deleteAdmin,
      formatPhone,
      firstNameError,
      secondNameError,
      isPopoverVisible,
      fullName,
        filteredStaff
    };
  }
}).mount('#app');
