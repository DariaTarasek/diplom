const { createApp } = Vue;

createApp({
    data() {
        return {
            admin: {
                second_name: '',
                first_name: '',
                surname: '',
                birthDate: '',
                gender: '',
                phone: '',
                email: '',
            },
            role: '',

            editingInfo: false,

            showEmailModal: false,
            showPasswordModal: false,

            newEmail: '',
            newPassword: '',
            confirmPassword: '',

            newEmailError: '',
            newPasswordError: '',
            confirmPasswordError: '',

            isPopoverVisible: false,
        };
    },
    computed: {
        fullName() {
            return [
                this.admin.first_name,
                this.admin.second_name
            ]
            .filter(Boolean)
            .join(' ');
        }
    },
    mounted() {
        this.fetchAdminData();
        document.addEventListener('click', this.handleClickOutside);

        this.$nextTick(() => {
            const tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'));
            tooltipTriggerList.map(el => new bootstrap.Tooltip(el));
        });
    },
    beforeUnmount() {
        document.removeEventListener('click', this.handleClickOutside);
    },
    watch: {
        newPassword(value) {
            const strongRegex = /^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,}$/;
            if (!value) {
                this.newPasswordError = '';
            } else if (!strongRegex.test(value)) {
                this.newPasswordError = 'Пароль должен содержать минимум 8 символов, включая латинские буквы и цифры.';
            } else {
                this.newPasswordError = '';
            }

            if (this.confirmPassword && value !== this.confirmPassword) {
                this.confirmPasswordError = 'Пароли не совпадают.';
            } else {
                this.confirmPasswordError = '';
            }
        },
        confirmPassword(value) {
            if (!value) {
                this.confirmPasswordError = '';
            } else if (value !== this.newPassword) {
                this.confirmPasswordError = 'Пароли не совпадают.';
            } else {
                this.confirmPasswordError = '';
            }
        },
        newEmail(value) {
            const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]{2,}$/;
            if (!value) {
                this.newEmailError = '';
            } else if (!emailRegex.test(value)) {
                this.newEmailError = 'Некорректный формат email.';
            } else {
                this.newEmailError = '';
            }
        }
    },
    methods: {
        fetchAdminData() {
            fetch('http://192.168.1.207:8080/api/admin-profile')
                .then(response => response.json())
                .then(data => {
                    this.admin = data.admin;
                    this.role = data.role;
                })
                .catch(error => {
                    console.error('Ошибка загрузки профиля администратора:', error);
                });
        },
        changeEmail() {
            fetch('http://192.168.1.207:8080/api/change-email', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email: this.newEmail })
            })
            .then(response => {
                if (response.ok) {
                    this.admin.email = this.newEmail;
                    this.showEmailModal = false;
                    this.newEmail = '';
                } else {
                    alert('Ошибка при смене email');
                }
            });
        },
        changePassword() {
            if (this.newPassword !== this.confirmPassword) return;

            fetch('http://192.168.1.207:8080/api/change-password', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ password: this.newPassword })
            })
            .then(response => {
                if (response.ok) {
                    this.showPasswordModal = false;
                    this.newPassword = '';
                    this.confirmPassword = '';
                } else {
                    alert('Ошибка при смене пароля');
                }
            });
        },
        saveAdminData() {
            fetch('http://192.168.1.207:8080/api/update-admin-profile', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(this.admin)
            })
            .then(response => {
                if (!response.ok) {
                    alert('Ошибка при сохранении профиля');
                } else {
                    this.editingInfo = false;
                }
            })
            .catch(error => {
                console.error('Ошибка при сохранении профиля:', error);
            });
        },
        togglePopover() {
            this.isPopoverVisible = !this.isPopoverVisible;
        },
        handleClickOutside(event) {
            const popover = document.getElementById('admin-profile');
            if (popover && !popover.contains(event.target)) {
                this.isPopoverVisible = false;
            }
        }
    }
}).mount('#app');
