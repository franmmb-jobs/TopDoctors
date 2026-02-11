# TopDoctors backend code challenge

## Índice
1. [El problema propuesto](#el-problema-propuesto)
2. [Nuestras expectativas](#nuestras-expectativas)
3. [La entrega](#la-entrega)
4. [Hazlo lo mejor que sabes](#hazlo-lo-mejor-que-sabes)

## El problema propuesto

Uno de nuestros stakeholders nos solicita que desarrollemos una solución para
consultar y almacenar datos de diagnóstico de pacientes en nuestros sistemas,
esta solución aparte de funcionar internamente, necesita tener la capacidad
para integrarse con otras aplicaciones. Ambas partes atesoran datos sensibles
por lo que será necesario implementar una forma de autenticar las peticiones.

Para lograr este objetivo, haremos lo siguiente:

- Un endpoint que requiera de un usuario y contraseña para generar un token de autenticación
- Un endpoint protegido que permita consultar diagnósticos y filtrarlos por: nombre del paciente y/o fecha
- Un endpoint protegido que permita almacenar diagnósticos para un paciente en concreto

Sobre la estructura de la base de datos, necesitamos que existan pacientes,
sobre los cuales guardaremos los siguientes datos:

- nombre
- dni
- email
- teléfono (opcional)
- dirección (opcional)

Estos pacientes tendrán una relación con sus datos de diagnóstico, sobre los
cuáles guardaremos los siguientes datos:

- paciente
- diagnóstico
- prescripción (opcional)
- fecha

## Nuestras expectativas

Aunque normalmente trabajamos con Golang, no tenemos problema si nos entregas la
solución en Javascript, Golang o Python. También siéntente libre de incorporar mejoras
al problema propuesto, explícanos que has detectado y que camino has decidido tomar.

Tu código debería ser capaz de hacer sonar las alarmas cuando hay un cambio que
rompe alguna pieza, para eso no hay mejor remedio que realizar unos buenos tests.
Se requiere tests en al menos la parte de la creaación/almacenamiento de diagnósticos.

Para que nos sea fácil comprobar que tu solución funciona, documenta
los pasos que tenemos que seguir para ponerla en marcha en el README.md del proyecto.
Se valorará el uso de contenedores para facilitar la ejecución de la solución.

Y por último no dudes en incluir todo aquello que consideres oportuno, o creas que
aporta valor extra. Por ejemplo, un Swagger integrado en el proyecto para ver y probar
los endpoints de la API, o un Postman que contenga los endpoints ya preparados.

# La entrega

Y una vez terminado, ¿cómo nos lo entregas?

Sencillo, puedes crear tu repositorio privado en github (u otro) y danos acceso, así también podemos
ver como trabajas con git :).
También puedes enviarnos la prueba por correo eletrónico si lo prefieres.

# Hazlo lo mejor que sabes

Sabemos por experiencia que estas pruebas no son triviales, por eso dispones de un
margen de tiempo bastante holgado para realizarla, no te apresures, disfruta de la
experiencia, expláyate y si puedes trata de aprender algo nuevo en el camino.

¡Un saludo!
El equipo de backend
