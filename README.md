# go-screenlist

util to create a screenlist file or take screenshots from a video

## cmd

```
Application Options:
  -v, --verbose             verbose mode
  -f, --force               force execution (ignore errors)
  -t, --treads=             number of treads (default: 4)
  -i, --input=              file destinations

screenlist:
      --screenlist.enabled  enable screenlist [$SCREENLIST_ENABLED]
      --screenlist.info     enable info [$SCREENLIST_INFO]
      --screenlist.images=  images directory (default: 15) [$SCREENLIST_IMAGES]
      --screenlist.grid=    grid size (default: 3) [$SCREENLIST_GRID]
      --screenlist.width=   resulting image width (default: 1200) [$SCREENLIST_WIDTH]
      --screenlist.result=  resulting image name. Use {{.Name}} for filename prefix (default: {{.Name}}_screenlist.jpg) [$SCREENLIST_RESULT]

delta:
      --delta.enabled       enable delta saving [$DELTA_ENABLED]
      --delta.images=       images directory (default: 15) [$DELTA_IMAGES]
      --delta.width=        resulting image width (default: 1200) [$DELTA_WIDTH]
      --delta.result=       resulting image name with counter. Use {{.Counter}} for counter value and {{.Name}} for filename prefix (default:
                            {{.Name}}_screenshot_{{.Counter}}.jpg) [$DELTA_RESULT]
```

## example

```bash
screenlist -i video.mp4 \
  --screenlist.enabled \
  --screenlist.info
```

## usage as library

example works as command higher

```bash
go get github.com/ReanSn0w/go-screenlist
```

```go
import (
	"github.com/ReanSn0w/go-screenlist/pkg/engine"
	"github.com/go-pkgz/lgr"
)

func CreateImages(lgr.Default()) error {
	engine := engine.NewEngine(
		logger, false, 4,
		engine.Screenlist{
			Enabled: true,
			Info:    true,
			Images:  15,
			Grid:    3,
			Width:   1200,
			Result:  "{{.Name}}_screenlist.jpg",
		},
		engine.Delta{
			Enabled: false,
		}
	)

	return engine.Run("video.mp4")
}
```

## install

### from source (requires go)

```bash
go install github.com/ReanSn0w/go-screenlist/cmd/screenlist@latest
```

### download binary

go to [releases](https://github.com/ReanSn0w/go-screenlist/releases)

## license

[MIT](https://github.com/ReanSn0w/go-screenlist/blob/main/LICENSE)
