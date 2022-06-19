# TP TDL Unbutú
----

Para correr el tp: `go run ./tp-tdl-unbutu`.

Desde el browser ir a `http://localhost:8080/ui`.

Rutas:
- `/date`: Crea un nuevo trabajo (obtiene la fecha actual después de 10 seg). Devuelve el Id de trabajo.
- `/date/:jobId/status`: Devuelve el estado del trabajo `jobId`
- `/date/:jobId/output`: Devuelve la salida del trabajo `jobId`
