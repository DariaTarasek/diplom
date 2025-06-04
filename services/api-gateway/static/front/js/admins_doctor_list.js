const { createApp, ref, reactive, watch, onMounted, computed } = Vue;

createApp({
    setup() {
        const search = ref('');
        const filters = reactive({ specialty: '' });
        const staff = ref([]);
        const allStaff = ref([]); // <-- Список всех сотрудников
        const specialties = ref([]);

        const form = reactive({
            user_id: '',
            secondName: '',
            firstName: '',
            surname: '',
            email: '',
            specialty: [],
            gender: '',
            phone: '',
            experience: '',
            education: ''
        });

        const admin = ref({
            secondName: '',
            firstName: '',
            role: ''
        });

        const fullName = computed(() =>
            [admin.value.firstName, admin.value.secondName].filter(Boolean).join(' ')
        );

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
            specialties.value = await (await fetch('/api/specialties')).json();
        }

        async function loadStaff() {
            const res = await fetch('/api/staff-doctors');
            allStaff.value = await res.json();
            applyFilters();
        }

        function applyFilters() {
            const searchTerm = search.value.trim().toLowerCase();
            const selectedSpecialty = filters.specialty;

            staff.value = allStaff.value.filter(s => {
                const fullName = `${s.secondName} ${s.firstName} ${s.surname}`.toLowerCase();
                const email = (s.email || '').toLowerCase();

                const matchesSearch =
                    !searchTerm ||
                    fullName.includes(searchTerm) ||
                    email.includes(searchTerm);

                const matchesSpecialty =
                    !selectedSpecialty ||
                    s.specialty.includes(Number(selectedSpecialty));

                return matchesSearch && matchesSpecialty;
            });
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

        function onRowClick(s) {
            selectedDoctor.value = s;
            selectedDoctorId.value = s.user_id;

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
            form.gender = s.gender === "м" ? "м" : "ж";
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
                    if (phoneInput.maskRef?.destroy) phoneInput.maskRef.destroy();

                    let digits = (form.phone || '').replace(/\D/g, '').slice(1);

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
                    });

                    mask.unmaskedValue = digits;
                    form.phone = mask.value;
                }
            }, 300);
        }

        function saveStaff() {
            const phoneValid = !phoneError.value && form.phone;

            if (!phoneValid || specError.value) return;

            secondNameError.value = '';
            firstNameError.value = '';

            if (!form.secondName?.trim()) secondNameError.value = 'Фамилия не может быть пустой';
            if (!form.firstName?.trim()) firstNameError.value = 'Имя не может быть пустым';
            if (secondNameError.value || firstNameError.value) return;

            const payload = {
                user_id: form.user_id,
                firstName: form.firstName,
                secondName: form.secondName,
                surname: form.surname,
                phone: form.phone.replace(/\D/g, '').trim(),
                email: form.email,
                specialty: Array.isArray(form.specialty) ? form.specialty : [form.specialty],
                experience: form.experience,
                education: form.education,
                gender: form.gender
            };

            fetch('/api/save-doctor', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload)
            }).then(response => {
                if (response.ok) {
                    loadStaff();
                    modalEdit?.hide();
                    selectedDoctorId.value = null;
                    selectedDoctor.value = null;
                } else {
                    alert('Ошибка при сохранении (код ' + response.status + ')');
                }
            });
        }

        function isSelected(s) {
            return selectedDoctorId.value === s.user_id;
        }

        function openLoginChange() {
            isSwitchingModals.value = true;
            modalChoice.hide();

            if (!modalLogin) {
                modalLogin = new bootstrap.Modal(document.getElementById('changeLoginModal'));
            }

            if (selectedDoctor.value) {
                form.user_id = selectedDoctor.value.user_id;
                form.email = selectedDoctor.value.email;
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

            await fetch(`/api/doctors-login/${form.user_id}`, {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email: form.email })
            });

            modalLogin.hide();
            selectedDoctorId.value = null;
            selectedDoctor.value = null;
            await loadStaff();
        }

        async function fetchAdminData() {
            try {
                const res = await fetch('/api/admin/me');
                const data = await res.json();

                admin.value.firstName = data.firstName || '';
                admin.value.secondName = data.secondName || '';
            } catch (err) {
                console.error('Ошибка при загрузке данных:', err);
            }
        }

        async function deleteDoctor() {
            if (!selectedDoctor.value) return;

            const confirmed = confirm(`Вы уверены, что хотите удалить пользователя: ${selectedDoctor.value.firstName} ${selectedDoctor.value.secondName}?`);
            if (!confirmed) return;

            try {
                const res = await fetch(`/api/doctors/${selectedDoctor.value.user_id}`, {
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

        // === WATCH ===

        watch(() => form.email, (val) => {
            const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
            emailError.value = val && !emailRegex.test(val) ? 'Неверный формат email' : '';
        });

        watch(() => form.specialty, (val) => {
            specError.value = val.length === 0 ? 'Выберите хотя бы одну специализацию' : '';
        });

        watch(() => form.firstName, (val) => {
            firstNameError.value = val.trim().length === 0 ? 'Имя не может быть пустым' : '';
        });

        watch(() => form.secondName, (val) => {
            secondNameError.value = val.trim().length === 0 ? 'Фамилия не может быть пустой' : '';
        });

        watch([search, () => filters.specialty], applyFilters);

        onMounted(async () => {
            const editEl = document.getElementById('editModal');
            const loginEl = document.getElementById('changeLoginModal');
            const actionsEl = document.getElementById('actionsModal');

            if (editEl) editEl.addEventListener('hidden.bs.modal', onModalHidden);
            if (loginEl) loginEl.addEventListener('hidden.bs.modal', onModalHidden);
            if (actionsEl) actionsEl.addEventListener('hidden.bs.modal', onModalHidden);

            document.addEventListener('click', handleClickOutside);

            await loadSpecialties();
            await loadStaff();
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
            fullName,
            admin
        };
    }
}).mount('#app');
