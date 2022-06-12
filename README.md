# ascii-art-web

<h1>Ascii Art WEB</h1>


## Installation
***
Step : .
```
$ git clone https://git.zone01normandie.org/maximediet/ascii-art-web.git
$ cd ascii-art-web
$ go run ascii-art-web.go

```
***
## Usage :
Go to <a href="http://localhost:8080">localhost:8080</a> and type whatever you want, click and you'll have your input in ascii
***
***
## Implementation details: algorithm :
Algorithm write in Go.
Modules:
```golang
import (
    "bufio"
	"color/color" //Local Package in folder color
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)
```
***