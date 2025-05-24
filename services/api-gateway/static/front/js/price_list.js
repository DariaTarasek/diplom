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
    }
  },

  methods: {
    async fetchData() {
      try {
       
        const userRes = await fetch('http://192.168.1.207:8080/api/admin-data');
        const userData = await userRes.json();
        this.first_name = userData.first_name || '';
        this.second_name = userData.second_name || '';
        this.userRole = userData.role.role || '';

        await this.fetchServices();
        await this.fetchMaterials();
      } catch (err) {
        console.error('Ошибка при загрузке данных:', err);
      }
    },

    async fetchServices() {
      try {
        const params = new URLSearchParams();
        if (this.searchQuery) params.append('search', this.searchQuery);
        if (this.selectedCategory) params.append('category', this.selectedCategory);

        const res = await fetch(`http://192.168.1.207:8080/api/services?${params}`);
        const data = await res.json();

        this.services = data.services || [];
      } catch (err) {
        console.error('Ошибка при загрузке услуг:', err);
      }
    },

    async fetchMaterials() {
      try {
        const params = new URLSearchParams();
        if (this.searchQuery) params.append('search', this.searchQuery);

        const res = await fetch(`http://192.168.1.207:8080/api/materials?${params}`);
        const data = await res.json();

        this.materials = data.materials || [];
      } catch (err) {
        console.error('Ошибка при загрузке материалов:', err);
      }
    },

    async fetchServiceCategories() {
        try {
            const res = await fetch('http://192.168.1.207:8080/api/service-categories');
            const data = await res.json();
            this.serviceCategories = data.categories || [];
        } catch (err) {
            console.error('Ошибка при загрузке категорий:', err);
        }
        },

    async onSearchOrFilterChange() {
      if (this.activeTab === 'services') {
        await this.fetchServices();
      } else if (this.activeTab === 'materials') {
        await this.fetchMaterials();
      }
    },

     openServiceModal(service) {
        this.selectedService = { ...service };
  },

    openMaterialModal(material) {
         this.selectedMaterial = { ...material };
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

  async saveServicePrice() {
    if (this.editServiceError) {
        return;
    }
    try {
      const res = await fetch(`http://192.168.1.207:8080/api/services/${this.selectedService.id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ price: this.editedServicePrice })
      });

      if (!res.ok) throw new Error('Ошибка при обновлении');

      this.isEditingService = false;
      this.closeServiceModal();
      await this.fetchServices();
    } catch (err) {
      alert('Не удалось сохранить изменения.');
    }
  },

   async saveMaterialPrice() {
    if (this.editMaterialError) {
        return;
    }
    try {
      const res = await fetch(`http://192.168.1.207:8080/api/materials/${this.selectedMaterial.id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ price: this.editedMaterialPrice })
      });

      if (!res.ok) throw new Error('Ошибка при обновлении');

      this.isEditingMaterial = false;
      this.closeMaterialModal();
      await this.fetchMaterials();
    } catch (err) {
      alert('Не удалось сохранить изменения.');
    }
  },

  async confirmServiceDelete() {
    if (!confirm(`Удалить услугу "${this.selectedService.name}"?`)) return;

    try {
      const res = await fetch(`http://192.168.1.207:8080/api/services/${this.selectedService.id}`, {
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
      const res = await fetch(`http://192.168.1.207:8080/api/materials/${this.selectedMaterial.id}`, {
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
    this.newItem = {
      name: '',
      category_id: '',
      price: 0
    };
  },

  async submitAdd() {
    if (this.addNameError || this.addCategoryError || this.addPriceError) {
        console.log('otleteli na proverke' + this.addNameError + '\n' + this.addCategoryError + '\n' + this.addPriceError);
        return;
    }
    const url =
      this.activeTab === 'services'
        ? 'http://192.168.1.207:8080/api/services'
        : 'http://192.168.1.207:8080/api/materials';

    const payload =
      this.activeTab === 'services'
        ? {
            name: this.newItem.name,
            price: this.newItem.price,
            category_id: this.newItem.category_id
          }
        : {
            name: this.newItem.name,
            price: this.newItem.price
          };

    try {
      const res = await fetch(url, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(payload)
      });

      if (!res.ok) {
        const data = await res.json();
        throw new Error(data.error || 'Ошибка при добавлении');
      }

      await this.onSearchOrFilterChange();

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
        searchQuery: 'fetchData',
        selectedCategory: 'fetchData',
        'editedServicePrice'(newVal) {
            if (newVal === '' || newVal === null) {
                this.editServiceError = 'Цена не может быть пустой';
            } else if (newVal < 0) {
                this.editServiceError = 'Цена не может быть отрицательной';
            } else {
                this.editServiceError = '';
            }
        },
        'editedMaterialPrice'(newVal) {
            if (newVal === '' || newVal === null) {
                this.editMaterialError = 'Цена не может быть пустой';
            } else if (newVal < 0) {
                this.editMaterialError = 'Цена не может быть отрицательной';
            } else {
                this.editMaterialError = '';
            }
        },
        'newItem.name'(newVal) {
            if (newVal === '' || newVal === null) {
                this.addNameError = 'Название не может быть пустым';
            } else {
                this.addNameError = '';
            }
        },
        'newItem.category_id'(newVal) {
            if (this.activeTab === 'services') {
            if (newVal === '' || newVal === null) {
                this.addCategoryError = 'Необходимо выбрать категорию';
            } else {
                this.addCategoryError = '';
            }
            }
        },
        'newItem.price'(newVal) {
            if (newVal === '' || newVal === null) {
                this.addPriceError = 'Цена не может быть пустой';
            } else if (newVal < 0) {
                this.addPriceError = 'Цена не может быть отрицательной';
            } else {
                this.addPriceError = '';
            }
        },
    },



  mounted() {
    this.fetchData().finally(() => {
  });
    this.fetchServiceCategories();
    document.addEventListener('click', this.handleClickOutside);
  }
}).mount('#app');
