# Adtelligent Internship

> **Performed by:** Eugeniu Popa

## Run the application

Before running the application, make sure [Docker](https://www.docker.com/) is installed.  
Type these commands in the root folder.

```bash
docker compose up
```

```bash
go run main.go  
```

## Endpoints:

- Get data:

  ```
  GET "http://localhost:8080/campaigns?source_id=1"
  ```

- Get cached data:

  ```
  GET "http://localhost:8080/cached-campaigns?source_id=1"
  ```

## Benchmarks:
Result without cache:
```
  BenchmarkCampaignHandler-16          	     358	   3293453 ns/op	         2.352 response_time_ms
```

Result with cache:
```
  BenchmarkCachedCampaignHandler-16    	 1000000	      1037 ns/op	         0 response_time_ms
```
