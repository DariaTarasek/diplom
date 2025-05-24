const { createApp } = Vue;

createApp({
  data() {
    return {
      activeTab: 'today',
      second_name: '',
      first_name: '',
      tabs: [
        { id: 'today', label: 'Записи на сегодня' },
        { id: 'upcoming', label: 'Расписание' },
        { id: 'completed', label: 'Завершенные приемы' }
      ],
      data: {
        today: [],
        upcoming: [],
        completed: []
      },
      scheduleDates: [],
      scheduleTimes: [],     
      isPopoverVisible: false,
    };
  },
  computed: {
    todaySorted() {
      return [...this.data.today].sort((a, b) => a.time.localeCompare(b.time));
    },
    fullName() {
      return [this.first_name, this.second_name].filter(Boolean).join(' ');
    }
  },
  methods: {
    fetchDoctorData() {
      fetch('http://192.168.1.207:8080/api/doctor-data')
        .then(res => {
          if (!res.ok) throw new Error('Ошибка загрузки');
          return res.json();
        })
        .then(json => {
          this.data.today = json.today || [];
          this.data.upcoming = json.upcoming || [];
          this.data.completed = json.completed || [];
          this.first_name = json.first_name;
          this.second_name = json.second_name;
        })
        .catch(err => {
          console.error('Ошибка при получении данных врача:', err);
        });
    },
    fetchSchedule() {
      fetch('http://192.168.1.207:8080/api/schedule-admin')
        .then(res => {
          if (!res.ok) throw new Error('Ошибка загрузки расписания');
          return res.json();
        })
        .then(json => {
          this.scheduleDates = json.days.map(d => d.date);
          this.scheduleTimes = json.timeSlots || [];
        })
        .catch(err => {
          console.error('Ошибка при получении расписания:', err);
        });
    },
    getPatientByDateTime(date, time) {
      return this.data.upcoming.find(x => x.date === date && x.time === time);
    },
    isToday(dateStr) {
      const today = new Date().toISOString().split('T')[0];
      return dateStr === today;
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
    }
  },
  mounted() {
    this.fetchDoctorData();
    this.fetchSchedule(); 
    document.addEventListener('click', this.handleClickOutside);
  }
}).mount('#app');
