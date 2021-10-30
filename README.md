# benzingatest
Simple Web Service.

#### Features
-  Endpoint /healthz -> Accepts GET requests and returns the health
-  Endpoint /log -> Accepts POST requests.
-  Batch size, batch interval and external endpoint are configurable.
-  Dockerfile and benzinga-compose.yaml files are present. One can run container either by Dockerfile or docker-compose.
-  Uses in-memory cache to store payload temporarily. POST as per batch size and interval. Retries 3 times in case of errors. Application exits couldn't POST to external endpoint after 3 retries.


#### Run
- With Docker-compose -> ```docker-compose -f benzinga-compose.yaml up```
- With Dockerfile -> ```docker build -t benzinga:1.0 .; docker run -p 9000:9000 --name benzinga benzinga:1.0```
  


#### Assumptions or Points to note
- This is very bare minimum application. There might be edge cases or special use-cases which are not handled. Application might give errors in that case.
- There is a lot of scope to improve in terms of modularity, unti tests, etc.
