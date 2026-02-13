
## Índice
1. [Resumen resolución](#resumen-resolución)
2. [Justificaciones](#justificaciones)
3. [Cómo ejecutar el proyecto](#cómo-ejecutar-el-proyecto)

## Resumen resolución

genralmente a la hora de construir la estructura de la informacion le daria importancia a la tabla de pacientes que entiendo que podrian tener varios diagnosticos asociados por como se podria entender que funcionarioa la informacion recolectada. Por esta vez se ha optado en declarar que la informacion viene dada como prioridad por los diagnosticos en vez de por los pacientes y se ha estructurado el codigo de esta manera. 

aunque seguramente vayan a hacer falta las opciones como modificar o eliminar registros para diagnosticos y posiblemente ususarios no se implemnete por falta de tiempo.


para simplicidad de la solucion se ha implementado una desision basada en monolito por lo que solo se basa en un solo repositorio y no se ha implementado una arquitectura de microservicios.

no se ha implementado cache para la base de datos por falta de tiempo y por que no se considero necesario para el funcionamiento del proyecto.

se ejecuta automigrate cada vez que se inicia el programa para no complicar y que compruebe la base de datos si esta generada correctamente




## Documentación API (Swagger)

La API cuenta con documentación automática generada con **swaggo**.

### Generar documentación
Si realizas cambios en las anotaciones de los handlers, regenera la documentación con:
```bash
swag init -g cmd/api/main.go
```

### Acceso a la UI
Una vez iniciado el servidor, puedes acceder a la interfaz de Swagger en:
[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

## Justificaciones
se ha decidido usar la libreria de go http en vez de un framwork u otra del estilo para demostrar el conocimiento de la libreria de go y no depender de frameworks que no son estandar de go

configuracion tanto .yml como .env, se ha usado .yml para configuración no sensible y .env para datos sensibles. con apoyo de la libreria de terceros cleanenv, ya que no considero que sea necesario implementar una configuración más compleja y pese a ser algo critico se puede delegar a la libreria y si falla es facilmente solucionable.

Logger
Se ha decidido usar el logger estandar de go slog, con la libreria de terceros tint para darle un formato mas legible y con colores, se ha decidido usar este logger por que es el estandar de go y es muy facil de usar, ademas de que es muy rapido y eficiente. (Aparte que queria probar este metodo de configuracion de logger)
Cambia el log por defecto de go por un logger personalizado que permite personalizar el formato, el nivel y el output, tambien altera el funcionamiento de log pero no se utilizará


## Cómo ejecutar el proyecto


go run .\cmd\api\main.go -config='configs/config.dev.yml'