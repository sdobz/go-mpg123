Go bindings for the mpg123 mp3 decoding library
http://www.mpg123.de/

These bingings were heavily influenced by weberc2's work here: https://bitbucket.org/weberc2/media

Usage:

First init the library
mpg123.Initialize()

Open a file and check the error
mp3, err := mpg123.Open("example.mp3")
if err != nil {
	panic(err)
}
errors are passed verbatim from mpg123

Get the format
rate, channels, encoding, format := mp3.Format()
fmt.Printf("Rate: %i Channels: %i Encoding: %i Format: %s\n", rate, channels, encoding, format)

format is compatible with audio.Format consts

Read 100 bytes of mp3 data into a buffer

data := make([]byte, 100)
err := mp3.Read(data)

mp3 conforms to the io.Reader, io.Seeker, io.Closer interfaces, though seeking is not yet implemented.

This means it can be used quite handily with the golang.org/x/mobile/exp/audio openal bindings.
