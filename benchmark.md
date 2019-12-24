```
ab -kc 10 -t 60 http://localhost:8080/
Requests per second:    1355.88 [#/sec] (mean)

ab -kc 10 -t 60 http://localhost:8080/blog/
Requests per second:    904.27 [#/sec] (mean)

ab -kc 10 -t 60 http://localhost:8080/blog/category/health-and-food/
Requests per second:    603.94 [#/sec] (mean)

ab -kc 10 -t 60 http://localhost:8080/blog/the-best-juices-for-a-child/
Requests per second:    1186.13 [#/sec] (mean)

ab -kc 10 -t 60 http://localhost:8080/shop/
Requests per second:    824.24 [#/sec] (mean)

ab -kc 10 -t 60 http://localhost:8080/shop/category/electronics/
Requests per second:    596.34 [#/sec] (mean)

ab -kc 10 -t 60 http://localhost:8080/shop/samsung-galaxy-s10-128-gb/
Requests per second:    620.03 [#/sec] (mean)
```
