// This package is a thin wrapper around libmpg123.
// API documentation [here](http://www.mpg123.de/api/)
// Based on https://bitbucket.org/weberc2/media (no license found)
package mpg123

/*
#cgo LDFLAGS: -lmpg123 -L/usr/local/lib/
#cgo CFLAGS: -I/usr/local/include/
#include "mpg123.h"
*/
import "C"

import (
	"fmt"
	"io"
	"runtime"
	"unsafe"
)

type Mp3ReaderSeekerCloser struct {
	handle *C.mpg123_handle
	outBuf []byte
	offset uint64
	io.Reader
	io.Seeker
	io.Closer
}

type Format int

const (
	Mono8 Format = iota + 1
	Mono16
	Stereo8
	Stereo16
)

func mpg123PlainStrError(errcode C.int) string {
	return C.GoString(C.mpg123_plain_strerror(errcode))
}

func isError(errcode C.int) bool {
	return errcode != C.MPG123_OK
}

func mpg123Error(errcode C.int) error {
	if isError(errcode) {

		// check for io.EOF
		if errcode == C.MPG123_DONE {
			return io.EOF
		}

		s := mpg123PlainStrError(errcode)
		return fmt.Errorf("MPG123 ERR %v: %s", errcode, s)
	}
	return nil
}

func Initialize() {
	if err := mpg123Error(C.mpg123_init()); err != nil {
		Exit()
		panic(err)
	}
}

func Exit() {
	C.mpg123_exit()
}

func Open(path string) (*Mp3ReaderSeekerCloser, error) {
	var errInt C.int
	handle := C.mpg123_new(nil, &errInt)
	if err := mpg123Error(errInt); err != nil {
		return nil, err
	}
	cs := C.CString(path)
	defer C.free(unsafe.Pointer(cs))
	if err := mpg123Error(C.mpg123_open(handle, cs)); err != nil {
		return nil, err
	}
	rsc := &Mp3ReaderSeekerCloser{
		handle: handle,
	}
	runtime.SetFinalizer(rsc, (*Mp3ReaderSeekerCloser).free)
	return rsc, nil
}

func (rsc *Mp3ReaderSeekerCloser) free() {
	C.mpg123_close(rsc.handle)
	C.mpg123_delete(rsc.handle)
}

func (rsc *Mp3ReaderSeekerCloser) Read(p []byte) (int, error) {
	var bytesRead C.size_t
	buffer := (*C.uchar)(unsafe.Pointer(&p[0]))
	bufferSize := (C.size_t)(len(p))
	err := mpg123Error(C.mpg123_read(rsc.handle, buffer, bufferSize, &bytesRead))
	if err != nil {
		return int(bytesRead), err
	}
	return int(bytesRead), nil
}

func (rsc *Mp3ReaderSeekerCloser) Seek(offset int64, whence int) (int64, error) {
	c_offset := (C.off_t)(offset)
	c_whence := (C.int)(whence)
	s_offset := (int64)(C.mpg123_seek(rsc.handle, c_offset, c_whence))
	return s_offset, nil
}

func (rsc *Mp3ReaderSeekerCloser) SeekFrame(offset int64, whence int) int64 {
	c_offset := (C.off_t)(offset)
	c_whence := (C.int)(whence)
	s_offset := (int64)(C.mpg123_seek_frame(rsc.handle, c_offset, c_whence))
	return s_offset
}

func (rsc *Mp3ReaderSeekerCloser) Timeframe(seconds int64) int64 {
	c_seconds := (C.double)(seconds)
	s_offset := (int64)(C.mpg123_timeframe(rsc.handle, c_seconds))
	return s_offset
}

func (rsc *Mp3ReaderSeekerCloser) Close() error {
	return mpg123Error(C.mpg123_close(rsc.handle))
}

func (rsc *Mp3ReaderSeekerCloser) Format() (rate int64, channels int, encoding int, format Format) {
	c_rate := (*C.long)(unsafe.Pointer(&rate))
	c_channels := (*C.int)(unsafe.Pointer(&channels))
	c_encoding := (*C.int)(unsafe.Pointer(&encoding))
	C.mpg123_getformat(rsc.handle, c_rate, c_channels, c_encoding)
	if channels == 1 {
		if encoding&C.MPG123_ENC_8 != 0 {
			format = Mono8
		} else if encoding&C.MPG123_ENC_16 != 0 {
			format = Mono16
		}
	} else if channels == 2 {
		if encoding&C.MPG123_ENC_8 != 0 {
			format = Stereo8
		} else if encoding&C.MPG123_ENC_16 != 0 {
			format = Stereo16
		}
	}
	return
}
