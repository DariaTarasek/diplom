const { createApp, ref, reactive, onMounted, onUnmounted, computed, watch, nextTick } = Vue;

createApp({
  setup() {
    const specialties = ref([]);
    const doctors = ref([]);
    const schedule = ref({});
    const maxSlots = ref(0);
    const selectedDoctorId = ref(null);
    const selectedSpecialization = ref(null);
    const selectedDoctor = ref(null);
    const selectedSlot = ref(null);
    const selectedDate = ref(null);
    const step = ref(1);

    const authFirstName = ref("");
    const authSecondName = ref("");

    const isPopoverVisible = ref(false)

    const isAuthorized = ref(false)

    const patient = reactive({
      user_id: '',
      secondName: '',
      firstName: '',
      surname: '',
      birthDate: '',
      gender: '',
      phone: ''
    });

    const errors = reactive({
      phone: '',
      firstName: '',
      secondName: ''
    });

    // Пагинация по неделям
    const currentPage = ref(0);
    const daysPerPage = 7;

    const totalPages = computed(() => {
      return Math.ceil(schedule.value.length / daysPerPage);
    });

    const paginatedSchedule = computed(() => {
      const start = currentPage.value * daysPerPage;
      const end = start + daysPerPage;
      return schedule.value.slice(start, end);
    });

    const nextPage = () => {
      if (currentPage.value < totalPages.value - 1) {
        currentPage.value++;
      }
    };

    const prevPage = () => {
      if (currentPage.value > 0) {
        currentPage.value--;
      }
    };




    const validatePhone = () => {
      if (errors.phone) {
        return false;
      } else {
        return true;
      }
    };

    const validateFirstName = () => {
      errors.first_name = patient.firstName.trim() ? '' : 'Имя не может быть пустым';
      return !errors.firstName;
    };

    const validateSecondName = () => {
      errors.second_name = patient.secondName.trim() ? '' : 'Фамилия не может быть пустой';
      return !errors.secondName;
    };


    const fetchSpecialties = async () => {
      const res = await fetch('/api/specialties');
      specialties.value = await res.json();
    };

    const fetchDoctors = async (specialtyId) => {
      selectedDoctorId.value = "";
      doctors.value = [];
  const res = await fetch(`/api/doctors/${specialtyId}`);
  const rawDoctors = await res.json();

 
  doctors.value = rawDoctors.map(doc => ({
    ...doc,
    fullName: `${doc.secondName} ${doc.firstName} ${doc.surname}`.trim()
  }));
};


    const fetchSchedule = async (doctorId) => {
      const res = await fetch(`/api/appointment-doctor-schedule/${doctorId}`);
      const data = await res.json();

      // гарантируем, что slots всегда массив
      data.forEach(day => {
        if (!Array.isArray(day.slots)) {
          day.slots = [];
        }
      });

      schedule.value = data;

      maxSlots.value = Math.max(...schedule.value.map(day => day.slots.length));
    };


    const selectDateSlot = (date, slot) => {
      selectedDate.value = date.label;
      selectedSlot.value = slot;
      step.value = 2;
    };

    const back = () => {
      step.value = 1;
      selectedSlot.value = null;
    };

    const submitForm = async () => {
      if (!selectedSlot.value || !selectedDoctorId.value) {
        alert('Выберите врача и время');
        return;
      }

      if (!validatePhone() || !validateFirstName() || !validateSecondName) return;

const phoneInput = document.getElementById('phone');
if (phoneInput) {
  patient.phone = phoneInput.value; 
}
patient.phone = patient.phone.replace(/\D/g, '')

      const payload = {
        doctor_id: selectedDoctorId.value,
        date: selectedDate.value,
        time: selectedSlot.value,
      ...patient
      };

      const res = await fetch('/api/appointments', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      });

      if (res.ok) {
        alert('Запись успешно создана!');
        if (isAuthorized.value === true) {
             window.location.href = "/patient_account.html";
        } else {
             window.location.href = "/index.html";
        }
      } else {
        alert('Ошибка при записи. Попробуйте позже.');
      }
    };

    const loadPatientData = async () => {
      try {
        const res = await fetch('/api/patient/me');
        if (!res.ok) return;
        isAuthorized.value = true;
        const data = await res.json();
        authFirstName.value = data.firstName;
        authSecondName.value = data.secondName;
        Object.assign(patient, {
          user_id: data.user_id || '',
          secondName: data.secondName || '',
          firstName: data.firstName || '',
          surname: data.surname || '',
          birthDate: data.birthDate || '',
          gender: data.gender || '',
          phone: data.phone || ''
        });
      } catch {
        console.warn('Пациент не авторизован');
      }
    };

    const fullName = computed(() => {
  return [authFirstName.value, authSecondName.value].filter(Boolean).join(' ');
    });

const birthDateAttrs = reactive({ min: '', max: '' });

const validateDateRange = () => {
  const today = new Date();
  const yyyy = today.getFullYear();
  const mm = String(today.getMonth() + 1).padStart(2, '0');
  const dd = String(today.getDate()).padStart(2, '0');
  birthDateAttrs.min = `${yyyy - 110}-${mm}-${dd}`;
  birthDateAttrs.max = `${yyyy - 18}-${mm}-${dd}`;
};


function handleClickOutside(event) {
            const popover = document.getElementById('patient-profile');
            if (popover && !popover.contains(event.target)) {
                isPopoverVisible.value = false;
            }
        };

onMounted(() => {
      fetchSpecialties();
      loadPatientData();
      validateDateRange();
      document.addEventListener('click', handleClickOutside);
    });

onUnmounted(() => {
  if (handleClickOutside) {
    document.removeEventListener('click', handleClickOutside);
  }

  const phoneInput = document.getElementById('phone');
  if (maskInstance && maskInstance.destroy) {
    maskInstance.destroy();
    maskInstance = null;
  }

  if (phoneInput) {
    delete phoneInput.dataset.masked;
  }
});

    
watch(step, (newVal) => {
  if (newVal === 2) {
    nextTick(() => {
      const phoneInput = document.getElementById('phone');
      if (phoneInput && !phoneInput.dataset.masked) {
        maskInstance = IMask(phoneInput, {
          mask: '+{7} (000) 000-00-00'
        });

        phoneInput.dataset.masked = "true";

        maskInstance.on('accept', () => {
          const digits = maskInstance.value.replace(/\D/g, '');
          const valid = digits.length === 11 && digits.startsWith('7');
          errors.phone = digits && !valid ? 'Неверный формат телефона' : '';
        });

        maskInstance.on('complete', () => {
          patient.phone = maskInstance.value;
        });
      }
    });
  }
});



watch(() => patient.firstName, () => {
  validateFirstName();
});

watch(() => patient.secondName, () => {
  validateSecondName();
});

    return {
      specialties,
      doctors,
      schedule,
      maxSlots,
      selectedDoctorId,
      selectedSpecialization,
      selectedDoctor,
      selectedSlot,
      step,
      patient,
      errors,
      fetchDoctors,
      fetchSchedule,
      selectDateSlot,
      back,
      submitForm,
      validatePhone,
      isAuthorized,
      isPopoverVisible,
      fullName,
      birthDateAttrs,
      currentPage,
      totalPages,
      paginatedSchedule,
      nextPage,
      prevPage
    };
  }
}).mount("#app");
