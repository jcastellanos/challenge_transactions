# Challenge

Demo: https://youtu.be/55sDjkcUlOg

## Arquitectura

La arquitectura del proyecto esta basada en una arquitectura limpia donde se tiene desarrollada la lógica de negocio en un caso de uso y se tienen adaptadores para las diferentes infraestructuras, se soporta infraestructura "standalone" que funciona como una aplicación binaria o dockerizada y una infraestructura en lambda.

La infraestructura standalone guarda los registros en una base de datos SQLite embebida y la infraestructura lambda guarda los registros en una base de datos dynamodb. Ambos envian correo a través del SMTP de gmail.

## Infraestructura standalone

Para ejecutar el proyecto como docker se debe editar el archivo docker-compose.yaml y modificar las siguientes variables:

EMAIL_TO: correo a donde se van enviar las estadisticas.
EMAIL_USERNAME: cuenta de correo de gmail desde donde se va a enviar el correo.
EMAIL_PASSWORD: clave de gmail donde se va a enviar el correo. Esta clave es una clave de tipo app especial que se puede generar desde gmail para poder enviar correos por smtp.

Y en el volume dentro del docker-compose se debe modificar el /folder_local_transactions por la carpeta local donde se van a procesar los archivos. Esta carpeta debe tener la siguiente estructura creada:

- /folder_local_transactions/pending: carpeta donde se copian los .csv con las transacciones para procesarlos.
- /folder_local_transactions/processed: carpeta donde se mueven los archivos .csv luego de procesarlos.

Para ejecutar el docker se puede utilizar el comando: 

```console
> docker-compuse up
```

### Infraestructura con lambda

Se debe crear paquete .zip para subir al servicio de lambda para esto se puede ejecutar:

```console
> GOOS=linux GOARCH=amd64 go build -o bootstrap cmd/challenge/main.go
> zip boostrap.zip template template/*  bootstrap
```

Luego de esto en la consola de lambda se puede subir el archivo bootstrap.zip.

En la lamba se deben configurar las siguientes variables de entorno:

RUNTIME: lambda
EMAIL_TO: correo a donde se van enviar las estadisticas.
EMAIL_USERNAME: cuenta de correo de gmail desde donde se va a enviar el correo.
EMAIL_PASSWORD: clave de gmail donde se va a enviar el correo. Esta clave es una clave de tipo app especial que se puede generar desde gmail para poder enviar correos por smtp.

La lambda se debe configurar para que el desencadenador sea la escritura de un archivo .csv en un bucket.

La lambda debe ser ejecutada con un rol que tenga la siguiente politica asociada (esta politica es necesario cerrarla mejor para mejorar la seguridad):

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "logs:PutLogEvents",
                "logs:CreateLogGroup",
                "logs:CreateLogStream"
            ],
            "Resource": "arn:aws:logs:*:*:*"
        },
        {
            "Effect": "Allow",
            "Action": [
                "s3:GetObject"
            ],
            "Resource": "arn:aws:s3:::*/*"
        },
        {
            "Effect": "Allow",
            "Action": [
                "dynamodb:PutItem",
                "dynamodb:BatchWriteItem"
            ],
            "Resource": "arn:aws:dynamodb:*:*:*"
        }
    ]
}
```