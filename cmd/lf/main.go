package main

import (
	"flag"
	"fmt"
	"github.com/kyleishie/logfind/pkg/logfind"
	"github.com/kyleishie/logfind/pkg/logfind/reader/csv"
	"os"
	"time"
)

func main() {

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] filepath\n\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}

	countConcernPtr := flag.String("count", "event", "Changes how lf counts events. Values are event, operation, user.  Event is default.")
	minTimestampPtr := flag.String("minTimestamp", "", "The minimum date to match.")
	maxTimestampPtr := flag.String("maxTimestamp", "", "The maximum date to match. Note this is exclusive.")
	usernamePtr := flag.String("username", "", "The usernamePtr to match.")
	operationPtr := flag.String("operation", "", "The operationPtr to match.")
	minSizePtr := flag.Int("minSize", -1, "The minimum size to match.")
	maxSizePtr := flag.Int("maxSize", -1, "The maximum size to match. Note this in inclusive.")
	verbosePtr := flag.Bool("verbose", false, "Use this flag to print more info about the matched events.")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("missing filepath")
		flag.Usage()
		os.Exit(1)
	}

	filepath := args[len(args)-1]

	csvFile, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	r := csv.NewReader(csvFile)
	f := logfind.NewFinder(r)

	var opts []logfind.FinderOptionFunc

	if countConcernPtr != nil {
		opts = append(opts, logfind.WithCountConcern(logfind.CountConcern(*countConcernPtr)))
	}

	if *minTimestampPtr != "" && *maxTimestampPtr != "" {
		minTimestamp, err := time.Parse(time.RFC3339, *minTimestampPtr)
		maxTimestamp, err := time.Parse(time.RFC3339, *maxTimestampPtr)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(2)
		}

		opts = append(opts, logfind.WhereTimestampIsBetween(minTimestamp, maxTimestamp))
	}

	if *usernamePtr != "" {
		opts = append(opts, logfind.WhereUsernameEquals(*usernamePtr))
	}

	if *operationPtr != "" {
		opts = append(opts, logfind.WhereOperationEquals(*operationPtr))
	}

	if *minSizePtr != -1 {
		opts = append(opts, logfind.WhereSizeGreaterThanOrEqual(*minSizePtr))
	}

	if *maxSizePtr != -1 {
		opts = append(opts, logfind.WhereSizeGreaterThanOrEqual(*maxSizePtr))
	}

	count, events, err := f.Find(opts...)

	fmt.Printf("count: %d\n", count)

	if *verbosePtr {
		for _, event := range events {
			fmt.Println(event)
		}
	}
}
