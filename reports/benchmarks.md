## Benchmarks:
Result without cache:
```
  BenchmarkCampaignHandler-16          	     358	   3293453 ns/op	         2.352 response_time_ms
```

Result with cache:
```
  BenchmarkCachedCampaignHandler-16    	 1000000	      1037 ns/op	         0 response_time_ms
```
