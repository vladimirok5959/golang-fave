## Short
```
ab -kc 10 -t 60 http://localhost:8080/
Requests per second:    1531.99 [#/sec] (mean)

ab -kc 10 -t 60 http://localhost:8080/blog/
Requests per second:    845.04 [#/sec] (mean)

ab -kc 10 -t 60 http://localhost:8080/blog/category/health-and-food/
Requests per second:    417.89 [#/sec] (mean)

ab -kc 10 -t 60 http://localhost:8080/blog/the-best-juices-for-a-child/
Requests per second:    1496.85 [#/sec] (mean)

ab -kc 10 -t 60 http://localhost:8080/shop/
Requests per second:    840.95 [#/sec] (mean)

ab -kc 10 -t 60 http://localhost:8080/shop/category/electronics/
Requests per second:    228.13 [#/sec] (mean)

ab -kc 10 -t 60 http://localhost:8080/shop/samsung-galaxy-s10/
Requests per second:    744.29 [#/sec] (mean)
```

## Full
```
ab -kc 10 -t 60 http://localhost:8080/
This is ApacheBench, Version 2.3 <$Revision: 1807734 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking localhost (be patient)
Completed 5000 requests
Completed 10000 requests
Completed 15000 requests
Completed 20000 requests
Completed 25000 requests
Completed 30000 requests
Completed 35000 requests
Completed 40000 requests
Completed 45000 requests
Completed 50000 requests
Finished 50000 requests


Server Software:        fave.pro/1.3.3
Server Hostname:        localhost
Server Port:            8080

Document Path:          /
Document Length:        7594 bytes

Concurrency Level:      10
Time taken for tests:   32.637 seconds
Complete requests:      50000
Failed requests:        0
Keep-Alive requests:    0
Total transferred:      391650000 bytes
HTML transferred:       379700000 bytes
Requests per second:    1531.99 [#/sec] (mean)
Time per request:       6.527 [ms] (mean)
Time per request:       0.653 [ms] (mean, across all concurrent requests)
Transfer rate:          11718.84 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.1      0       4
Processing:     2    6   2.4      6      36
Waiting:        1    6   2.3      6      36
Total:          2    6   2.4      6      37

Percentage of the requests served within a certain time (ms)
  50%      6
  66%      7
  75%      8
  80%      8
  90%     10
  95%     11
  98%     13
  99%     14
 100%     37 (longest request)

ab -kc 10 -t 60 http://localhost:8080/blog/
This is ApacheBench, Version 2.3 <$Revision: 1807734 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking localhost (be patient)
Completed 5000 requests
Completed 10000 requests
Completed 15000 requests
Completed 20000 requests
Completed 25000 requests
Completed 30000 requests
Completed 35000 requests
Completed 40000 requests
Completed 45000 requests
Completed 50000 requests
Finished 50000 requests


Server Software:        fave.pro/1.3.3
Server Hostname:        localhost
Server Port:            8080

Document Path:          /blog/
Document Length:        8555 bytes

Concurrency Level:      10
Time taken for tests:   59.169 seconds
Complete requests:      50000
Failed requests:        0
Keep-Alive requests:    0
Total transferred:      439700000 bytes
HTML transferred:       427750000 bytes
Requests per second:    845.04 [#/sec] (mean)
Time per request:       11.834 [ms] (mean)
Time per request:       1.183 [ms] (mean, across all concurrent requests)
Transfer rate:          7257.08 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.0      0       3
Processing:     4   12   4.3     11      42
Waiting:        4   12   4.3     11      42
Total:          4   12   4.3     11      42

Percentage of the requests served within a certain time (ms)
  50%     11
  66%     13
  75%     14
  80%     15
  90%     17
  95%     20
  98%     23
  99%     26
 100%     42 (longest request)

ab -kc 10 -t 60 http://localhost:8080/blog/category/health-and-food/
This is ApacheBench, Version 2.3 <$Revision: 1807734 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking localhost (be patient)
Completed 5000 requests
Completed 10000 requests
Completed 15000 requests
Completed 20000 requests
Completed 25000 requests
Finished 25075 requests


Server Software:        fave.pro/1.3.3
Server Hostname:        localhost
Server Port:            8080

Document Path:          /blog/category/health-and-food/
Document Length:        7798 bytes

Concurrency Level:      10
Time taken for tests:   60.003 seconds
Complete requests:      25075
Failed requests:        0
Keep-Alive requests:    0
Total transferred:      201527775 bytes
HTML transferred:       195534850 bytes
Requests per second:    417.89 [#/sec] (mean)
Time per request:       23.929 [ms] (mean)
Time per request:       2.393 [ms] (mean, across all concurrent requests)
Transfer rate:          3279.90 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.0      0       0
Processing:     7   24   7.0     23      61
Waiting:        7   24   7.0     23      61
Total:          7   24   7.0     23      61

Percentage of the requests served within a certain time (ms)
  50%     23
  66%     26
  75%     28
  80%     29
  90%     33
  95%     37
  98%     41
  99%     44
 100%     61 (longest request)

ab -kc 10 -t 60 http://localhost:8080/blog/the-best-juices-for-a-child/
This is ApacheBench, Version 2.3 <$Revision: 1807734 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking localhost (be patient)
Completed 5000 requests
Completed 10000 requests
Completed 15000 requests
Completed 20000 requests
Completed 25000 requests
Completed 30000 requests
Completed 35000 requests
Completed 40000 requests
Completed 45000 requests
Completed 50000 requests
Finished 50000 requests


Server Software:        fave.pro/1.3.3
Server Hostname:        localhost
Server Port:            8080

Document Path:          /blog/the-best-juices-for-a-child/
Document Length:        8270 bytes

Concurrency Level:      10
Time taken for tests:   33.404 seconds
Complete requests:      50000
Failed requests:        0
Keep-Alive requests:    0
Total transferred:      425450000 bytes
HTML transferred:       413500000 bytes
Requests per second:    1496.85 [#/sec] (mean)
Time per request:       6.681 [ms] (mean)
Time per request:       0.668 [ms] (mean, across all concurrent requests)
Transfer rate:          12438.16 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.1      0       4
Processing:     2    7   2.1      6      36
Waiting:        2    6   2.0      6      36
Total:          2    7   2.1      6      36

Percentage of the requests served within a certain time (ms)
  50%      6
  66%      7
  75%      8
  80%      8
  90%      9
  95%     11
  98%     12
  99%     13
 100%     36 (longest request)

ab -kc 10 -t 60 http://localhost:8080/shop/
This is ApacheBench, Version 2.3 <$Revision: 1807734 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking localhost (be patient)
Completed 5000 requests
Completed 10000 requests
Completed 15000 requests
Completed 20000 requests
Completed 25000 requests
Completed 30000 requests
Completed 35000 requests
Completed 40000 requests
Completed 45000 requests
Completed 50000 requests
Finished 50000 requests


Server Software:        fave.pro/1.3.3
Server Hostname:        localhost
Server Port:            8080

Document Path:          /shop/
Document Length:        5710 bytes

Concurrency Level:      10
Time taken for tests:   59.456 seconds
Complete requests:      50000
Failed requests:        0
Keep-Alive requests:    0
Total transferred:      297450000 bytes
HTML transferred:       285500000 bytes
Requests per second:    840.95 [#/sec] (mean)
Time per request:       11.891 [ms] (mean)
Time per request:       1.189 [ms] (mean, across all concurrent requests)
Transfer rate:          4885.58 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.0      0       2
Processing:     4   12   3.8     11      87
Waiting:        4   12   3.8     11      86
Total:          4   12   3.8     11      87

Percentage of the requests served within a certain time (ms)
  50%     11
  66%     13
  75%     14
  80%     15
  90%     17
  95%     19
  98%     21
  99%     23
 100%     87 (longest request)

ab -kc 10 -t 60 http://localhost:8080/shop/category/electronics/
This is ApacheBench, Version 2.3 <$Revision: 1807734 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking localhost (be patient)
Completed 5000 requests
Completed 10000 requests
Finished 13689 requests


Server Software:        fave.pro/1.3.3
Server Hostname:        localhost
Server Port:            8080

Document Path:          /shop/category/electronics/
Document Length:        5856 bytes

Concurrency Level:      10
Time taken for tests:   60.005 seconds
Complete requests:      13689
Failed requests:        0
Keep-Alive requests:    0
Total transferred:      83434455 bytes
HTML transferred:       80162784 bytes
Requests per second:    228.13 [#/sec] (mean)
Time per request:       43.834 [ms] (mean)
Time per request:       4.383 [ms] (mean, across all concurrent requests)
Transfer rate:          1357.87 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.0      0       1
Processing:     8   44  22.6     40     172
Waiting:        8   43  22.6     40     172
Total:          8   44  22.7     40     172

Percentage of the requests served within a certain time (ms)
  50%     40
  66%     51
  75%     57
  80%     62
  90%     75
  95%     86
  98%    101
  99%    111
 100%    172 (longest request)

ab -kc 10 -t 60 http://localhost:8080/shop/samsung-galaxy-s10/
This is ApacheBench, Version 2.3 <$Revision: 1807734 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking localhost (be patient)
Completed 5000 requests
Completed 10000 requests
Completed 15000 requests
Completed 20000 requests
Completed 25000 requests
Completed 30000 requests
Completed 35000 requests
Completed 40000 requests
Finished 44658 requests


Server Software:        fave.pro/1.3.3
Server Hostname:        localhost
Server Port:            8080

Document Path:          /shop/samsung-galaxy-s10/
Document Length:        10598 bytes

Concurrency Level:      10
Time taken for tests:   60.000 seconds
Complete requests:      44658
Failed requests:        0
Keep-Alive requests:    0
Total transferred:      483958746 bytes
HTML transferred:       473285484 bytes
Requests per second:    744.29 [#/sec] (mean)
Time per request:       13.436 [ms] (mean)
Time per request:       1.344 [ms] (mean, across all concurrent requests)
Transfer rate:          7876.87 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.0      0       0
Processing:     4   13   4.3     13     110
Waiting:        4   13   4.3     13     110
Total:          4   13   4.3     13     110

Percentage of the requests served within a certain time (ms)
  50%     13
  66%     14
  75%     16
  80%     16
  90%     18
  95%     21
  98%     23
  99%     26
 100%    110 (longest request)
```
