# audio

The audio module provides opus/pcm audio encoding/decoding as C bindings based on the [libopus](https://github.com/xiph/opus) library.
It also lets you combine multiple pcm streams into a single pcm stream.
This module requires [CGO](https://go.dev/blog/cgo) to be enabled.

## Getting Started

### Installing

```sh
$ go get github.com/disgoorg/disgo/audio
```
