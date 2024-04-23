# Adtelligent Internship

> **Performed by:** Eugeniu Popa

## Run the application

Before running the application, make sure [Docker](https://www.docker.com/) is installed.  
Type these commands in the root folder.

```bash
docker compose up
```

```bash
go run .  
```

## Endpoints:

- Get data:

  ```
  GET "http://localhost:8080/campaigns_map?source_id=1&domain=gmail.com"
  ```

  ```
  GET "http://localhost:8080/campaigns_slice?source_id=1&domain=gmail.com"
  ```

## [Benchmarks](reports/Benchmarks.md)

## [Profiling](reports/Profiling.md)
