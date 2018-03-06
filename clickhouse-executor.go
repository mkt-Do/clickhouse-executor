package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/joho/sqltocsv"
	_ "github.com/kshvakov/clickhouse"
	"log"
	"os"
)

const BUFSIZE = 1024

type Config struct {
	ClickHouseConfig ClickHouseConfig
}

type ClickHouseConfig struct {
	User     string
	Password string
	Host     string
	Port     int
}

func main() {
	confFile := flag.String("conf", "clickhouse.conf", "specify clickhouse configuration toml file")
	queryFile := flag.String("query", "default.sql", "sepcify query file")
	outFile := flag.String("out", "output.csv", "specify output file name(extension is only csv). if not specify default file name.")
	flag.Parse()

	if *confFile == "" || *queryFile == "" {
		log.Fatal("Please specify --conf file --query file")
		os.Exit(1)
	}

	fmt.Println(*outFile)

	var config Config
	_, err := toml.DecodeFile(*confFile, &config)
	if err != nil {
		log.Fatal("cannot read configuration this toml file")
		os.Exit(1)
	}

	chUrl := fmt.Sprintf("tcp://%v:%v", config.ClickHouseConfig.Host, config.ClickHouseConfig.Port)
	connect, err := sql.Open("clickhouse", chUrl)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	readFile, _ := os.Open(*queryFile)
	defer readFile.Close()

	buf := make([]byte, BUFSIZE)
	num := 0
	for {
		n, err := readFile.Read(buf)
		if n == 0 {
			break
		}
		if err != nil {
			log.Fatal("cannot read file")
			os.Exit(1)
		}
		num = n
	}
	query := string(buf[:num])
	fmt.Println(query)

	rows, err := connect.Query(query)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	if err := sqltocsv.WriteFile(*outFile, rows); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	rows.Close()
}
