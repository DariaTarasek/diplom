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
      data: {
        upcoming: [],
        history: [],
        tests: []
      }
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
    this.fetchPatientData();
    document.addEventListener('click', this.handleClickOutside);
  },
  beforeUnmount() {
    document.removeEventListener('click', this.handleClickOutside);
  },
  methods: {
    fetchPatientData() {
      fetch('http://192.168.1.207:8080/api/patient-data', {method: "GET",
      credentials: "include"})
        .then(res => {
          if (!res.ok) throw new Error('Ошибка загрузки');
          return res.json();
        })
        .then(json => {
          this.patient.secondName = json.second_name;
          this.patient.firstName = json.first_name;
          this.patient.surname = json.surname; 
          this.data.upcoming = json.upcoming || [];
          this.data.history = json.history || [];
          this.data.tests = json.tests || [];
        })
        .catch(err => {
          console.error('Ошибка при получении данных:', err);
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
