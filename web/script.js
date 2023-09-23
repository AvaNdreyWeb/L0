const form = document.querySelector('.get-by-id-form');
const input = document.querySelector('.id');
const output = document.querySelector('.order');

const serverError = {code: 500, message: "Сервис временно недоступен"}

form.addEventListener('submit', function (e) {
    e.preventDefault();

    const ID = input.value.trim();
    if (ID != '') {
        fetch(`http://localhost:8080/?id=${encodeURIComponent(ID)}`)
        .then(response => response.json())
        .then(jsonResult => {
            output.value = JSON.stringify(jsonResult, null, 2);
        })
        .catch(error => {
            output.value = JSON.stringify(serverError, null, 2);
        });
    }
    
});
