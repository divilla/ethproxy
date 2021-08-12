#!/bin/bash

(echo "127.0.0.1:8080/block/12979835/transaction/1"; echo "127.0.0.1:8080/block/12979835/transaction/2"; echo "127.0.0.1:8080/block/12979835/transaction/3"; echo "127.0.0.1:8080/block/12979835/transaction/4"; echo "127.0.0.1:8080/block/12979835/transaction/5"; echo "127.0.0.1:8080/block/12979835/transaction/6") \
  | parallel 'ab -v 3 -n 20000 -c 200 {}'
