# imgbb

`imgbb` is a toolkit to help upload images to [ImgBB](https://imgbb.com).

## Installation

Via Golang package get command

```sh
go get -u github.com/wabarc/imgbb/cmd/imgbb
```

From [gobinaries.com](https://gobinaries.com):

```sh
$ curl -sf https://gobinaries.com/wabarc/imgbb | sh
```

## Usage

Command-line:

```sh
$ imgbb
A CLI tool help upload images to ImgBB.

Usage:

  imgbb [options] [file1] ... [fileN]

  -k string
    	ImgBB api key, optional.
```

Go package:
```go
import (
        "fmt"

        "github.com/wabarc/imgbb"
)

func main() {
        if url, err := i.Upload(path); err != nil {
            fmt.Fprintf(os.Stderr, "imgbb: %v\n", err)
        } else {
            fmt.Fprintf(os.Stdout, "%s  %s\n", url, path)
        }
}
```

## License

This software is released under the terms of the GNU General Public License v3.0. See the [LICENSE](https://github.com/wabarc/imgbb/blob/main/LICENSE) file for details.
