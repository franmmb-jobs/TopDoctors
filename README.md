
## Índice
1. [Resumen resolución](#resumen-resolución)
2. [Justificaciones](#justificaciones)
3. [Cómo ejecutar el proyecto](#cómo-ejecutar-el-proyecto)

## Resumen resolución



## Justificaciones
se ha decidido usar la libreria de go http en vez de un framwork u otra del estilo para demostrar el conocimiento de la libreria de go y no depender de frameworks que no son estandar de go

configuracion tanto .yml como .env, se ha usado .yml para configuración no sensible y .env para datos sensibles. con apoyo de la libreria de terceros cleanenv, ya que no considero que sea necesario implementar una configuración más compleja y pese a ser algo critico se puede delegar a la libreria y si falla es facilmente solucionable.

Logger
Se ha decidido usar el logger estandar de go slog, con la libreria de terceros tint para darle un formato mas legible y con colores, se ha decidido usar este logger por que es el estandar de go y es muy facil de usar, ademas de que es muy rapido y eficiente. (Aparte que queria probar este metodo de configuracion de logger)
Cambia el log por defecto de go por un logger personalizado que permite personalizar el formato, el nivel y el output, tambien altera el funcionamiento de log pero no se utilizará


## Cómo ejecutar el proyecto


go run .\cmd\api\main.go -config='configs/config.dev.yml'