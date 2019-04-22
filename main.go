package main

import (
	"fmt"
	"os"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

// exec executes pipeline and ignore error by merely printing it
func exec(pipe redis.Pipeliner) {
	_, err := pipe.Exec()
	if err != nil {
		errors.Wrap(err, "pipeline exec failed")
	}
}

func main() {
	var (
		pattern string
		addrs   []string
		count   int64
		batch   int
		dryRun  bool
	)

	pflag.Usage = func() {
		fmt.Fprintln(os.Stderr, "redis-iter-del iterates over redis keys with SCAN matched by pattern and then")
		fmt.Fprintln(os.Stderr, "DEL the keys in pipelined commands")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		pflag.PrintDefaults()
	}

	pflag.StringVarP(&pattern, "pattern", "p", "", "Pattern to delete")
	pflag.StringSliceVarP(&addrs, "addrs", "a", []string{":6379"}, "Redis addrs, comma separated for cluster")
	pflag.Int64VarP(&count, "count", "c", 10, "Count for SCAN command")
	pflag.IntVarP(&batch, "batch", "b", 10, "Batch size for pipelined commands")
	pflag.BoolVarP(&dryRun, "dryrun", "d", false, "Dry run")
	pflag.Parse()

	if pattern == "" {
		fmt.Fprintf(os.Stderr, "pattern argument is required\n\n")
		pflag.Usage()
		os.Exit(1)
	}

	redisdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: addrs,
	})

	var total int

	var cursor uint64
	var n int

	pipe := redisdb.Pipeline()
	for {
		var keys []string
		var err error

		// Grab the next portion of keys
		keys, cursor, err = redisdb.Scan(cursor, pattern, count).Result()
		if err != nil {
			errors.Wrap(err, "scan failed")
			break
		}

		// Add DEL commands to the pipeline for each key
		if !dryRun {
			for _, key := range keys {
				pipe.Del(key)
				n++
			}

			// Execute full batch
			if n > batch {
				exec(pipe)
				n = 0
			}
		}

		total += len(keys)

		// Check for the end of iteration
		if cursor == 0 {
			exec(pipe)
			break
		}
	}

	fmt.Printf("iterated over %d keys\n", total)
}
