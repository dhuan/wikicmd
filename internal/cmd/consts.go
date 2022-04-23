package cmd

import "github.com/dhuan/wikicmd/pkg/mw"

var MAP_UPLOAD_WARNING_MESSAGE = map[mw.UploadWarning]string{
	mw.UPLOAD_WARNING_SAME_FILE_NO_CHANGE: "File was not uploaded because the existing image is exactly the same.",
}

var export_types = []string{
	export_type_all,
	export_type_page,
	export_type_image,
}

var export_type_all = "all"
var export_type_page = "page"
var export_type_image = "image"
