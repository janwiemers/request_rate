# request_rate

Purpose of this project is to provide a library that can be used to calculate a rate as well as an average.
Both information is calculated on a per second basis and a history is captured on a per second basis as well.

## Get started

This section explains how to integrate and use this library.

```go
include "github.com/janwiemers/request_rate"

rps := NewRequestRate()
rps.Start()

// business logic

log.Printf("You have %i per second", rps.Rate())
```

If you want to measure the request rate

```go
resp, err := http.Get("https://your-url.goes/here")
if err != nil {
   log.Fatalln(err)
}
rps.Incr(1)
```

If you want to measure the duration of your requests or computation alongside with the rate.

```go

uuid := rps.Observe(1)

resp, err := http.Get("https://your-url.goes/here")
if err != nil {
   log.Fatalln(err)
}

duration, err := rps.Finish(uuid)
if err != nil {
  log.Fatal(err)
}
log.Println(duration)

```

## How to Contribute

We'd love to accept your patches and contributions to this project. There are
just a few small guidelines you need to follow.

### Contributor License Agreement

Contributions to this project must be accompanied by a Contributor License
Agreement. You (or your employer) retain the copyright to your contribution;
this simply gives us permission to use and redistribute your contributions as
part of the project. Head over to <https://cla.developers.google.com/> to see
your current agreements on file or to sign a new one.

You generally only need to submit a CLA once, so if you've already submitted one
(even if it was for a different project), you probably don't need to do it
again.

### Code reviews

All submissions, including submissions by project members, require review. We
use GitHub pull requests for this purpose. Consult
[GitHub Help](https://help.github.com/articles/about-pull-requests/) for more
information on using pull requests.

### Community Guidelines

This project follows
[Google's Open Source Community Guidelines](https://opensource.google.com/conduct/).