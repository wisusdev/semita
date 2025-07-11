# Sistema de Autenticación OAuth2 (Inspirado en Laravel Passport)

Este sistema implementa autenticación OAuth2 para APIs en Go, inspirado en Laravel Passport. Permite la emisión, validación y revocación de tokens de acceso y refresh, así como la gestión de clientes y scopes.

## Tabla de Contenidos

- [Sistema de Autenticación OAuth2 (Inspirado en Laravel Passport)](#sistema-de-autenticación-oauth2-inspirado-en-laravel-passport)
  - [Tabla de Contenidos](#tabla-de-contenidos)
  - [Introducción](#introducción)
  - [Flujo de Autenticación](#flujo-de-autenticación)
  - [Estructura de Tablas](#estructura-de-tablas)
    - [oauth\_clients](#oauth_clients)
    - [oauth\_tokens](#oauth_tokens)
    - [oauth\_scopes](#oauth_scopes)
  - [Endpoints Principales](#endpoints-principales)
    - [Registro](#registro)
    - [Login](#login)
    - [Logout](#logout)
    - [Refresh Token](#refresh-token)
  - [Ejemplo de Uso de Token](#ejemplo-de-uso-de-token)
  - [Scopes](#scopes)
  - [Revocación de Tokens](#revocación-de-tokens)
  - [Notas de Seguridad](#notas-de-seguridad)
  - [Referencias](#referencias)

---

## Introducción

El sistema utiliza JWT como formato de token y almacena los tokens en la base de datos para permitir su revocación y control. Los clientes OAuth pueden ser gestionados desde la base de datos y cada token puede tener uno o varios scopes.

---

## Flujo de Autenticación

1. **Registro/Login de Usuario:**  
   El usuario se registra o inicia sesión y recibe un `access_token` y un `refresh_token`.

2. **Acceso a Recursos Protegidos:**  
   El cliente usa el `access_token` en la cabecera `Authorization: Bearer {token}` para acceder a rutas protegidas.

3. **Renovación de Token:**  
   Cuando el `access_token` expira, el cliente puede usar el `refresh_token` para obtener un nuevo token.

4. **Revocación de Token:**  
   El usuario puede cerrar sesión, lo que revoca el token en la base de datos.

---

## Estructura de Tablas

### oauth_clients

| Campo         | Tipo         | Descripción                |
|---------------|--------------|----------------------------|
| id            | INT          | Identificador único        |
| name          | VARCHAR(255) | Nombre del cliente         |
| client_id     | VARCHAR(100) | ID público del cliente     |
| client_secret | VARCHAR(255) | Secreto del cliente        |
| redirect_uri  | VARCHAR(255) | URI de redirección         |
| grant_types   | VARCHAR(255) | Tipos de grant permitidos  |
| scopes        | VARCHAR(255) | Scopes permitidos          |

### oauth_tokens

| Campo         | Tipo         | Descripción                |
|---------------|--------------|----------------------------|
| id            | INT          | Identificador único        |
| user_id       | INT          | ID del usuario             |
| client_id     | INT          | ID del cliente             |
| access_token  | VARCHAR(512) | Token de acceso JWT        |
| refresh_token | VARCHAR(512) | Token de refresco JWT      |
| scopes        | VARCHAR(255) | Scopes asignados           |
| revoked       | TINYINT(1)   | Si el token está revocado  |
| expires_at    | DATETIME     | Fecha de expiración        |

### oauth_scopes

| Campo         | Tipo         | Descripción                |
|---------------|--------------|----------------------------|
| id            | INT          | Identificador único        |
| name          | VARCHAR(100) | Nombre del scope           |
| description   | VARCHAR(255) | Descripción del scope      |

---

## Endpoints Principales

### Registro

```
POST /api/auth/register
{
  "name": "Juan",
  "email": "juan@ejemplo.com",
  "password": "0123456789",
  "confirm_password": "0123456789"
}
```

**Respuesta:**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 86400,
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "scope": ""
}
```

### Login

```
POST /api/auth/login
{
  "email": "juan@ejemplo.com",
  "password": "secreto"
}
```

**Respuesta:** igual que el registro.

### Logout

```
POST /api/auth/logout
Authorization: Bearer {access_token}
```

**Respuesta:**
```json
{
  "message": "Sesión cerrada correctamente"
}
```

### Refresh Token

```
POST /api/auth/refresh-token
Authorization: Bearer {refresh_token}
```

**Respuesta:** igual que login.

---

## Ejemplo de Uso de Token

Para acceder a rutas protegidas:

```
GET /api/user/profile
Authorization: Bearer {access_token}
```

---

## Scopes

Los scopes permiten limitar los permisos de los tokens.  
Ejemplo de scopes: `read`, `write`, `admin`.

- Al crear un token, puedes especificar los scopes permitidos.
- El middleware verifica que el token tenga los scopes requeridos para cada endpoint.

**Ejemplo de uso en middleware:**

```go
router.GET("/admin", middlewares.AuthMiddleware(), middlewares.ScopeMiddleware("admin"), adminHandler)
```

---

## Revocación de Tokens

- Al hacer logout, el token se marca como revocado en la base de datos.
- Los tokens revocados no pueden ser usados para acceder a recursos protegidos.

---

## Notas de Seguridad

- Los tokens JWT se validan y además se verifica su existencia y estado en la base de datos.
- Los refresh tokens también son JWT y se almacenan para permitir su revocación.
- Los secretos de los clientes deben mantenerse privados.

---

## Referencias

- [OAuth 2.0 RFC](https://datatracker.ietf.org/doc/html/rfc6749)
- [Laravel Passport](https://laravel.com/docs/10.x/passport)
