```
ab -kc 10 -t 60 http://localhost:8080/
Requests per second:    1231.18 [#/sec] (mean)

ab -kc 10 -t 60 http://localhost:8080/blog/
Requests per second:    904.83 [#/sec] (mean)

ab -kc 10 -t 60 http://localhost:8080/blog/category/health-and-food/
Requests per second:    587.40 [#/sec] (mean)

ab -kc 10 -t 60 http://localhost:8080/blog/the-best-juices-for-a-child/
Requests per second:    1197.49 [#/sec] (mean)

ab -kc 10 -t 60 http://localhost:8080/shop/
Requests per second:    830.99 [#/sec] (mean)

ab -kc 10 -t 60 http://localhost:8080/shop/category/electronics/
Requests per second:    589.96 [#/sec] (mean)

ab -kc 10 -t 60 http://localhost:8080/shop/samsung-galaxy-s10-128-gb/
Requests per second:    624.73 [#/sec] (mean)
```
