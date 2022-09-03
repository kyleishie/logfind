# logfind

## Project Structure
My first thought is to create a simple CLI for server admin to use against an input CSV file, however,
I will always implement such a tool as a sharable package first then build the user interface on top via
an executable.  See pkg and cmd respectively.

### PKG API
The package's API is broken into two components:

1. `reader` - Which is responsible for parsing the log stream. This interface may be reimplemented to support other file types, e.g., log, json, etc.
2. `finder` - Which is responsible for applying a query to the output of the parser. This interface may be reimplemented to support other query languages, e.g., SQL, LogQL, etc.

The intention with these layers is not to show off or be complex for the sake of complexity, rather it serves to create a extensible foundation. See `test/challenge_scenario_test.go`
for proof of scenario of requirement conformance as well as example golang api usage.

#### Example
```
csvFile, _ := os.Open("/path/to/log.csv")

r := csv.NewReader(csvFile)
f := logfind.NewReaderFinder(r)
count, _, _ := f.Find(
  logfind.WhereUsernameEquals("jeff22"),
  logfind.WhereOperationEqual("upload"),
  logfind.WhereTimestampIsBetween(
    time.Date(2020, 04, 15, 00, 00, 00, 0, time.UTC),
    time.Date(2020, 04, 16, 00, 00, 00, 0, time.UTC),
  ),
)
```

### CLI
I have dubbed the cli `lf`, short for `logfind`. `lf` is intended to be simple. See the below examples for how to answer the challenge's scenarios.

#### Scenario 1:
```
How many users accessed the server?

Command:
lf --count=user /path/to/log.csv

Output:
6
```

#### Scenario 2:
```
How many uploads were larger than 50kB?

Command:
lf --minSize=50 --operation=upload /path/to/log.csv

Output:
657 
```

#### Scenario 3:
```
How many times did jeff22 upload to the server on April 15th, 2020?

Command:
lf --username=jeff22 --operation=upload --minTimestamp=2020-04-15T00:00:00.000Z --maxTimestamp=2020-04-16T00:00:00.000Z /path/to/log.csv

Output:
3
```


## TODO
- [ ] CLI shorthand args 
  - [ ] -u instead of --username
  - [ ] -op instead of --operation
- [ ] Support other file types
  - [ ] json
  - [ ] log
  - [ ] any line based file using a regex with named groups to parse each line?
  - [ ] Automatically detect file type based on file extension.
- [ ] Add the ability to output the log events to a file.
- [ ] Make timestamp format configurable


## Notes
If I had the chance to redo this I would probably build something less dependent on known column names and go for a label based approach.
I would then start to question whether that is even necessary given tools like Loki, Prometheus, and Grafana exist.


