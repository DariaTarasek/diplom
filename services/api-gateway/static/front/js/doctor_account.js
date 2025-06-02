const { createApp } = Vue;

createApp({
  data() {
    return {
      activeTab: 'today',
      first_name: '',
      second_name: '',
      tabs: [
        { id: 'today', label: 'Записи на сегодня' },
        { id: 'upcoming', label: 'Расписание' },
      ],
      data: {
        today: null,
        upcoming: null,
      },
      scheduleDates: [],
      scheduleTimes: [],
      table: {},
      currentPageX: 1,
      pageSizeX: 7,
      isPopoverVisible: false,
    };
  },
  computed: {
    todaySorted() {
      return (this.data.today || []).slice().sort((a, b) => a.time.localeCompare(b.time));
    },
    fullName() {
      return [this.first_name, this.second_name].filter(Boolean).join(' ');
    },
    paginatedDates() {
      const start = (this.currentPageX - 1) * this.pageSizeX;
      return this.scheduleDates.slice(start, start + this.pageSizeX);
    },
    totalPagesX() {
      return Math.ceil(this.scheduleDates.length / this.pageSizeX);
    }
  },
  watch: {
    activeTab(newTab) {
      if (newTab === 'today' && this.data.today === null) this.fetchTodayAppointments();
      if (newTab === 'upcoming' && this.data.upcoming === null) this.fetchUpcomingSchedule();
    }
  },
  methods: {
    fetchTodayAppointments() {
      fetch('/api/appointments-today')
          .then(res => res.json())
          .then(json => { this.data.today = json || []; })
          .catch(err => console.error('Ошибка при загрузке записей на сегодня:', err));
    },
    fetchUpcomingSchedule() {
      fetch('/api/schedule-with-appointments')
          .then(res => res.json())
          .then(data => {
            this.scheduleDates = data.dates;
            this.scheduleTimes = data.times;
            this.table = data.table;
          })
          .catch(err => console.error('Ошибка при загрузке расписания:', err));
    },

    getPatientByDateTime(date, time) {
      return (this.data.upcoming || []).find(x => x.date === date && x.time === time);
    },
    formatDateToIso(dateStr) {
      const [day, month, year] = dateStr.split('.');
      return `${year}-${month}-${day}`;
    },

    isToday(dateStr) {
      const today = new Date().toISOString().split('T')[0];
      return this.formatDateToIso(dateStr) === today;
    },

    startConsultation(appointmentId) {
      if (appointmentId) {
        window.location.href = `doctors_consultation.html?appointment_id=${appointmentId}`;
      } else {
        alert('Идентификатор записи не найден.');
      }
    },
    togglePopover() {
      this.isPopoverVisible = !this.isPopoverVisible;
    },
    handleClickOutside(event) {
      const popover = document.getElementById('doctor-profile');
      if (popover && !popover.contains(event.target)) {
        this.isPopoverVisible = false;
      }
    },
    nextPageX() {
      if (this.currentPageX < this.totalPagesX) this.currentPageX++;
    },
    prevPageX() {
      if (this.currentPageX > 1) this.currentPageX--;
    }
  },
  mounted() {
    // подгружаем только первую вкладку
    this.fetchTodayAppointments();
    document.addEventListener('click', this.handleClickOutside);
  }
}).mount('#app');
