# ClickHouse Executor(Written By Golang)
Execute clickhouse query, and output csv file

## Command Options
`conf`(required): clickhouse config(Written by toml file)  
`query`(required): clickhouse execute query file  
`out`(option): filename to write results  

## How to Use
```
go run clickhouse-executor.go --conf clickhouse.conf --query query.sql --out output.csv
```
if build
```
go build clickhouse-executor.go
./clickhouse-executor --conf clickhouse.conf --query query.sql --out output.csv
```

