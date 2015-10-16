Go bindings for the [mpg123](http://www.mpg123.de/) mp3 decoding library

These bingings were heavily influenced by [weberc2's work](https://bitbucket.org/weberc2/media)

Usage:

First init the library
```golang
mpg123.Initialize()
```

Open a file and check the error
```golang
mp3, err := mpg123.Open("example.mp3")
if err != nil {
	panic(err)
}
```
errors are passed verbatim from mpg123

Get the format
```golang
rate, channels, encoding, format := mp3.Format()
fmt.Printf("Rate: %i Channels: %i Encoding: %i Format: %s\n", rate, channels, encoding, format)
```

format is compatible with audio.Format consts

Read 100 bytes of mp3 data into a buffer
```golang
data := make([]byte, 100)
err := mp3.Read(data)
```

mp3 conforms to the io.Reader, io.Seeker, io.Closer interfaces, though seeking is not yet implemented.

This means it can be used quite handily with the [go audio library](https://godoc.org/golang.org/x/mobile/exp/audio).
