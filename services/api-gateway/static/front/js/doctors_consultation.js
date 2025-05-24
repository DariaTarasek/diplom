const { createApp } = Vue;

createApp({
  data() {
    return {
      tabs: [
        { id: 'patient', label: 'Пациент' },
        { id: 'history', label: 'Анамнез' },
        { id: 'current', label: 'Прием' },
        { id: 'materials', label: 'Материалы' }
      ],
      activeTab: 'patient',
      patient: {
        id: null,
        first_name: '',
        second_name: '',
        surname: '',
        birthDate: '',
        gender: '',
        allergies: [],
        chronics: []
      },
      doctor: {
        first_name: '',
        second_name: ''
      },
      newAllergy: '',
      newAllergies: [],
      newChronic: '',
      newChronics: [],

      history: [],
      currentVisit: {
        complaints: '',
        treatment: '',
        selectedServices: []
      },

      services: [],
      materials: [],
      selectedMaterialIds: [],
      serviceSearch: '',
      materialSearch: '',
      materialQuantities: {},
      serviceQuantities: {},
      isPopoverVisible: false,
      icdSearch: '',
      icdCodes: [], 
      selectedICD: [], 
      showServiceModal: false,
      showMaterialModal: false,
      showICDModal: false,
      tempSelectedServices: [],
      tempSelectedMaterials: [],
    };
  },

  computed: {
    fullName() {
      return [this.patient.second_name, this.patient.first_name, this.patient.surname].filter(Boolean).join(' ');
    },
    patientAge() {
      const birth = new Date(this.patient.birthDate);
      const now = new Date();
      return now.getFullYear() - birth.getFullYear() - (now < new Date(now.getFullYear(), birth.getMonth(), birth.getDate()) ? 1 : 0);
    },
    totalPrice() {
      const serviceTotal = this.services
        .filter(s => this.currentVisit.selectedServices.includes(s.id))
        .reduce((sum, s) => sum + s.price, 0);

      const materialTotal = Object.entries(this.materialQuantities).reduce((sum, [id, qty]) => {
        const mat = this.materials.find(m => m.id == id);
        return mat ? sum + mat.price * qty : sum;
      }, 0);

      return serviceTotal + materialTotal;
    },
    filteredICDCodes() {
        const term = this.icdSearch.toLowerCase();
        return this.icdCodes.filter(icd =>
            icd.code.toLowerCase().includes(term) || icd.description.toLowerCase().includes(term)
        );
    },
    selectedServicesDetails() {
        return this.services.filter(s => this.currentVisit.selectedServices.includes(s.id));
    },
selectedMaterials() {
    return this.materials.filter(m => this.selectedMaterialIds.includes(m.id));
    },
      accountName() {
        return [
            this.doctor.first_name,
            this.doctor.second_name
        ]
        .filter(Boolean)
        .join(' ');
    },
    visibleICDCodes() {
  return this.icdCodes;
},
    filteredServices() {
        const term = this.serviceSearch.toLowerCase();
        return this.services.filter(service =>
            service.name.toLowerCase().includes(term)
        );
    },
    filteredMaterials() {
        const term = this.materialSearch.toLowerCase();
        return this.materials.filter(material =>
            material.name.toLowerCase().includes(term)
        );
    },

  combinedAllergies() {
    return [...this.patient.allergies, ...this.newAllergies];
  },
  combinedChronics() {
    return [...this.patient.chronics, ...this.newChronics];
  },
  },

  methods: {
    async loadData() {
        const urlParams = new URLSearchParams(window.location.search);
        const appointmentId = urlParams.get('appointment_id') || 1;

        try {
            const res = await fetch(`http://192.168.1.207:8080/api/appointments/${appointmentId}`);
            const data = await res.json();
            this.patient = data.patient || {};
            this.patient.allergies = this.patient.allergies || [];
            this.patient.chronics = this.patient.chronics || [];
            this.history = data.history || [];
        } catch (err) {
            console.error('Ошибка загрузки пациента:', err);
        }

        try {
            const [servicesRes, materialsRes] = await Promise.all([
            fetch('http://192.168.1.207:8080/api/services'),
            fetch('http://192.168.1.207:8080/api/materials')
            ]);

            const rawServices = await servicesRes.json();
            const rawMaterials = await materialsRes.json();

            this.services = Array.isArray(rawServices.services) ? rawServices.services : [];

            this.materials = Array.isArray(rawMaterials.materials) ? rawMaterials.materials : [];

            if (!Array.isArray(this.services)) {
            console.error('services не массив:', this.services);
            }

            if (!Array.isArray(this.materials)) {
            console.error('materials не массив:', this.materials);
            }

            this.services.forEach(s => {
                this.serviceQuantities[s.id] = 1; 
                });


            this.materials.forEach(mat => {
                this.materialQuantities[mat.id] = 0;
            });

        } catch (err) {
            console.error('Ошибка загрузки услуг или материалов:', err);
        }
        },

     fetchDoctorData() {
      fetch('http://192.168.1.207:8080/api/doctor-data')
        .then(res => {
          if (!res.ok) throw new Error('Ошибка загрузки');
          return res.json();
        })
        .then(json => {
          this.doctor.first_name = json.first_name;
          this.doctor.second_name = json.second_name;
        })
        .catch(err => {
          console.error('Ошибка при получении данных:', err);
        });
    },
      toggleMaterial(mat) {
    if (this.materialQuantities[mat.id] > 0) {
      this.materialQuantities[mat.id] = 0;

    } else {
      this.materialQuantities[mat.id] = 1;

    }
  },
  isMaterialSelected(id) {
    return this.materialQuantities[id] > 0;
  },

  removeService(id) {
  this.currentVisit.selectedServices = this.currentVisit.selectedServices.filter(sid => sid !== id);
},
removeNewAllergyByValue(allergy) {
  this.newAllergies = this.newAllergies.filter(a => a !== allergy);
},
 removeNewChronicByValue(chronic) {
    this.newChronics = this.newChronics.filter(c => c !== chronic);
  },
onMaterialToggle(mat) {
    if (this.selectedMaterialIds.includes(mat.id)) {
        this.materialQuantities[mat.id] = 1;
    } else {
        this.materialQuantities[mat.id] = 0;
    }
  },
  removeMaterial(id) {
    this.selectedMaterialIds = this.selectedMaterialIds.filter(mid => mid !== id);
    this.materialQuantities[id] = 0;
  },

addICDCode(icd) {
  if (!this.selectedICD.some(item => item.code === icd.code)) {
    this.selectedICD.push({ code: icd.code, description: icd.description, comment: '' });
  }
},

removeICDCode(code) {
  this.selectedICD = this.selectedICD.filter(item => item.code !== code);
},
async fetchICDCodes() {
    try {
      const res = await fetch('http://192.168.1.207:8080/api/icd-codes');
      this.icdCodes = await res.json();
    } catch (err) {
      console.error('Ошибка загрузки МКБ кодов:', err);
    }
  },
  applySelectedServices() {
    this.tempSelectedServices.forEach(id => {
      if (!this.currentVisit.selectedServices.includes(id)) {
        this.currentVisit.selectedServices.push(id);
      }
    });
    this.tempSelectedServices = [];
    this.showServiceModal = false;
  },
  applySelectedMaterials() {
    this.tempSelectedMaterials.forEach(id => {
      if (!this.selectedMaterialIds.includes(id)) {
        this.selectedMaterialIds.push(id);
        const mat = this.materials.find(m => m.id === id);
        this.selectedMaterials.push(mat);
        this.materialQuantities[id] = 1;
      }
    });
    this.tempSelectedMaterials = [];
    this.showMaterialModal = false;
  },
  selectMaterial(mat) {
  if (!this.selectedMaterialIds.includes(mat.id)) {
    this.selectedMaterialIds.push(mat.id);
    this.materialQuantities[mat.id] = 1;
  }
},
  addNewAllergy() {
        const trimmed = this.newAllergy.trim();
        if (trimmed && !this.patient.allergies.includes(trimmed) && !this.newAllergies.includes(trimmed)) {
            this.newAllergies.push(trimmed);
            this.newAllergy = '';
        }
    },

    removeNewAllergy(index) {
        this.newAllergies.splice(index, 1);
    },
    addNewChronic() {
        const trimmed = this.newChronic.trim();
        if (
            trimmed &&
            !this.patient.chronics.includes(trimmed) &&
            !this.newChronics.includes(trimmed)
        ) {
            this.newChronics.push(trimmed);
            this.newChronic = '';
        }
     },
selectService(service) {
  if (!this.currentVisit.selectedServices.includes(service.id)) {
    this.currentVisit.selectedServices.push(service.id);
    this.serviceQuantities[service.id] = 1;
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

    async submitVisit() {
        const errors = [];

        if (!this.currentVisit.complaints.trim()) {
            errors.push('Заполните поле "Жалобы".');
        }
        if (!this.currentVisit.treatment.trim()) {
            errors.push('Заполните поле "Лечение".');
        }
        if (this.currentVisit.selectedServices.length === 0) {
            errors.push('Выберите хотя бы одну услугу.');
        }

        for (const materialId of this.selectedMaterialIds) {
    const qty = this.materialQuantities[materialId];
    if (!Number.isInteger(qty) || qty <= 0) {
        errors.push('Количество для выбранных материалов должно быть больше 0');
        break;
    }
}

if (errors.length) {
    alert('Пожалуйста, исправьте ошибки:\n\n' + errors.join('\n'));
    return;
}

try {
    const payload = {
        patient_id: this.patient.id,
        complaints: this.currentVisit.complaints,
        treatment: this.currentVisit.treatment,
        manipulations: this.currentVisit.selectedServices,
        materials: this.selectedMaterialIds.map(id => ({
            id,
            quantity: this.materialQuantities[id]
        }))
    };

    payload.icd_codes = this.selectedICD.map(item => ({
        code: item.code,
        comment: item.comment
    }));


    const res = await fetch(`http://192.168.1.207:8080/api/visits`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
            });

            if (!res.ok) throw new Error('Ошибка при сохранении приема');

    
            if (this.newAllergies.length > 0) {
                await fetch(`http://192.168.1.207:8080/api/patients/${this.patient.id}/allergies`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ allergies: this.newAllergies })
                });
            }

    
            if (this.newChronics.length > 0) {
                await fetch(`http://192.168.1.207:8080/api/patients/${this.patient.id}/chronics`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ chronics: this.newChronics })
                });
            }

            alert('Прием сохранен');
            window.location.href = '/doctor_account.html';

        } catch (err) {
            console.error(err);
            alert('Ошибка при сохранении.');
        }
    }

},


  mounted() {
    this.loadData();
    this.fetchDoctorData();
    this.fetchICDCodes();
   
  }
}).mount('#app');
