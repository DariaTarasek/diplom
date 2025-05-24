document.addEventListener("DOMContentLoaded", () => {
  fetch("http://192.168.1.207:8080/api/doctors")
    .then(res => res.json())
    .then(doctors => {
      const container = document.getElementById("doctorsContainer");
      const template = document.getElementById("doctorTemplate");

      container.innerHTML = ""; // очищаем контейнер перед добавлением

      doctors.forEach(doctor => {
        const clone = template.cloneNode(true);
        clone.style.display = "block";
        clone.id = ""; 

      
        const title = clone.querySelector(".card-title");
        const text = clone.querySelector(".card-text");
        const img = clone.querySelector(".card-img-top");

        if (title) title.textContent =  doctor.second_name + " " + doctor.first_name + " " + doctor.surname;
        if (text) text.textContent = "Специализация: " + doctor.specialty.name;
        if (img) img.setAttribute("src", doctor.photo_url);

        container.appendChild(clone);
      });
    })
    .catch(err => {
      console.error("Ошибка:", err);
    });
});
