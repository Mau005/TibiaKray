function sendForm() {
    // Obtén los valores del formulario
    var username = $('#user').val();
    var password = $('#passworduser').val();
    var messageChange = document.getElementById(error-message);

    // Construye los datos del formulario en formato JSON
    var formData = {
        username: username,
        password: password
    };

    // Realiza la solicitud utilizando jQuery.ajax
    $.ajax({
        url: '/login',
        type: 'POST',
        contentType: 'application/json',
        data: JSON.stringify(formData),
        success: function(data) {
            console.log('Respuesta del servidor:', data);
            messageChange.innerHTML = "El RUT ingresado existe: ";
            messageChange.style.color = "red";

            // Puedes realizar más acciones basadas en la respuesta aquí
        },
        error: function(error) {
            console.error('Error en la solicitud:', error);
        }
    });
}