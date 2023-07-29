# go-screenlist

Утилита для создания изображения с сеткой скриншотов из видеофайла.

## cmd
```
-v, --verbose  verbose mode
    --count=   number of screenshots (default: 16)
    --width=   resulting image width (default: 1920)
    --delta    save delta images
-f, --force    force execution (ignore errors)
    --grid=    grid size (default: 3)
-i, --input=   file destinations
```

## example
```bash
screenlist -i video.mp4
```

## install

### from source
```bash
go install github.com/ReanSn0w/go-screenlist
```