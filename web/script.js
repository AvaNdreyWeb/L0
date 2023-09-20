const form = document.querySelector('.get-by-id-form');
const input = document.querySelector('.id');
const output = document.querySelector('.order');

form.addEventListener('submit', function (e) {
    e.preventDefault();

    const ID = input.value;
    fetch(`http://localhost:8080/?id=${encodeURIComponent(ID)}`)
        .then(response => response.json())
        .then(jsonResult => {
            output.value = JSON.stringify(jsonResult, null, 2);
        })
        .catch(error => {
            console.error('Error:', error);
        });
});
