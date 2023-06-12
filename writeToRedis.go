package main

import (
	"fmt"
	"time"
	"io/ioutil"
	"github.com/go-redis/redis/v8"
	"context"
)

var client *redis.Client
var ctx=context.Background()

func initRedis() {
	// Init redis client
	var (
		redisAddr   = "10.202.86.89:6379"
		redisDB     = 0
		redisPasswd = "Saber_Redis_Passwd"
	)
	client = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPasswd, // no password set
		DB:       redisDB,     // use default DB
                // PoolSize: 100,
                // DialTimeout:  10 * time.Second,
                // ReadTimeout:  10 * time.Second,
                // WriteTimeout: 10 * time.Second,
	})
	fmt.Println(ctx, "Redis initialization succeeded")
}

func readFile(file string) []byte{
	f, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("read file %v failed", file)
	}
	return f
}

func writeToFile(data string, fileName string) {
        err := ioutil.WriteFile(fmt.Sprintf("../output_data/%s", fileName), []byte(data), 0644)
        if err != nil {
                fmt.Println("写文件时发生错误:", err)
                return
        }
        fmt.Printf("已向文件写入数据：%s\n", fileName)
}
/*
func writeToRedis(data []byte, key string) error{
        dataSize := len(data)
        if dataSize > maxWriteSize {
                for start := 0; start < dataSize; start += maxWriteSize {
                        end := start + maxWriteSize
                        if end > dataSize {
                                end = dataSize
                        }
                        chunk := data[start : end]
                        err := rdb.Append(ctx, key, string(chunk)).Err()
                        if err != nil {
                                log.Error(ctx, fmt.Sprintf("write in batches to redis error, chunk is start: %v, end: %v, err: %v", start, end, err))
                                panic(fmt.Sprintf("write in batches to redis error, chunk is start: %v, end: %v, err: %v", start, end, err))
                        }
                rdb.Expire(ctx, key, 24*time.Hour)
                }
        } else {
                err := rdb.Set(ctx, key, data, 24*time.Hour).Err()
                if err != nil {
                        log.Error(ctx, fmt.Sprintf("set to redis error, %s", err))
                        panic("Cache set checkResult failed")
                }
        }
        return nil
}
*/


func writeToRedis(data []byte) error {
	const maxWriteSize = 1024 * 1024 // 每次写入最大字节数
	dataSize := len(data)
	value, err := client.ConfigGet(ctx, "proto-max-bulk-len").Result()
	fmt.Println("dataSize: ", dataSize, "proto-max-bulk-len: ", value)
	if err != nil {
		fmt.Println("get config proto-max-bulk-len is: ", value)
	}
	for i := 0; i < dataSize; i += maxWriteSize {
		time.Sleep(1 * time.Second)
		end := i + maxWriteSize
		if end > dataSize {
			end = dataSize
		}
		chunk := data[i:end]
		err = client.Append(ctx, "key1", string(chunk)).Err()
		if err != nil {
			fmt.Println(fmt.Sprintf("write error, start: %v, end: %v", i, end))
			return err
		}
		fmt.Println(fmt.Sprintf("have write %v", end))
	}

	return nil
}

func getDataFromRedis() {
	value, e := client.Get(ctx, "key1").Result()
	if e != nil {
		fmt.Println("get data from redis error: ", e)
	}
	fmt.Println("len value: ", len(value))
	writeToFile(value, "get_data_1.txt")
}


func main() {
	initRedis()

	// data := []byte("large amount of data")
	data := readFile("/root/gocode/output_data/tmp_render_data_storge.txt")
	err := writeToRedis(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Data has been written to Redis successfully!")
	getDataFromRedis()

}

/*
func main() {
	initRedis()
	getDataFromRedis()
}
*/

