# SEMITA

Semita es un pequeño marco de trabajo inspirado en Laravel

## Guia de instalación bascia

### Copiar archivo .nev

```bash
cp .env.example .env
```

Establece tus credenciales en el archivo de configuración

### Generar un llave unica para semita

```bash
go run . key:generate
```

### Realizamos la migración de nuestras tablas

```bash
go run . migrate
```
