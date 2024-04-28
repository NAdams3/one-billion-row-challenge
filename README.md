# one billion row challenge

This is a custom version of the one billion row challenge to better understand concurrency in go

Also running bets on how long it takes me to fry my CPU lol ...

Status:
Attempt 1: Unable to process 1_000_000_000 rows, OOM error
Attempt 2: 709.18s

> [!CAUTION]
> Use at your own risk!

### usage
run:
```Shell
go build -o build/app
```

to generate a file with n lines:
```Shell
./build/app -generate -count n
```

to process a measurements file with n rows, Note: you'll have to generate the file firest if it doesn't exist
```Shell
./build/app -count n
```

