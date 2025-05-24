const { createApp } = Vue;

createApp({
    data() {
        return {
            patient: {
                second_name: '',
                first_name: '',
                surname: '',
                birthDate: '',
                gender: '',
                phone: '',
                email: '',
            },
            showEmailModal: false,
            showPasswordModal: false,
            newEmail: '',
            newPassword: '',
            confirmPassword: '',
            newEmailError: '',
            newPasswordError: '',
            confirmPasswordError: '',
            isPopoverVisible: false, // добавили для поповера
        };
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
    computed: {
        fullName() {
            return [
                this.patient.first_name,
                this.patient.second_name,
            ]
            .filter(Boolean)
            .join(' ');
        }
    },
    mounted() {
        this.fetchPatientData();
        document.addEventListener('click', this.handleClickOutside);

    // Вот добавка для тултипов:
    this.$nextTick(() => {
        const tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'));
        tooltipTriggerList.map(tooltipTriggerEl => new bootstrap.Tooltip(tooltipTriggerEl));
    });
    },
    beforeUnmount() {
        document.removeEventListener('click', this.handleClickOutside);
    },
    methods: {
        fetchPatientData() {
            fetch('http://192.168.1.207:8080/api/patient-profile')
                .then(response => response.json())
                .then(data => {
                    this.patient = data;
                })
                .catch(error => {
                    console.error('Ошибка загрузки профиля:', error);
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
                    this.patient.email = this.newEmail;
                    this.showEmailModal = false;
                    this.newEmail = '';
                } else {
                    alert('Ошибка при смене email');
                }
            });
        },
        changePassword() {
            if (this.newPassword !== this.confirmPassword) {
                return;
            }
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
        togglePopover() {
            this.isPopoverVisible = !this.isPopoverVisible;
        },
        handleClickOutside(event) {
            const popover = document.getElementById('patient-profile');
            if (popover && !popover.contains(event.target)) {
                this.isPopoverVisible = false;
            }
        }
    }
}).mount('#app');
