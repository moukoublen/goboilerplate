package httplog

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"sync"
)

type BytesBufferPool struct {
	pool *sync.Pool
}

func NewBytesBufferPool(defaultCap int) *BytesBufferPool {
	b := &BytesBufferPool{
		pool: &sync.Pool{
			New: func() any {
				buf := &bytes.Buffer{}
				buf.Grow(defaultCap)
				return buf
			},
		},
	}

	return b
}

func (b *BytesBufferPool) Get() *bytes.Buffer {
	item := b.pool.Get()
	buf, is := item.(*bytes.Buffer)
	if !is {
		return nil
	}
	buf.Reset()

	return buf
}

func (b *BytesBufferPool) Put(buf *bytes.Buffer) {
	b.pool.Put(buf)
}

//
// Tee functionality
//

type TeeCallBack func(readErr, closeErr error, buf *bytes.Buffer)

func NewTeeReadCloserPooled(r io.Reader, pool *BytesBufferPool, cb TeeCallBack) io.ReadCloser {
	teeBuffer := pool.Get()

	t := &teeReadCloser{
		inner: r,
		buf:   teeBuffer,
		cb: func(readErr, closeErr error, buf *bytes.Buffer) {
			cb(readErr, closeErr, buf)
			pool.Put(buf)
		},
	}

	if asWriteTo, isWriteTo := r.(io.WriterTo); isWriteTo {
		writeTo := &teeWriteTo{
			teeReader: t,
			from:      asWriteTo,
		}

		return &teeReadCloserAndWriteTo{teeReadCloser: t, teeWriteTo: writeTo}
	}

	return &teeReadCloserAndWriteTo{teeReadCloser: t}
}

func NewTeeReadCloser(r io.Reader, teeBuffer *bytes.Buffer, cb TeeCallBack) io.ReadCloser {
	t := &teeReadCloser{
		inner: r,
		buf:   teeBuffer,
		cb:    cb,
	}

	if asWriteTo, isWriteTo := r.(io.WriterTo); isWriteTo {
		writeTo := &teeWriteTo{
			teeReader: t,
			from:      asWriteTo,
		}

		return &teeReadCloserAndWriteTo{teeReadCloser: t, teeWriteTo: writeTo}
	}

	return &teeReadCloserAndWriteTo{teeReadCloser: t}
}

// ### io.ReadCloser.
var _ io.ReadCloser = (*teeReadCloser)(nil)

type teeReadCloser struct {
	inner    io.Reader
	readErr  error
	closeErr error
	buf      *bytes.Buffer
	cb       TeeCallBack
	cbOnce   sync.Once
}

func (r *teeReadCloser) Read(p []byte) (n int, err error) {
	n, err = r.inner.Read(p) // tunnel read

	isEOF := errors.Is(err, io.EOF)

	if err != nil && !isEOF {
		r.readErr = err
	}

	if n > 0 {
		if bufN, bufErr := r.buf.Write(p[:n]); bufErr != nil {
			_ = bufN // TODO (or not): bufN != n -> error
			if !isEOF {
				return n, errors.Join(err, bufErr)
			}

			return n, err
		}
	}

	return
}

func (r *teeReadCloser) Buffer() *bytes.Buffer { return r.buf }

func (r *teeReadCloser) Close() error {
	if cl, is := r.inner.(io.Closer); is {
		r.closeErr = cl.Close()
	}

	r.doCB()

	return r.closeErr
}

func (r *teeReadCloser) doCB() {
	if r.cb != nil {
		r.cbOnce.Do(func() {
			r.cb(r.readErr, r.closeErr, r.buf)
		})
	}
}

// ### io.WriterTo.
var (
	_ io.WriterTo = (*teeWriteTo)(nil)
)

type teeWriteTo struct {
	teeReader *teeReadCloser
	from      io.WriterTo //nolint:unused
}

func (r *teeWriteTo) WriteTo(w io.Writer) (n int64, err error) {
	return io.Copy(w, r.teeReader) // tunnel to teeReadCloser.Read
}

// ### io.ReadCloser + io.WriterTo.
var (
	_ io.ReadCloser = (*teeReadCloserAndWriteTo)(nil)
	_ io.WriterTo   = (*teeReadCloserAndWriteTo)(nil)
)

type teeReadCloserAndWriteTo struct {
	*teeReadCloser
	*teeWriteTo
}

//
// Drain functionality
//

func drainRequestBody(req *http.Request) ([]byte, error) {
	if req.GetBody != nil {
		body, err := req.GetBody()
		if err != nil {
			return nil, err
		}

		return io.ReadAll(body)
	}

	buf := &bytes.Buffer{}

	defer func() {
		req.Body = io.NopCloser(buf)
	}()

	if _, err := buf.ReadFrom(req.Body); err != nil {
		return buf.Bytes(), err
	}

	if err := req.Body.Close(); err != nil {
		return buf.Bytes(), err
	}

	return buf.Bytes(), nil
}

func drainResponseBody(req *http.Response) ([]byte, error) {
	buf := &bytes.Buffer{}

	defer func() {
		req.Body = io.NopCloser(buf)
	}()

	if _, err := buf.ReadFrom(req.Body); err != nil {
		return buf.Bytes(), err
	}

	if err := req.Body.Close(); err != nil {
		return buf.Bytes(), err
	}

	return buf.Bytes(), nil
}
