const { createApp } = Vue;

createApp({
  data() {
    return {
        dayLabels: ['Воскресенье', 'Понедельник', 'Вторник', 'Среда', 'Четверг', 'Пятница', 'Суббота'],
      clinicSchedule: [],
      clinicSlotMinutes: null,
      doctors: [],
      selectedDoctor: null,
        selectedDoctorId: null,
        doctorSchedule: [],
      doctorSlotMinutes: null,
      overrideForm: {
            doctor_id: null,
            date: '',
            type: 'off',
            start_time: '',
            end_time: '',
        },
    isPopoverVisible: false,
    admin: {
        first_name: '',
        second_name: '',
        role: '',
        },
    invalidDays: [],
    invalidClinicDays: [],
    clinicSlotError: false,
    doctorSlotError: false,
    //overrideSlotError: false,
    overrideErrors: {
        date: false,
        start_time: false,
        end_time: false,
        },
     overrideClinicForm: {
        date: '',
        type: 'off',
        start_time: '',
        end_time: '',
         },
    overrideClinicErrors: {
        date: false,
        start_time: false,
        end_time: false,
        },
    };
},
  computed: {
    fullName() {
      return [this.admin.first_name, this.admin.second_name].filter(Boolean).join(' ');
        },
    isSuperAdmin() {
         return this.admin.role === "superadmin";
  },
      clinicScheduleSorted() {
              return [...this.clinicSchedule].sort((a, b) => {
                  // Сначала понедельник (1), затем вторник (2), ..., воскресенье (0) в конце
                  const order = [1, 2, 3, 4, 5, 6, 0];
                  return order.indexOf(a.day) - order.indexOf(b.day);
              });
      },
      doctorScheduleSorted() {
          return [...this.doctorSchedule].sort((a, b) => {
              const order = [1, 2, 3, 4, 5, 6, 0];
              return order.indexOf(a.day) - order.indexOf(b.day);
          });
      },
    },
  mounted() {
    this.loadClinicSchedule();
    this.loadDoctors();
    this.fetchData(); 
    document.addEventListener('click', this.handleClickOutside);
  },
        watch: {
            'overrideForm.date'(newDate) {
                if (!/^\d{4}-\d{2}-\d{2}$/.test(newDate)) return;
                if (this.overrideForm.doctor_id && newDate.length === 10) {
                    this.loadDoctorOverride(this.overrideForm.doctor_id, newDate);
                }
            },
            'overrideClinicForm.date'(newDate) {
                if (!/^\d{4}-\d{2}-\d{2}$/.test(newDate)) return;
                if (newDate.length === 10) this.loadClinicOverride(newDate);
            },
            'overrideForm.doctor_id'(newDoctorId) {
                if (newDoctorId && this.overrideForm.date && this.overrideForm.date.length === 10) {
                    this.loadDoctorOverride(newDoctorId, this.overrideForm.date);
                }
            },
            selectedDoctor(newVal) {
                if (newVal) {
                    this.fetchDoctorSchedule();
                }
            }
        },
  methods: {
    togglePopover() {
        this.isPopoverVisible = !this.isPopoverVisible;
    },
      async fetchData() {
      try {
        const res = await fetch('/api/admin-data');
        const data = await res.json();

        this.admin.first_name = data.first_name || '';
        this.admin.second_name = data.second_name || '';
        this.admin.role = data.role || '';
        }
        catch (err) {
      }
    },

      async loadClinicOverride(date) {
          this.overrideClinicErrors.date = false;
          const datePattern = /^\d{4}-\d{2}-\d{2}$/;


          try {
              const res = await fetch(`/api/clinic-overrides/${date}`);
              if (!res.ok) {
                  if (res.status === 404) {
                      this.overrideClinicForm.type = 'off';
                      this.overrideClinicForm.start_time = '';
                      this.overrideClinicForm.end_time = '';
                      return;
                  }
                  console.error('Ошибка при загрузке переопределения клиники:', res.statusText);
                  return;
              }

              const data = await res.json();
              this.overrideClinicForm.type = data.is_day_off ? 'off' : 'work';
              this.overrideClinicForm.start_time = data.start_time || '';
              this.overrideClinicForm.end_time = data.end_time || '';
          } catch (err) {
              console.error('Ошибка при загрузке переопределения клиники:', err);
          }
      },

      async loadDoctorOverride(doctorId, date) {
          this.overrideErrors.date = false;
          const datePattern = /^\d{4}-\d{2}-\d{2}$/;

          if (!date || !doctorId || !datePattern.test(date)) return;

          try {
              const res = await fetch(`/api/doctor-overrides/${doctorId}/${date}`);
              if (!res.ok) {
                  if (res.status === 404) {
                      this.overrideForm.type = 'off';
                      this.overrideForm.start_time = '';
                      this.overrideForm.end_time = '';
                      return;
                  }
                  console.error('Ошибка при загрузке переопределения врача:', res.statusText);
                  return;
              }

              const data = await res.json();
              this.overrideForm.type = data.is_day_off ? 'off' : 'work';
              this.overrideForm.start_time = data.start_time || '';
              this.overrideForm.end_time = data.end_time || '';
          } catch (err) {
              console.error('Ошибка при загрузке переопределения врача:', err);
          }
      },

    handleClickOutside(event) {
         const popover = document.getElementById('admin-profile');
         if (popover && !popover.contains(event.target)) {
        this.isPopoverVisible = false;
        }
    },

   async loadClinicSchedule() {
  try {
    const res = await fetch('/api/clinic-schedule');
    const data = await res.json();

      const fullSchedule = Array(7).fill(null).map((_, i) => ({
          day: i,
          start_time: '',
          end_time: '',
          is_day_off: false
      }));

      for (const day of data.schedule) {
          fullSchedule[day.day] = { ...day, day: day.day };
      }

    this.clinicSchedule = fullSchedule;
    this.clinicSlotMinutes = data.slot_minutes;
  } catch (err) {
    console.error('Ошибка при загрузке расписания клиники', err);
  }
},

    validateTimeInterval(start, end, minDuration) {
        if (!start || !end) return false;
        const [sh, sm] = start.split(':').map(Number);
        const [eh, em] = end.split(':').map(Number);
        const startMin = sh * 60 + sm;
        const endMin = eh * 60 + em;
        return endMin > startMin && (endMin - startMin) >= minDuration;
    },

    isValidSlot(slot) {
        return slot >= 5 && slot <= 180;
    },


    isValidDateOverride(dateStr) {
        const today = new Date();
        const date = new Date(dateStr);
        const in3Months = new Date();
        in3Months.setMonth(in3Months.getMonth() + 3);
        return date >= today.setHours(0,0,0,0) && date <= in3Months;
    },

    clinicTimeBoundsForDay(day) {
        const clinicDay = this.clinicSchedule.find(d => d.day === day);
        if (!clinicDay || !clinicDay.is_day_off) return null; // если день не рабочий
        return {
            start: clinicDay.start_time,
            end: clinicDay.end_time,
        };
    },

    isWithinClinicTime(start, end, clinicBounds) {
        if (!clinicBounds) {
            this.overrideErrors.date = true;
            return false
        };
        return (
            start >= clinicBounds.start &&
            end <= clinicBounds.end &&
            this.validateTimeInterval(start, end, this.clinicSlotMinutes)
        );
    },


    async saveClinicSchedule() {
        this.invalidClinicDays = [];
        this.clinicSlotError = false;

        if (!this.isValidSlot(this.clinicSlotMinutes)) {
            this.clinicSlotError = true;
            alert('Продолжительность приема должна составлять от 10 до 180 минут');
            return;
        }

        for (const day of this.clinicSchedule) {
            if (day.is_day_off && !this.validateTimeInterval(day.start_time, day.end_time, this.clinicSlotMinutes)) {
            this.invalidClinicDays.push(day.day);
            }
        }

        if (this.invalidClinicDays.length > 0) {
            alert('Проверьте корректность расписания клиники');
            return;
        }


        try {
            const res = await fetch('api/clinic-schedule', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({schedule: this.clinicSchedule,
          slot_duration_minutes: this.clinicSlotMinutes}),
        });
        if (!res.ok) throw new Error('Ошибка при сохранении');
        alert('Расписание клиники сохранено');
      } catch (err) {
        console.error(err);
        alert('Не удалось сохранить расписание клиники');
      }
    },
    async loadDoctors() {
      try {
        const res = await fetch('/api/doctors');
        this.doctors = await res.json();
      } catch (err) {
        console.error('Ошибка при загрузке списка врачей', err);
      }
    },

      async fetchDoctorSchedule() {
          if (!this.selectedDoctor) return;

          try {
              const res = await fetch(`/api/doctor-schedule/${this.selectedDoctor.user_id}`);
              const data = await res.json();
              console.log("Ответ от сервера", data);

              const fullSchedule = Array(7).fill(null).map((_, i) => ({
                  day: i,
                  start_time: '',
                  end_time: '',
                  is_day_off: false // все дни по умолчанию выходные
              }));

              if (Array.isArray(data.schedule)) {
                  for (const day of data.schedule) {
                      fullSchedule[day.day] = {
                          ...fullSchedule[day.day],
                          ...day,
                          day: day.day,
                          is_day_off: day.is_day_off !== undefined ? day.is_day_off : true
                      };
                  }
              }

              this.doctorSchedule = fullSchedule;
              this.doctorSlotMinutes = data.slot_minutes ?? null; // если нет, будет null
          } catch (err) {
              console.error('Ошибка при загрузке расписания врача', err);
              // fallback: показать пустое расписание, где всё выходные
              this.doctorSchedule = Array(7).fill(null).map((_, i) => ({
                  day: i,
                  start_time: '',
                  end_time: '',
                  is_day_off: true
              }));
              this.doctorSlotMinutes = null;
          }
      },
    async saveDoctorSchedule() {
        this.invalidDays = [];
        this.doctorSlotError = false;

        if (!this.isValidSlot(this.doctorSlotMinutes)) {
            this.doctorSlotError = true;
            alert('Продолжительность приема должна составлять от 10 до 180 минут');
            return;
        }

        for (const day of this.doctorSchedule) {
            if (day.is_day_off) { // если день рабочий
            const clinicBounds = this.clinicTimeBoundsForDay(day.day);

            if (!this.validateTimeInterval(day.start_time, day.end_time, this.doctorSlotMinutes)) {
                this.invalidDays.push(day.day);
                continue;
            }

            if (!this.isWithinClinicTime(day.start_time, day.end_time, clinicBounds)) {
                this.invalidDays.push(day.day);
                continue;
            }
            }
        }

        if (this.invalidDays.length > 0) {
            alert('Проверьте корректность расписания врача');
            return;
        }

    try {
        const res = await fetch(`/api/doctor-schedule/${this.selectedDoctor.user_id}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            schedule: this.doctorSchedule,
            slot_minutes: this.doctorSlotMinutes,
        }),
        });
        if (!res.ok) throw new Error('Ошибка при сохранении');
        alert('Расписание врача сохранено');
    } catch (err) {
        console.error(err);
        alert('Не удалось сохранить расписание врача');
    }
    },

      async handleOverrideClinicDateChange() {
          this.overrideClinicErrors.date = false;

          const datePattern = /^\d{4}-\d{2}-\d{2}$/;

          if (!this.overrideClinicForm.date || !datePattern.test(this.overrideClinicForm.date)) {
              return;
          }

          try {
              const res = await fetch(`/api/clinic-overrides/${this.overrideClinicForm.date}`);
              if (!res.ok) {
                  if (res.status === 404) {
                      this.overrideClinicForm.type = 'off';
                      this.overrideClinicForm.start_time = '';
                      this.overrideClinicForm.end_time = '';
                      return;
                  }
                  console.error('Ошибка при загрузке переопределения клиники:', res.statusText);
                  return;
              }

              const data = await res.json();

              this.overrideClinicForm.type = data.is_day_off ? 'off' : 'work';
              this.overrideClinicForm.start_time = data.start_time || '';
              this.overrideClinicForm.end_time = data.end_time || '';
          } catch (err) {
              console.error('Ошибка при загрузке переопределения клиники:', err);
          }
      },

      async handleOverrideDateChange() {
          this.overrideErrors.date = false;

          const datePattern = /^\d{4}-\d{2}-\d{2}$/;

          if (!this.overrideForm.date || !this.overrideForm.doctor_id || !datePattern.test(this.overrideForm.date)) {
              return;
          }
          try {
              const res = await fetch(`/api/doctor-overrides/${this.overrideForm.doctor_id}/${this.overrideForm.date}`);
              if (!res.ok) {
                  if (res.status === 404) {
                      this.overrideForm.type = 'off';
                      this.overrideForm.start_time = '';
                      this.overrideForm.end_time = '';
                      return;
                  }
                  console.error('Ошибка при загрузке переопределения клиники:', res.statusText);
                  return;
              }



              const data = await res.json();

              this.overrideForm.type = data.is_day_off ? 'off' : 'work';
              this.overrideForm.start_time = data.start_time || '';
              this.overrideForm.end_time = data.end_time || '';
          } catch (err) {
              console.error('Ошибка при загрузке переопределения врача:', err);
          }
      },

      async saveOverrideClinic() {
        this.overrideClinicErrors = {
            date: false,
            start_time: false,
            end_time: false,
        };

        if (!this.isValidDateOverride(this.overrideClinicForm.date)) {
            this.overrideClinicErrors.date = true;
            alert('Недопустимая дата. Выберите дату не позднее трех месяцев от текущей.');
            return;
        }

        if (this.overrideClinicForm.type === 'work') {
            if (!this.validateTimeInterval(this.overrideClinicForm.start_time, this.overrideClinicForm.end_time, 0)) {
            this.overrideClinicErrors.start_time = true;
            this.overrideClinicErrors.end_time = true;
            alert('Недопустимый интервал времени');
            return;
            }
        }

        try {
            const res = await fetch('/api/clinic-overrides', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(this.overrideClinicForm),
            });

            if (!res.ok) throw new Error('Ошибка при сохранении переопределения');
            alert('Переопределение сохранено');

            this.overrideClinicForm = {
            date: '',
            type: 'off',
            start_time: '',
            end_time: ''
            };
        } catch (err) {
            console.error(err);
            alert('Не удалось сохранить переопределение');
        }
        }, 

    async saveOverride() {
        this.overrideErrors = {
            date: false,
            start_time: false,
            end_time: false,
            slot_minutes: false
        };


        if (!this.isValidDateOverride(this.overrideForm.date)) {
             this.overrideErrors.date = true;
             alert('Недопустимая дата. Выберите дату не позднее трех месяцев от текущей.');
            return;
    }

    if (this.overrideForm.type === 'work') {
        if (!this.validateTimeInterval(this.overrideForm.start_time, this.overrideForm.end_time, 0)) {
            this.overrideErrors.start_time = true;
            this.overrideErrors.end_time = true;
            alert('Недопустимый интервал времени');
            return;
            }

        const day = new Date(this.overrideForm.date).getDay(); // 0 = воскресенье
        const clinicBounds = this.clinicTimeBoundsForDay(day);


        if (!this.isWithinClinicTime(this.overrideForm.start_time, this.overrideForm.end_time, clinicBounds)) {
                this.overrideErrors.start_time = true;
                this.overrideErrors.end_time = true;
                alert('Выбранные параметры не соответствуют расписанию клиники');
                return;
            }
        }
      try {
        const res = await fetch('/api/doctor-overrides', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(this.overrideForm),
        });
        if (!res.ok) throw new Error('Ошибка при сохранении переопределения');
        alert('Переопределение сохранено');
        this.overrideForm = {
          doctor_id: null,
          date: '',
          type: 'off',
          start_time: '',
          end_time: '',
          slot_minutes: null,
        };
      } catch (err) {
        console.error(err);
        alert('Не удалось сохранить переопределение');
      }
    },
  }
}).mount('#app');
