package form

import (
	"github.com/smartwalle/binding"
	"net/http"
	"strings"
)

//
const (
	kFormTag = "form"

	//	k_FORM_NO_TAG              = "-"
	//	k_FORM_CLEANED_DATA        = "CleanedData"
	//	k_FORM_DEFAULT_FUNC_PREFIX = "Default"
)

func BindWithRequest(request *http.Request, dst interface{}) (err error) {
	err = request.ParseForm()
	if err != nil {
		return err
	}

	var contentType = request.Header.Get("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		request.ParseMultipartForm(32 << 20)
	}

	err = Bind(request.Form, dst)
	return err
}

func Bind(src map[string][]string, dst interface{}) (err error) {
	var source = make(map[string]interface{})
	for key, value := range src {
		if len(value) > 1 {
			source[key] = value
		} else {
			source[key] = value[0]
		}
	}
	return binding.BindWithTag(source, dst, kFormTag)
}
