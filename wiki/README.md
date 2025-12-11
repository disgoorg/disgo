# DisGo Wiki

## Quick Start

If you have DevContainers set up (and i highly recommend it), you can simply open this repository in a compatible IDE (like VSCode) and it will set up the environment for you.

Pre-requisites: [Hugo](https://gohugo.io/getting-started/installing/), [Go](https://golang.org/doc/install) and [Git](https://git-scm.com)


## Running the Wiki Locally
```shell
hugo mod tidy
hugo server --logLevel debug --disableFastRender -p 1313
```

Then open your browser to `http://localhost:1313` to view the wiki locally.
