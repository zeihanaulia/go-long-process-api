package response

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var (
	ErrInternalServer = NewError("Internal Server Error", http.StatusInternalServerError)
)

type Error struct {
	ErrMessage string `json:"error,omitempty"`
	ErrCode    int    `json:"error_code,omitempty"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s with code %d", e.ErrMessage, e.ErrCode)
}

func NewError(msg string, code int) *Error {
	return &Error{
		ErrMessage: msg,
		ErrCode:    code,
	}
}

type HTTPResponse interface {
	WriteResponse(w http.ResponseWriter)
}

type BasicResponse struct {
	Body        []byte
	StatusCode  int
	ContentType string
}

func (b *BasicResponse) WriteResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", b.ContentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(b.Body)))
	w.WriteHeader(b.StatusCode)
	if _, err := w.Write(b.Body); err != nil {
		log.Println("unable to write byte.", err)
	}
}

type JSONResponse struct {
	BasicResponse
	JSONBody JSONBody
	Error    Error
}

type JSONBody struct {
	*Error
	Data interface{} `json:"data,omitempty"`
}

const JSONContentType = "application/json"

func NewJSONResponse() *JSONResponse {
	return &JSONResponse{
		BasicResponse: BasicResponse{
			ContentType: JSONContentType,
			StatusCode:  http.StatusOK,
		},
	}
}

func (r *JSONResponse) WriteJSONResponse(rw http.ResponseWriter, err error) {
	rw.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(err)
	rw.Write(b)
}

func (r *JSONResponse) SetBody(data interface{}) *JSONResponse {
	r.JSONBody.Data = data
	return r
}

func (r *JSONResponse) SetError(err error) *JSONResponse {
	if respErr, ok := err.(*Error); ok {
		r.JSONBody.Error = respErr
	} else {
		// when unspecified error is provided it will categorize the response as internal server error
		r.JSONBody.Error = NewError(err.Error(), http.StatusInternalServerError)
	}
	return r
}

func (r *JSONResponse) WriteResponse(w http.ResponseWriter) {
	b, err := json.Marshal(r.JSONBody)
	if err != nil {
		JSONBody := JSONBody{
			Error: NewError(err.Error(), http.StatusInternalServerError),
		}
		b, _ = json.Marshal(JSONBody)
	}
	r.Body = b
	if r.JSONBody.Error != nil {
		r.StatusCode = r.JSONBody.ErrCode
	}
	r.BasicResponse.WriteResponse(w)
}
