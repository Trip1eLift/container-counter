# container-counter
PoC to count active containers of a cluster.

## Client Mock
Mimic client traffics to continue to send traffic to cluster.

1. Interval is default to 3
2. mode can go from 0 to 10

### 2 clusters with client traffic interval of 3
```text
curl --header "Content-Type: application/json" --request POST --data '{"mode":"1","interval":"3"}' http://localhost:7000/toggle
```

### 5 clusters with client traffic interval of 5
```text
curl --header "Content-Type: application/json" --request POST --data '{"mode":"5","interval":"6"}' http://localhost:7000/toggle
```

## TODO:
1. Fix redis connection timeout from long connection...