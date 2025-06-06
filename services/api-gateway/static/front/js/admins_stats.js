const { createApp } = Vue;

createApp({
    data() {
        return {
            stats: {
                totalPatients: 0,
                totalVisits: 0,
                topServices: [],
                doctorAvgVisit: [],
                doctorCheckStat: [],
                doctorUniquePatients: [],
                newPatientsThisMonth: 0,
                avgVisitPerPatient: 0,
                ageGroupStat: [],
                totalIncome: 0,
                monthlyIncome: 0,
                clinicAvgCheck: 0
            },
            fullName: '',
            isPopoverVisible: false
        };
    },

    methods: {
        async fetchStats() {
            try {
                const res = await fetch('/api/statistics');
                const data = await res.json();

                this.stats = {
                    ...this.stats,
                    ...data
                };
            } catch (err) {
                console.error('Ошибка при загрузке статистики:', err);
            }
        },

        async fetchAdminData() {
            try {
                const res = await fetch('/api/admin/me');
                const data = await res.json();

                const firstName = data.firstName || '';
                const secondName = data.secondName || '';
                this.fullName = `${firstName} ${secondName}`.trim();
            } catch (err) {
                console.error('Ошибка при загрузке данных администратора:', err);
            }
        }
    },

    mounted() {
        this.fetchStats();
        this.fetchAdminData();
    }
}).mount('#app');
