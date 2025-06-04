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
                firstName: '',
                secondName: ''
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
            documents: [],
            showDocumentsModal: false,
        };
    },

    computed: {
        fullName() {
            return [this.patient.second_name, this.patient.first_name, this.patient.surname].filter(Boolean).join(' ');
        },
        patientAge() {
            if (!this.patient.birthDate) return '';

            const parts = this.patient.birthDate.split('.');
            if (parts.length !== 3) return '';

            const day = parseInt(parts[0], 10);
            const month = parseInt(parts[1], 10) - 1; // Месяцы начинаются с 0
            const year = parseInt(parts[2], 10);

            const birth = new Date(year, month, day);
            if (isNaN(birth)) return '';

            const now = new Date();
            let age = now.getFullYear() - birth.getFullYear();
            const hasBirthdayPassed =
                now.getMonth() > birth.getMonth() ||
                (now.getMonth() === birth.getMonth() && now.getDate() >= birth.getDate());

            if (!hasBirthdayPassed) {
                age--;
            }

            return age;
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
            return [this.doctor.firstName, this.doctor.secondName].filter(Boolean).join(' ');
        },
        visibleICDCodes() {
            return this.icdCodes;
        },
        filteredServices() {
            const term = this.serviceSearch.toLowerCase();
            return this.services.filter(service => service.name.toLowerCase().includes(term));
        },
        filteredMaterials() {
            const term = this.materialSearch.toLowerCase();
            return this.materials.filter(material => material.name.toLowerCase().includes(term));
        },
        combinedAllergies() {
            return [
                ...this.patient.allergies.map(a => typeof a === 'string' ? a : a.title),
                ...this.newAllergies.map(a => typeof a === 'string' ? a : a.title)
            ];
        },

        combinedChronics() {
            return [
                ...this.patient.chronics.map(c => typeof c === 'string' ? c : c.title),
                ...this.newChronics.map(c => typeof c === 'string' ? c : c.title)
            ];
        },
    },

    methods: {
        async loadData() {
            const urlParams = new URLSearchParams(window.location.search);
            const appointmentId = urlParams.get('appointment_id') || 0;

            try {
                const res = await fetch(`/api/appointments/${appointmentId}`);
                const data = await res.json();
                this.patient = {
                    id: data.id,
                    first_name: data.firstName,
                    second_name: data.secondName,
                    surname: data.surname,
                    birthDate: data.birthDate,
                    gender: data.gender,
                    allergies: [],
                    chronics: []
                };
            } catch (err) {
                console.error('Ошибка загрузки пациента:', err);
            }

            try {
                //this.patient.id = 20;
                const condRes = await fetch(`/api/patient-notes/${this.patient.id}`);
                const conditions = await condRes.json();
                this.patient.allergies = conditions.filter(c => c.type === 'allergy');
                this.patient.chronics = conditions.filter(c => c.type === 'chronic');
            } catch (err) {
                console.error('Ошибка загрузки мед. записей:', err);
            }

            try {
                const historyRes = await fetch(`/api/patient-history/${this.patient.id}`);
                this.history = await historyRes.json();
                this.history = this.history.map(entry => ({
                    id: entry.id,
                    date: entry.created_at,
                    diagnosis: entry.diagnoses
                        .map(d => `${d.icd_code}${d.notes ? ' (' + d.notes + ')' : ''}`)
                        .join('; '),
                    treatment: entry.treatment,
                    doctor: entry.doctor || ''
                }));
            } catch (err) {
                console.error('Ошибка загрузки анамнеза пациента:', err);
            }

            try {
                const [servicesRes, materialsRes] = await Promise.all([
                    fetch('/api/services'),
                    fetch('/api/materials')
                ]);

                const rawServices = await servicesRes.json();
                const rawMaterials = await materialsRes.json();

                this.services = Array.isArray(rawServices.services) ? rawServices.services : [];
                this.materials = Array.isArray(rawMaterials.materials) ? rawMaterials.materials : [];

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

        getYearWord(age) {
            const lastDigit = age % 10;
            const lastTwoDigits = age % 100;

            if (lastTwoDigits >= 11 && lastTwoDigits <= 14) return 'лет';
            if (lastDigit === 1) return 'год';
            if (lastDigit >= 2 && lastDigit <= 4) return 'года';
            return 'лет';
        },

        fetchDoctorData() {
            fetch('/api/doctor/me')
                .then(res => {
                    if (!res.ok) throw new Error('Ошибка загрузки');
                    return res.json();
                })
                .then(json => {
                    this.doctor.firstName = json.firstName;
                    this.doctor.secondName = json.secondName;
                })
                .catch(err => {
                    console.error('Ошибка при получении данных:', err);
                });
        },

        async fetchDocuments() {
            if (!this.patient.id) return;
            try {
                const res = await fetch(`/api/doctor/consultation/patient-tests/${this.patient.id}`);
                const rawDocs = await res.json();
                this.documents = Array.isArray(rawDocs)
                    ? rawDocs.map(doc => ({
                        id: doc.id,
                        date: doc.date || doc.created_at || '—',
                        description: doc.description || 'Без описания'
                    }))
                    : [];
            } catch (err) {
                console.error('Ошибка загрузки документов:', err);
            }
        },


        toggleMaterial(mat) {
            this.materialQuantities[mat.id] = this.materialQuantities[mat.id] > 0 ? 0 : 1;
        },

        isMaterialSelected(id) {
            return this.materialQuantities[id] > 0;
        },

        removeService(id) {
            this.currentVisit.selectedServices = this.currentVisit.selectedServices.filter(sid => sid !== id);
        },

        removeNewAllergyByValue(title) {
            this.newAllergies = this.newAllergies.filter(a => a !== title);
        },

        removeNewChronicByValue(title) {
            this.newChronics = this.newChronics.filter(c => c !== title);
        },

        onMaterialToggle(mat) {
            this.materialQuantities[mat.id] = this.selectedMaterialIds.includes(mat.id) ? 1 : 0;
        },

        removeMaterial(id) {
            this.selectedMaterialIds = this.selectedMaterialIds.filter(mid => mid !== id);
            this.materialQuantities[id] = 0;
        },

        addICDCode(icd) {
            if (!this.selectedICD.some(item => item.code === icd.code)) {
                this.selectedICD.push({ id: icd.id, code: icd.code, description: icd.description, comment: '' });
            }
        },

        removeICDCode(code) {
            this.selectedICD = this.selectedICD.filter(item => item.code !== code);
        },

        async fetchICDCodes() {
            try {
                const res = await fetch(`/api/icd-codes`);
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
            if (
                trimmed &&
                !this.patient.allergies.some(a => a.title === trimmed) &&
                !this.newAllergies.includes(trimmed)
            ) {
                this.newAllergies.push(trimmed);
                this.newAllergy = '';
            }
        },

        addNewChronic() {
            const trimmed = this.newChronic.trim();
            if (
                trimmed &&
                !this.patient.chronics.some(c => c.title === trimmed) &&
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
                    appointment_id: this.appointment_id || 19,
                    patient_id: this.patient.id,
                    doctor_id: this.doctor_id,
                    complaints: this.currentVisit.complaints,
                    treatment: this.currentVisit.treatment,
                    manipulations: this.currentVisit.selectedServices.map(id => ({
                        id,
                        quantity: this.serviceQuantities[id] || 1
                    })),
                    materials: this.selectedMaterialIds.map(id => ({
                        id,
                        quantity: this.materialQuantities[id]
                    })),
                    icd_codes: this.selectedICD.map(item => ({
                        code: item.id,
                        comment: item.comment
                    }))
                };

                const res = await fetch(`/api/visits`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(payload)
                });

                if (!res.ok) throw new Error('Ошибка при сохранении приема');

                const newConditions = [];
                this.newAllergies.forEach(title => newConditions.push({ type: 'allergy', title }));
                this.newChronics.forEach(title => newConditions.push({ type: 'chronic', title }));

                if (newConditions.length > 0) {
                    await fetch(`/api/patient-notes/${this.patient.id}`, {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify(newConditions)
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

    async mounted() {
        await this.loadData();
        this.fetchDoctorData();
        this.fetchICDCodes();
        this.fetchDocuments();
    },
}).mount('#app');