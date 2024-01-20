# container-counter
PoC to count active containers of a cluster.

## Client Mock
Mimic client traffics to continue to send traffic to cluster.

1. Interval is default to 3
2. mode can go from 0 to 10 (0-4 for now)

### 2 clusters with client traffic interval of 3
```text
curl --header "Content-Type: application/json" --request POST --data '{"mode":"2","interval":"3"}' http://localhost:7000/toggle
```

### 4 clusters with client traffic interval of 2
```text
curl --header "Content-Type: application/json" --request POST --data '{"mode":"4","interval":"2"}' http://localhost:7000/toggle
```

### Visit [toggle.rest](toggle.rest) for all toggle mode suggestion

The current setting only spawn 4 cluster containers due to computation resource congestion.

Use rest client vscode plugin to run http requests in [toggle.rest](toggle.rest).

## TODO:
1. Fix redis connection timeout from long connection...