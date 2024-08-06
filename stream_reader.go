package openai

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"

	utils "github.com/iimeta/go-openai/internal"
)

var (
	headerData  = []byte("data: ")
	errorPrefix = []byte(`data: {"error":`)
)

type streamable interface {
	ChatCompletionStreamResponse | CompletionResponse
}

type streamReader[T streamable] struct {
	emptyMessagesLimit uint
	isFinished         bool

	reader         *bufio.Reader
	response       *http.Response
	errAccumulator utils.ErrorAccumulator
	unmarshaler    utils.Unmarshaler

	httpHeader
}

func (stream *streamReader[T]) Recv() (responseBytes []byte, response T, err error) {
	if stream.isFinished {
		err = io.EOF
		return
	}

	responseBytes, response, err = stream.processLines()
	return
}

//nolint:gocognit
func (stream *streamReader[T]) processLines() ([]byte, T, error) {
	var (
		emptyMessagesCount uint
		hasErrorPrefix     bool
	)

	for {
		rawLine, readErr := stream.reader.ReadBytes('\n')
		if readErr != nil || hasErrorPrefix {
			respErr := stream.unmarshalError()
			if respErr != nil {
				return nil, *new(T), fmt.Errorf("error, %w", respErr.Error)
			}
			return nil, *new(T), readErr
		}

		noSpaceLine := bytes.TrimSpace(rawLine)
		if bytes.HasPrefix(noSpaceLine, errorPrefix) {
			hasErrorPrefix = true
		}
		if !bytes.HasPrefix(noSpaceLine, headerData) || hasErrorPrefix {
			if hasErrorPrefix {
				noSpaceLine = bytes.TrimPrefix(noSpaceLine, headerData)
			}
			writeErr := stream.errAccumulator.Write(noSpaceLine)
			if writeErr != nil {
				return nil, *new(T), writeErr
			}
			emptyMessagesCount++
			if emptyMessagesCount > stream.emptyMessagesLimit {
				return nil, *new(T), ErrTooManyEmptyStreamMessages
			}

			continue
		}

		noPrefixLine := bytes.TrimPrefix(noSpaceLine, headerData)
		if string(noPrefixLine) == "[DONE]" {
			stream.isFinished = true
			return nil, *new(T), io.EOF
		}

		var response T
		unmarshalErr := stream.unmarshaler.Unmarshal(noPrefixLine, &response)
		if unmarshalErr != nil {
			return nil, *new(T), unmarshalErr
		}

		return noPrefixLine, response, nil
	}
}

func (stream *streamReader[T]) unmarshalError() (errResp *ErrorResponse) {
	errBytes := stream.errAccumulator.Bytes()
	if len(errBytes) == 0 {
		return
	}

	err := stream.unmarshaler.Unmarshal(errBytes, &errResp)
	if err != nil {
		errResp = nil
	}

	return
}

func (stream *streamReader[T]) Close() error {
	return stream.response.Body.Close()
}
