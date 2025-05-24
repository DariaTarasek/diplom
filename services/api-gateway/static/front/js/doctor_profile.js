const { createApp } = Vue;

createApp({
    data() {
        return {
            doctor: {
                second_name: '',
                first_name: '',
                surname: '',
                birthDate: '',
                gender: '',
                education: '',
                experience: '',
                phone: '',
                email: '',
            },
            showPasswordModal: false,
            newPassword: '',
            confirmPassword: '',
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
    }
  },
    computed: {
        fullName() {
            return [
                this.doctor.first_name,
                this.doctor.second_name
            ]
            .filter(Boolean)
            .join(' ');
        }
    },
    mounted() {
        this.fetchDoctorData();
        document.addEventListener('click', this.handleClickOutside);

    // добавка для тултипов:
    this.$nextTick(() => {
        const tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'));
        tooltipTriggerList.map(tooltipTriggerEl => new bootstrap.Tooltip(tooltipTriggerEl));
    });
    },
    beforeUnmount() {
        document.removeEventListener('click', this.handleClickOutside);
    },
    methods: {
        fetchDoctorData() {
            fetch('http://192.168.1.207:8080/api/doctor-profile')
                .then(response => response.json())
                .then(data => {
                    this.doctor = data;
                })
                .catch(error => {
                    console.error('Ошибка загрузки профиля:', error);
                });
        },

        changePassword() {
            if (this.newPasswordError || this.confirmPasswordError) {
            return;
        }

            fetch('http://192.168.1.207:8080/api/change-password', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ password: this.newPassword })
        }).then(response => {
            if (response.ok) {
            this.showPasswordModal = false;
            this.newPassword = '';
            this.confirmPassword = '';
            } else {
            this.newPasswordError = 'Ошибка при смене пароля.';
            }
        }).catch(() => {
            this.newPasswordError = 'Ошибка сети. Повторите попытку позже.';
        });
        },
        togglePopover() {
            this.isPopoverVisible = !this.isPopoverVisible;
        },
        handleClickOutside(event) {
            const popover = document.getElementById('doctor-profile');
            if (popover && !popover.contains(event.target)) {
                this.isPopoverVisible = false;
            }
        }
    }
}).mount('#app');
