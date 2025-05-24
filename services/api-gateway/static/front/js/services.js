document.addEventListener("DOMContentLoaded", () => {
  fetch("http://192.168.1.207:8080/api/services")
    .then((res) => res.json())
    .then((data) => {
      const container = document.getElementById("servicesContainer");
      if (!container) {
        console.error("Элемент #servicesContainer не найден");
        return;
      }

      data.categories.forEach((category) => {
        // Заголовок категории
        const catHeading = document.createElement("h4");
        catHeading.className = "mt-5 mb-3";
        catHeading.textContent = category.name;
        container.appendChild(catHeading);

        // Список услуг
        const ul = document.createElement("ul");
        ul.className = "list-unstyled";

        category.services.forEach((service) => {
          const li = document.createElement("li");
          li.className = "d-flex justify-content-between border-bottom py-2";

          li.innerHTML = `
            <span>${service.name}</span>
            <span>${service.price}₽</span>
          `;
          ul.appendChild(li);
        });

        container.appendChild(ul);
      });
    })
    .catch((err) => {
      console.error("Ошибка при загрузке услуг:", err);
    });
});
