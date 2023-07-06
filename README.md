# red

![red](https://user-images.githubusercontent.com/141232/54882450-bb85b200-4e8c-11e9-8bd9-37cf43b5b1ed.gif)

_Red_ is a terminal log analysis tools.

## Usage

<a href="https://webpod.dev/?from=red"><img src="https://webpod.dev/img/banner.png" alt="Webpod - deploy JavaScript apps" width="190" align="right"></a>

Pipe JSON stream logs into _red_ and specify a few fields to display. For example using with kubernetes:

```bash
kubectl logs ... | red level message
```

You will see combined logs with trend sparkline and total count.

## Install

```bash
go get github.com/antonmedv/red
```

## License

MIT
