### enviroments variable

- POSTGRES_HOST
- POSTGRES_PORT
- POSTGRES_DB
- POSTGRES_USER
- POSTGRES_PASSWORD
- GTFS_QUERY_API_KEY
- GOOGLE_GEOCODE_API_KEY


### health

`GET /health` returns `{"status":"ok"}` without touching the database. The
Kubernetes readiness and liveness probes in `k8s/deployment.yaml` hit it on
port 8080.
