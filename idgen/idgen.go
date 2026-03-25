package idgen

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/aldge/gopkg/idgen/snowflake"
)

type IDGenerator interface {
	GenID() (int64, error)
}

var DefaultIDGenerator IDGenerator

func init() {
	workID := os.Getenv("WORK_ID")
	if workID == "" {
		log.Fatal("WORK_ID environment variable is not set")
	}

	dataCenterID := os.Getenv("DATA_CENTER_ID")
	if dataCenterID == "" {
		log.Fatal("DATA_CENTER_ID environment variable is not set")
	}

	workIDInt64, err := strconv.ParseInt(workID, 10, 64)
	if err != nil {
		log.Fatalf("invalid WORK_ID: %v", err)
	}

	dataCenterIDInt64, err := strconv.ParseInt(dataCenterID, 10, 64)
	if err != nil {
		log.Fatalf("invalid DATA_CENTER_ID: %v", err)
	}

	DefaultIDGenerator, err = snowflake.NewSnowflake(workIDInt64, dataCenterIDInt64)
	if err != nil {
		panic("failed to initialize snowflake ID generator: " + err.Error())
	}
}

func GenID() (int64, error) {
	return DefaultIDGenerator.GenID()
}

func GenIDStr(prefix string) (string, error) {
	id, err := DefaultIDGenerator.GenID()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%d", prefix, id), nil
}
