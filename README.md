# Challenge

## Estructura
Se debe crear una carpeta local para procesar las transacciones, dentro de esta carpeta se debe crear la siguiente estructura:

- /pending
- /processed

Por ejemplo si la carpeta local se llama /transactions, la estructura deberia quedar de la siguiente forma:

- /transactions/pending
- /transactions/processed

## Ejecución

### Ejecución con docker

Edite el archivo docker-compose.yaml y modifique las siguientes variables:

EMAIL_TO: correo a donde se van enviar las estadisticas.
EMAIL_USERNAME: cuenta de correo de gmail desde donde se va a enviar el correo.
EMAIL_PASSWORD: clave de gmail donde se va a enviar el correo. Esta clave es una clave de tipo app especial que se puede generar desde gmail para poder enviar correos por smtp.

Y en el volume modifique /folder_local_transactions por la carpeta local donde se van a procesar los archivos. Esta carpeta debe tener la estructura creada.

### Ejecución con lambda

Para crear el paquete .zip para subir a lambda ejecutar:

GOOS=linux GOARCH=amd64 go build -o bootstrap cmd/challenge/main.go
zip boostrap.zip template template/*  bootstrap

Luego de esto en la consola de lambda se puede subir el .zip que contiene el binario.

En la lamba configure las siguientes variables de entorno:

RUNTIME: lambda
EMAIL_TO: correo a donde se van enviar las estadisticas.
EMAIL_USERNAME: cuenta de correo de gmail desde donde se va a enviar el correo.
EMAIL_PASSWORD: clave de gmail donde se va a enviar el correo. Esta clave es una clave de tipo app especial que se puede generar desde gmail para poder enviar correos por smtp.