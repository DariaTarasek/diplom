const { createApp } = Vue;

createApp({
    data() {
        return {
            activeTab: 'services',
            tabs: [
                { id: 'services', label: 'Услуги' },
                { id: 'materials', label: 'Материалы' }
            ],
            services: [],
            materials: [],
            searchQuery: '',
            selectedCategory: '',
            serviceCategories: [],
            isPopoverVisible: false,
            first_name: '',
            second_name: '',
            selectedService: null,
            hoveredServiceId: null,
            isEditingService: false,
            editedServicePrice: 0,
            selectedMaterial: null,
            hoveredMaterialId: null,
            isEditingMaterial: false,
            editedMaterialPrice: 0,
            newItem: {
                name: '',
                category_id: '',
                price: 0
            },
            addError: '',
            editServiceError: '',
            editMaterialError: '',
            addNameError: '',
            addCategoryError: '',
            addPriceError: '',
            userRole: ''
        };
    },

    computed: {
        fullName() {
            return [this.first_name, this.second_name].filter(Boolean).join(' ');
        },
        filteredServices() {
            return this.services.filter(service => {
                const matchesQuery = service.name.toLowerCase().includes(this.searchQuery.toLowerCase());
                const matchesCategory = this.selectedCategory ? service.category_id === this.selectedCategory : true;
                return matchesQuery && matchesCategory;
            });
        },
        filteredMaterials() {
            return this.materials.filter(material => {
                return material.name.toLowerCase().includes(this.searchQuery.toLowerCase());
            });
        }
    },

    methods: {
        async fetchData() {
            try {
                const userRes = await fetch('/api/admin-data');
                const userData = await userRes.json();
                this.first_name = userData.first_name || '';
                this.second_name = userData.second_name || '';
                this.userRole = userData.role || '';

                await this.fetchServices();
                await this.fetchMaterials();
            } catch (err) {
                console.error('Ошибка при загрузке данных:', err);
            }
        },

        async fetchServices() {
            try {
                const res = await fetch(`/api/services`);
                const data = await res.json();
                this.services = data.services || [];
            } catch (err) {
                console.error('Ошибка при загрузке услуг:', err);
            }
        },

        async fetchMaterials() {
            try {
                const res = await fetch(`/api/materials`);
                const data = await res.json();
                this.materials = data.materials || [];
            } catch (err) {
                console.error('Ошибка при загрузке материалов:', err);
            }
        },

        async fetchServiceCategories() {
            try {
                const res = await fetch('/api/service-categories');
                const data = await res.json();
                this.serviceCategories = data.categories || [];
            } catch (err) {
                console.error('Ошибка при загрузке категорий:', err);
            }
        },

        openServiceModal(service) {
            this.selectedService = {
                id: service.id,
                price: service.price,
                name: service.name,
                category_id: service.category_id
            };
        },

        openMaterialModal(material) {
            this.selectedMaterial = {
                id: material.id,
                price: material.price,
                name: material.name
            };
        },

        closeServiceModal() {
            this.selectedService = null;
            this.isEditingService = false;
        },

        closeMaterialModal() {
            this.selectedMaterial = null;
            this.isEditingMaterial = false;
        },

        startEditingService() {
            this.editedServicePrice = this.selectedService.price;
            this.isEditingService = true;
        },

        startEditingMaterial() {
            this.editedMaterialPrice = this.selectedMaterial.price;
            this.isEditingMaterial = true;
        },

        getCategoryName(categoryId) {
            const cat = this.serviceCategories.find(c => c.id === categoryId);
            return cat ? cat.name : '';
        },

        async saveServicePrice() {
            if (this.editServiceError) return;
            try {
                const res = await fetch(`/api/services/${this.selectedService.id}`, {
                    method: 'PUT',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({name: this.selectedService.name, price: this.editedServicePrice, category_id: this.selectedService.category_id })
                });

                if (!res.ok) throw new Error('Ошибка при обновлении');

                this.closeServiceModal();
                await this.fetchServices();
            } catch (err) {
                alert('Не удалось сохранить изменения.');
            }
        },

        async saveMaterialPrice() {
            if (this.editMaterialError) return;
            try {
                const res = await fetch(`/api/materials/${this.selectedMaterial.id}`, {
                    method: 'PUT',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ name: this.selectedMaterial.name, price: this.editedMaterialPrice })
                });

                if (!res.ok) throw new Error('Ошибка при обновлении');

                this.closeMaterialModal();
                await this.fetchMaterials();
            } catch (err) {
                alert('Не удалось сохранить изменения.');
            }
        },

        async confirmServiceDelete() {
            if (!confirm(`Удалить услугу "${this.selectedService.name}"?`)) return;

            try {
                const res = await fetch(`/api/services/${this.selectedService.id}`, {
                    method: 'DELETE'
                });

                if (!res.ok) throw new Error('Ошибка при удалении');

                this.closeServiceModal();
                await this.fetchServices();
            } catch (err) {
                alert('Не удалось удалить услугу.');
            }
        },

        async confirmMaterialDelete() {
            if (!confirm(`Удалить расходный материал "${this.selectedMaterial.name}"?`)) return;

            try {
                const res = await fetch(`/api/materials/${this.selectedMaterial.id}`, {
                    method: 'DELETE'
                });

                if (!res.ok) throw new Error('Ошибка при удалении');

                this.closeMaterialModal();
                await this.fetchMaterials();
            } catch (err) {
                alert('Не удалось удалить расходный материал.');
            }
        },

        openAddModal() {
            this.resetNewItem();
            this.addError = '';
            const modal = new bootstrap.Modal(document.getElementById('addModal'));
            modal.show();
        },

        resetNewItem() {
            this.newItem = { name: '', category_id: '', price: 0 };
        },

        async submitAdd() {
            if (this.addNameError || this.addCategoryError || this.addPriceError) return;

            const url = this.activeTab === 'services' ? '/api/services' : '/api/materials';
            const payload = this.activeTab === 'services'
                ? { name: this.newItem.name, price: this.newItem.price, category_id: this.newItem.category_id }
                : { name: this.newItem.name, price: this.newItem.price };

            try {
                const res = await fetch(url, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(payload)
                });

                if (!res.ok) {
                    const data = await res.json();
                    throw new Error(data.error || 'Ошибка при добавлении');
                }

                await this.fetchData();
                const modalEl = document.getElementById('addModal');
                bootstrap.Modal.getInstance(modalEl).hide();
            } catch (err) {
                this.addError = err.message;
            }
        },

        togglePopover() {
            this.isPopoverVisible = !this.isPopoverVisible;
        },

        handleClickOutside(event) {
            const popover = document.getElementById('admin-profile');
            if (popover && !popover.contains(event.target)) {
                this.isPopoverVisible = false;
            }
        }
    },

    watch: {
        activeTab: 'fetchData',
        'editedServicePrice'(val) {
            this.editServiceError = (val === '' || val === null) ? 'Цена не может быть пустой' : val < 0 ? 'Цена не может быть отрицательной' : '';
        },
        'editedMaterialPrice'(val) {
            this.editMaterialError = (val === '' || val === null) ? 'Цена не может быть пустой' : val < 0 ? 'Цена не может быть отрицательной' : '';
        },
        'newItem.name'(val) {
            this.addNameError = (!val) ? 'Название не может быть пустым' : '';
        },
        'newItem.category_id'(val) {
            if (this.activeTab === 'services') {
                this.addCategoryError = (!val) ? 'Необходимо выбрать категорию' : '';
            }
        },
        'newItem.price'(val) {
            this.addPriceError = (val === '' || val === null) ? 'Цена не может быть пустой' : val < 0 ? 'Цена не может быть отрицательной' : '';
        }
    },

    mounted() {
        this.fetchData();
        this.fetchServiceCategories();
        document.addEventListener('click', this.handleClickOutside);
    }
}).mount('#app');
