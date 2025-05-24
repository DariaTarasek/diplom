const { createApp } = Vue;

createApp({
  data() {
    return {
      activeTab: 'schedule',
      second_name: '',
      first_name: '',
      tabs: [
        { id: 'schedule', label: 'Расписание приёмов' },
        { id: 'pending', label: 'Неподтверждённые записи' }
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
      appointmentSchedule: {},
      maxSlots: 0,
      selectedSlot: null,
      step: 1,

      patient: {
        id: '',
        second_name: '',
        first_name: '',
        surname: '',
        birth_date: '',
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

      currentWeekStartIndex: 0,  // индекс начала текущей видимой недели

    };
  },

  computed: {
    fullName() {
      return [this.first_name, this.second_name].filter(Boolean).join(' ');
    },
    visibleWeekDays() {
        return this.schedule.days.slice(this.currentWeekStartIndex, this.currentWeekStartIndex + 7);
    }
  },

  methods: {
    async fetchData() {
      try {
        const res = await fetch('http://192.168.1.207:8080/api/admin-data');
        const data = await res.json();
        this.appointments = data.appointments || {};
        this.pending = data.pending || [];
        this.first_name = data.first_name || '';
        this.second_name = data.second_name || '';

        const scheduleRes = await fetch('http://192.168.1.207:8080/api/schedule-admin');
        const scheduleData = await scheduleRes.json();
        this.schedule = {
          days: scheduleData.days || [],
          timeSlots: scheduleData.timeSlots || []
        };
      } catch (err) {
        console.error('Ошибка при загрузке данных:', err);
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
      const res = await fetch('http://192.168.1.207:8080/api/specialties');
      this.specialties = await res.json();
    },

   async fetchDoctors(specialtyId) {
  const res = await fetch(`http://192.168.1.207:8080/api/doctors?specialty=${specialtyId}`);
  const rawDoctors = await res.json();
  this.doctors = rawDoctors.map(doc => ({
    ...doc,
    fullName: `${doc.second_name} ${doc.first_name} ${doc.surname}`.trim()
  }));
},

async fetchDoctorSchedule(doctorId) {
  const res = await fetch(`http://192.168.1.207:8080/api/schedule?doctor_id=${doctorId}`);
  this.appointmentSchedule = await res.json();
  this.maxSlots = Math.max(...Object.values(this.appointmentSchedule).map(day => day.length));
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
        slot: this.selectedSlot,
        patient: { ...this.patient, phone: this.patient.phone.replace(/\D/g, '') }
      };

      const res = await fetch('http://192.168.1.207:8080/api/appointments', {
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

    resetModalData() {
        this.selectedSpecialization = null;
        this.selectedDoctorId = null;
        this.selectedDoctor = null;
        this.appointmentSchedule = {};
        this.selectedSlot = null;
        this.step = 1;
        this.patient = {
            second_name: '',
            first_name: '',
            surname: '',
            birth_date: '',
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

    showTransferModal() {
    if (!this.selectedAppt) return;

    this.appointmentSchedule = {};
    this.selectedTransferSlot = { date: null, time: null };


    this.fetchDoctorSchedule(this.selectedAppt.doctor.id); // <-- нужен doctor_id в `appt`

    const manageModal = bootstrap.Modal.getInstance(document.getElementById('manageAppointmentModal'));
    if (manageModal) manageModal.hide();

    const transferModal = new bootstrap.Modal(document.getElementById('transferAppointmentModal'));
    transferModal.show();
    },

    async rescheduleAppointment(newDate, newTime) {
    if (!this.selectedAppt) return;

    const payload = {
        doctor_id: this.selectedAppt.doctor.id,
        patient_id: this.selectedAppt.patient.id,
        old_slot: { day: this.selectedAppt.day, time: this.selectedAppt.time },
        new_slot: { day: newDate, time: newTime }
    };

    const res = await fetch('http://192.168.1.207:8080/api/reschedule', {
        method: 'POST',
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

    const payload = {
        doctor_id: this.selectedAppt.doctor.id,
        patient_id: this.selectedAppt.patient.id,
        day: this.selectedAppt.day,
        time: this.selectedAppt.time
    };

    const res = await fetch('http://192.168.1.207:8080/api/cancel-appointment', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
    });

    if (res.ok) {
        alert('Запись отменена');
        location.reload();
    } else {
        alert('Ошибка отмены записи');
    }
    },
     },

  watch: {
    'patient.first_name': 'validateFirstName',
    'patient.second_name': 'validateSecondName'
  },

  mounted() {
    this.fetchData();
    document.addEventListener('click', this.handleClickOutside);
  }
}).mount('#app');
