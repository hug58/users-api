#!/bin/bash


#este es un cron job para eliminar tokens vencidos y cambiar el estaus del usuario a false
echo "* * * * * psql -c 'SELECT delete_expired_tokens();' -d postgres -U postgres" | crontab -
# Muestra un mensaje de confirmación
echo "El cronjob ha sido creado y programado correctamente."
