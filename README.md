# benzingatest
Simple Web Service.

#### Endpoints
-  /healthz -> Accepts GET requests and returns the health
-  /log -> Accepts POST requests.
-  Batch size, batch interval and external endpoint are configurable.
-  Dockerfile and benzinga-compose.yaml files are present. One can run container either by Dockerfile or docker-compose.

#### Assumptions or Points to note
- This is very bare minimum application. There might be edge cases or special use-cases which are not handled. Application might give errors in that case.
- There is a lot of scope to improve in terms of modularity, unti tests, etc.
