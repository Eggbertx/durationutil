# durationutil
A wrapper around time.ParseDuration that adds flexibility and larger time units (days, weeks, months, years)

## Usage
durationutil provides the function `ParseLongerDuration`, which is meant for parsing strings representing longer durations than the `time` package provides, and in a more verbose, but more readable format, though it does not recognize any unit shorter than a second

```Go
// with units separated by spaces
duration1, err := durationutil.ParseLongerDuration("1 year 2 months 3 days 4 hours 5 minutes 6 seconds")
if err != nil {
	panic(err.Error())
}

// with units together
duration2, err := durationutil.ParseLongerDuration("1year2months3days4hours5minutes6seconds")
if err != nil {
	panic(err.Error())
}

// with units abbreviated and separated by spaces
duration3, err := durationutil.ParseLongerDuration("1y 2mo 3d 4h 5m 6s")
if err != nil {
	panic(err.Error())
}

// with units abbreviated and not separated by spaces
duration4, err := durationutil.ParseLongerDuration("1y2mo3d4h5m6s")
if err != nil {
	panic(err.Error())
}
```
In the above example, `duration1`, `duration2`, `duration3`, and `duration4` are all equal. `ParseLongerDuration` allows you to skip units (for example, `"1y 2d 3s"`), but expects them to be already sorted longest to shortest in descending order.

See [durationutil_test.go](./durationutil_test.go) for more usage info.