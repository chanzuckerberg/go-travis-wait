# go-travis-wait

This is a small utility tool meant to prevent travis from killing long running steps in your build.
Unlike travis_wait, this tool will stream logs back to the user instead of accumulating them until the end.
The downside is we print keep-alive messages periodically to work around this travis limitation.

Usage:
```
go-travis-wait -timeout 30m -interval 1m sleep 10
```
`timeout` specifies the max execution time for this command, we forcefully cancel the build once this deadline exceeds.
`interval` specifies the frequency at which we will print keep-alive messages
