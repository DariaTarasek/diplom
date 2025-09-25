const { createApp } = Vue;

createApp({
    data() {
        return {
            activeTab: 'schedule',
            secondName: '',
            firstName: '',
            tabs: [
                { id: 'schedule', label: 'Расписание приёмов' },
                { id: 'pending', label: 'Неподтверждённые записи' },
                { id: 'completed', label: 'Завершённые приёмы' }
            ],
            schedule: {
                days: [],
                timeSlots: []
            },
            appointments: {},
            pending: [],
            isPopoverVisible: false,

            // Модалка записи
            showModal: false,
            specialties: [],
            doctors: [],
            selectedSpecialization: null,
            selectedDoctorId: null,
            selectedDoctor: null,
            appointmentSchedule: [],
            maxSlots: 0,
            selectedSlot: null,
            step: 1,

            patient: {
                id: '',
                second_name: '',
                first_name: '',
                surname: '',
                birthDate: '',
                gender: '',
                phone: ''
            },
            doctor: {
                id: '',
                second_name: '',
                first_name: '',
                surname: '',
                specialty: ''
            },
            errors: {
                phone: '',
                first_name: '',
                second_name: ''
            },
            birthDateAttrs: {
                min: '',
                max: ''
            },

            selectedAppt: null,

            selectedTransferSlot: {
                date: null,
                time: null
            },

            currentWeekStartIndex: 0,
            completed: [],
            currentPage: 0
        };
    },

    computed: {
        fullName() {
            return [this.firstName, this.secondName].filter(Boolean).join(' ');
        },
        visibleWeekDays() {
            return this.schedule.days.slice(this.currentWeekStartIndex, this.currentWeekStartIndex + 7);
        },
        paginatedSchedule() {
            const start = this.currentPage * 7;
            const end = start + 7;
            return this.appointmentSchedule.slice(start, end);
        },
        totalPages() {
            return Math.ceil(this.appointmentSchedule.length / 7);
        }
    },

    methods: {
        async fetchData() {
            try {
                const scheduleRes = await fetch('/api/schedule-admin');
                const scheduleData = await scheduleRes.json();
                this.schedule = {
                    days: scheduleData.schedule?.days || [],
                    timeSlots: scheduleData.schedule?.timeSlots || []
                };
                this.appointments = scheduleData.appointments || {};

                const res = await fetch('/api/admin/me');
                const data = await res.json();
                this.firstName = data.firstName || '';
                this.secondName = data.secondName || '';
            } catch (err) {
                console.error('Ошибка при загрузке данных:', err);
            }
        },

        formatDateTime(dt) {
            const d = new Date(dt);
            return d.toLocaleString('ru-RU', {
                day: '2-digit',
                month: '2-digit',
                year: 'numeric',
                hour: '2-digit',
                minute: '2-digit'
            });
        },

        async fetchCompletedVisits() {
            try {
                const res = await fetch('/api/completed-visits');
                const data = await res.json();
                this.completed = data || [];
            } catch (err) {
                console.error('Ошибка при загрузке завершённых приёмов:', err);
            }
        },

        confirmEntry(index) {
            const entry = this.pending[index];
            this.pending.splice(index, 1);
        },

        togglePopover() {
            this.isPopoverVisible = !this.isPopoverVisible;
        },

        handleClickOutside(event) {
            const popover = document.getElementById('admin-profile');
            if (popover && !popover.contains(event.target)) {
                this.isPopoverVisible = false;
            }
        },

        validatePhone() {
            const phone = this.patient.phone.replace(/\D/g, '');
            this.errors.phone = phone.length === 11 && phone.startsWith('7') ? '' : 'Неверный формат телефона';
            return !this.errors.phone;
        },

        validateFirstName() {
            this.errors.first_name = this.patient.first_name.trim() ? '' : 'Имя не может быть пустым';
            return !this.errors.first_name;
        },

        validateSecondName() {
            this.errors.second_name = this.patient.second_name.trim() ? '' : 'Фамилия не может быть пустой';
            return !this.errors.second_name;
        },

        async fetchSpecialties() {
            const res = await fetch('/api/specialties');
            this.specialties = await res.json();
        },

        async fetchDoctors(specialtyId) {
            const res = await fetch(`/api/doctors/${specialtyId}`);
            const rawDoctors = await res.json();
            this.doctors = rawDoctors.map(doc => ({
                ...doc,
                fullName: `${doc.secondName} ${doc.firstName} ${doc.surname}`.trim()
            }));
        },

        async fetchDoctorSchedule(doctorId) {
            const res = await fetch(`/api/appointment-doctor-schedule/${doctorId}`);
            const data = await res.json();

            data.forEach(day => {
                if (!Array.isArray(day.slots)) {
                    day.slots = [];
                }
            });

            this.appointmentSchedule = data;
            this.maxSlots = Math.max(...data.map(d => d.slots?.length || 0));
        },

        selectSlot(slot) {
            this.selectedSlot = slot;
            this.step = 2;
            this.$nextTick(() => this.initPhoneMask());
        },

        back() {
            this.step = 1;
            this.selectedSlot = null;
        },

        async submitForm() {
            if (!this.selectedSlot || !this.selectedDoctorId) {
                alert('Выберите врача и время');
                return;
            }

            if (!this.validatePhone() || !this.validateFirstName() || !this.validateSecondName()) return;

            const payload = {
                doctor_id: this.selectedDoctorId,
                date: this.selectedSlot.date,
                time: this.selectedSlot.time,
                ...this.patient,
                phone: this.patient.phone.replace(/\D/g, '')
            };

            const res = await fetch('/api/appointments', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload)
            });

            if (res.ok) {
                alert('Запись успешно создана!');
                location.reload();
            } else {
                alert('Ошибка при записи. Попробуйте позже.');
            }
        },

        validateDateRange() {
            const today = new Date();
            const yyyy = today.getFullYear();
            const mm = String(today.getMonth() + 1).padStart(2, '0');
            const dd = String(today.getDate()).padStart(2, '0');
            this.birthDateAttrs.min = `${yyyy - 110}-${mm}-${dd}`;
            this.birthDateAttrs.max = `${yyyy - 18}-${mm}-${dd}`;
        },

        initPhoneMask() {
            const phoneInput = document.getElementById('phone');
            if (phoneInput && !phoneInput.dataset.masked) {
                const mask = IMask(phoneInput, {
                    mask: '+{7} (000) 000-00-00'
                });
                phoneInput.dataset.masked = "true";

                mask.on('accept', () => {
                    const digits = mask.value.replace(/\D/g, '');
                    const valid = digits.length === 11 && digits.startsWith('7');
                    this.errors.phone = digits && !valid ? 'Неверный формат телефона' : '';
                });

                mask.on('complete', () => {
                    this.patient.phone = mask.value;
                });
            }
        },

        openModal() {
            this.resetModalData();
            this.fetchSpecialties();
            this.validateDateRange();

            const modalEl = document.getElementById('appointmentModal');
            const modal = new bootstrap.Modal(modalEl);
            modal.show();
        },

        async confirmVisit(entry) {
            try {
                const res = await fetch(`/api/completed-visits/${entry.visit_id}`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        price: entry.price,
                        status: 'confirmed'
                    })
                });

                if (!res.ok) throw new Error('Ошибка при подтверждении');

                alert('Приём подтверждён');
                await this.fetchCompletedVisits(); // Обновить данные
            } catch (err) {
                console.error(err);
                alert('Не удалось подтвердить приём');
            }
        },

        resetModalData() {
            this.selectedSpecialization = null;
            this.selectedDoctorId = null;
            this.selectedDoctor = null;
            this.appointmentSchedule = [];
            this.selectedSlot = null;
            this.step = 1;
            this.patient = {
                second_name: '',
                first_name: '',
                surname: '',
                birthDate: '',
                gender: '',
                phone: ''
            };
            this.errors = {
                phone: '',
                first_name: '',
                second_name: ''
            };
        },

        closeModal() {
            this.showModal = false;
        },

        openAppointmentModal(day, time, appt) {
            this.selectedAppt = { ...appt, day, time };
            const modal = new bootstrap.Modal(document.getElementById('manageAppointmentModal'));
            modal.show();
        },

        async fetchUnconfirmedAppointments() {
            try {
                const res = await fetch('/api/unconfirmed-appointments');
                const rawData = await res.json();

                if (!Array.isArray(rawData)) {
                    this.pending = [];
                    return;
                }

                this.pending = rawData.map(entry => {
                    return {
                        id: entry.id,
                        date: entry.date,
                        time: entry.time,
                        name: `${entry.patient_second_name} ${entry.patient_first_name} ${entry.patient_surname}`,
                        birthDate: entry.patient_birth_date,
                        phone: `+${entry.phone_number}`, // Уже с "7" впереди
                        doctor: entry.doctor // Строка — ок
                    };
                });
            } catch (err) {
                console.error('Ошибка при загрузке неподтверждённых записей:', err);
                this.pending = [];
            }
        },

        async updateAppointmentStatus(entry) {
            const payload = {
                id: entry.id,
                date: entry.date,
                time: entry.time,
                status: 'confirmed',
                updated_at: new Date().toISOString()
            };
            try {
                const res = await fetch(`/api/unconfirmed-appointments/${entry.id}`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(payload)
                });

                if (!res.ok) throw new Error('Ошибка при подтверждении записи');
                alert('Запись подтверждена!');
                await this.fetchUnconfirmedAppointments();
            } catch (err) {
                console.error('Ошибка:', err);
                alert('Не удалось обновить статус записи');
            }
        },


        showTransferModal() {
            if (!this.selectedAppt) return;

            this.appointmentSchedule = [];
            this.selectedTransferSlot = { date: null, time: null };
            this.fetchDoctorSchedule(this.selectedAppt.doctor.id);

            const manageModal = bootstrap.Modal.getInstance(document.getElementById('manageAppointmentModal'));
            if (manageModal) manageModal.hide();

            const transferModal = new bootstrap.Modal(document.getElementById('transferAppointmentModal'));
            transferModal.show();
        },

        async rescheduleAppointment(newDate, newTime) {
            if (!this.selectedAppt) return;

            const payload = {
                id: this.selectedAppt.id,
                doctor_id: this.selectedAppt.doctor.id,
                patient_id: this.selectedAppt.patient.id,
                date: newDate,
                time: newTime
            };

            const res = await fetch('/api/appointments/transfer', {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload)
            });

            if (res.ok) {
                alert('Запись успешно перенесена');
                location.reload();
            } else {
                alert('Ошибка переноса');
            }
        },

        confirmTransfer() {
            if (!this.selectedTransferSlot.date || !this.selectedTransferSlot.time) return;
            this.rescheduleAppointment(this.selectedTransferSlot.date, this.selectedTransferSlot.time);
        },

        selectTransferSlot(date, time) {
            this.selectedTransferSlot = { date, time };
        },

        prevWeek() {
            if (this.currentWeekStartIndex >= 7) {
                this.currentWeekStartIndex -= 7;
            }
        },
        nextWeek() {
            if (this.currentWeekStartIndex + 7 < this.schedule.days.length) {
                this.currentWeekStartIndex += 7;
            }
        },

        async confirmCancelAppointment() {
            if (!confirm('Вы уверены, что хотите отменить запись?')) return;

            const res = await fetch(`/api/appointments/cancel/${this.selectedAppt.id}`, {
                method: 'GET',
                headers: { 'Content-Type': 'application/json' }
            });

            if (res.ok) {
                alert('Запись отменена');
                location.reload();
            } else {
                alert('Ошибка отмены записи');
            }
        },

        prevPage() {
            if (this.currentPage > 0) {
                this.currentPage--;
            }
        },

        nextPage() {
            if ((this.currentPage + 1) * 7 < this.appointmentSchedule.length) {
                this.currentPage++;
            }
        },

        selectDateSlot(daySchedule, slot) {
            this.selectedSlot = {
                date: daySchedule.label,
                time: slot
            };
            this.step = 2;
            this.$nextTick(() => this.initPhoneMask());
        }
    },

    watch: {
        'patient.first_name': 'validateFirstName',
        'patient.second_name': 'validateSecondName'
    },

    mounted() {
        this.fetchData();
        this.fetchCompletedVisits();
        this.fetchUnconfirmedAppointments();
        document.addEventListener('click', this.handleClickOutside);
    }
}).mount('#app');
