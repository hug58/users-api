#!/bin/bash


# Añade la línea al final del archivo crontab, utilizando las variables de entorno cargadas
echo "* * * * * /usr/bin/docker exec -it postgres-db psql -c 'SELECT delete_expired_tokens();' -d postgres -U postgres" | crontab -

# Muestra un mensaje de confirmación
echo "El cronjob ha sido creado y programado correctamente."
