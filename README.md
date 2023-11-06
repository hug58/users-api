
# USERS API


Requisitos
Go 1.19 o superior instalado en tu sistema.

## Instalación
Clona este repositorio en tu máquina local:

```
git clone https://github.com/hug58/users-api.git
```
Navega hasta el directorio del proyecto:

```
cd users-api
```

Instala las dependencias del proyecto:

```sh
go mod tidy 
```

crea un archivo .env y configura la db
```sh
export PGHOST="localhost"
export PGPORT=5432
export PGDATABASE="postgres"
export PGUSER="postgres"
export PGPASSWORD="postgres"
export SECRET="SUPERSECRET!"
export TOKEN_EXPIRATION_MINUTES="5m"
```
carga las variables de entorno
```sh

source .env
```

## Inicia el servidor

```
go run main.go
```


Una vez que el servidor esté en funcionamiento, puedes acceder a la API REST a través de http://localhost:8080. A continuación, se muestran algunos ejemplos de rutas disponibles:

* <span style="color: yellow">POST</span> api/v1/users/login 
```json
{
    "phone":"+580000000",
    "password":"Gugosss"
}
```
* <span style="color: green">GET</span> api/v1/users: Obtiene todos los usuarios.
* <span style="color: green">GET</span> api/v1/users/{id}: Obtiene un usuario específico por su ID.

* <span style="color: yellow">POST</span> api/v1/users: Crea un nuevo usuario:
```json
{
    "name":"test",
    "addres":"casa de hugo",
    "email":"dosss@gmail.com",
    "phone":"+580000000",
    "password":"Gugosss"
}
```

* <span style="color: blue">PUT</span> api/v1/users/{id}: Actualiza un usuario existente por su ID.
```json
{
    "name":"test",
    "addres":"casa de hugo",
    "email":"dosss@gmail.com",
    "phone":"+580000000",
    "password":"Gugosss"
}
```


* <span style="color: red">DELETE</span> api/v1/users/{id}: Elimina un usuario por su ID.