const { createApp } = Vue;

createApp({
  data() {
    return {
      patient: {
        secondName: '',
        firstName: '',
        surname: ''
      },
      isPopoverVisible: false,
      activeTab: 'upcoming',
      tabs: [
        { id: 'upcoming', label: 'Предстоящие записи' },
        { id: 'history', label: 'История посещений' },
        { id: 'tests', label: 'Анализы / Исследования' }
      ],
      upcoming: [],
      history: [],
      tests: [],
      selectedAppt: null,
      appointmentSchedule: {},
      selectedTransferSlot: null,
      paginatedSchedule: [],
      maxSlots: 0,
      currentPage: 0,
      daysPerPage: 7,
    };
  },

  computed: {
    fullName() {
      return [this.patient.firstName, this.patient.secondName]
        .filter(Boolean)
        .join(' ');
    }
  },
  mounted() {
    this.fetchPatientInfo();
    this.fetchUpcoming();
    this.fetchHistory();
    this.fetchTests();
    document.addEventListener('click', this.handleClickOutside);
  },

  beforeUnmount() {
    document.removeEventListener('click', this.handleClickOutside);
  },
  methods: {
    fetchPatientInfo() {
      fetch('/api/patient/me', { method: 'GET', credentials: 'include' })
          .then(res => res.json())
          .then(json => {
            this.patient.secondName = json.secondName;
            this.patient.firstName = json.firstName;
            this.patient.surname = json.surname;
          })
          .catch(err => console.error('Ошибка при получении информации о пациенте:', err));
    },
    fetchUpcoming() {
      fetch('/api/patient/upcoming', { method: 'GET', credentials: 'include' })
          .then(res => res.json())
          .then(json => {
            this.upcoming = json || [];
          })
          .catch(err => console.error('Ошибка при получении предстоящих записей:', err));
    },
    fetchHistory() {
      fetch('/api/patient/history', { method: 'GET', credentials: 'include' })
          .then(res => res.json())
          .then(json => {
            this.history = json || [];
          })
          .catch(err => console.error('Ошибка при получении истории посещений:', err));
    },
    fetchTests() {
      fetch('/api/patient/tests', { method: 'GET', credentials: 'include' })
          .then(res => res.json())
          .then(json => {
            this.tests = json || [];
          })
          .catch(err => console.error('Ошибка при получении анализов:', err));
    },
    togglePopover() {
      this.isPopoverVisible = !this.isPopoverVisible;
    },
    handleClickOutside(event) {
      const popover = document.getElementById('patient-profile');
      if (popover && !popover.contains(event.target)) {
        this.isPopoverVisible = false;
      }
    },
    openManageModal(item) {
      this.selectedAppt = item;
      const modal = new bootstrap.Modal(document.getElementById('manageAppointmentModal'));
      modal.show();
    },

    showTransferModal() {
      // Пример запроса слотов, замени URL и формат под себя
      fetch(`/api/appointment-doctor-schedule/${this.selectedAppt.doctorId}`, {
        method: 'GET',
        credentials: 'include'
      })
          .then(res => res.json())
          .then(json => {
            this.appointmentSchedule = json;
            this.prepareScheduleTable(); // обновляем таблицу
            const modal = new bootstrap.Modal(document.getElementById('transferAppointmentModal'));
            modal.show();
          })
          .catch(err => console.error('Ошибка при получении слотов:', err));
    },

    selectTransferSlot(date, time) {
      this.selectedTransferSlot = { date, time };
    },

    confirmTransfer() {
      const payload = {
        id: this.selectedAppt.id,
        date: this.selectedTransferSlot.date,
        time: this.selectedTransferSlot.time
      };

      fetch('/api/appointments/transfer', {
        method: 'PUT',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      })
          .then(res => res.json())
          .then(() => {
            bootstrap.Modal.getInstance(document.getElementById('transferAppointmentModal')).hide();
            bootstrap.Modal.getInstance(document.getElementById('manageAppointmentModal')).hide();
            this.fetchUpcoming(); // обновим список
          })
          .catch(err => console.error('Ошибка при переносе записи:', err));
    },

    confirmCancelAppointment() {
      if (!confirm('Вы уверены, что хотите отменить запись?')) return;

      fetch(`/api/appointments/cancel/${this.selectedAppt.id}`, {
        method: 'GET',
        credentials: 'include'
      })
          .then(() => {
            bootstrap.Modal.getInstance(document.getElementById('manageAppointmentModal')).hide();
            this.fetchUpcoming();
          })
          .catch(err => console.error('Ошибка при отмене записи:', err));
    },
    prepareScheduleTable() {
      if (!Array.isArray(this.appointmentSchedule)) {
        this.paginatedSchedule = [];
        this.maxSlots = 0;
        return;
      }

      const normalized = this.appointmentSchedule.map(day => ({
        label: day.label,
        slots: Array.isArray(day.slots) ? day.slots : [],
      }));

      // пагинация
      const start = this.currentPage * this.daysPerPage;
      const end = start + this.daysPerPage;
      this.paginatedSchedule = normalized.slice(start, end);

      this.maxSlots = Math.max(0, ...this.paginatedSchedule.map(day => day.slots.length));
    },
    nextPage() {
      const maxPage = Math.ceil(this.appointmentSchedule.length / this.daysPerPage) - 1;
      if (this.currentPage < maxPage) {
        this.currentPage++;
        this.prepareScheduleTable();
      }
    },
    prevPage() {
      if (this.currentPage > 0) {
        this.currentPage--;
        this.prepareScheduleTable();
      }
    },
    isSelectedSlot(date, time) {
      return this.selectedTransferSlot &&
          this.selectedTransferSlot.date === date &&
          this.selectedTransferSlot.time === time;
    }
  }
}).mount('#app');
