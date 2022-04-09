package cmd

import "github.com/dhuan/wikicmd/pkg/mw"

var imageExtensions = []string{
	"png",
	"jpg",
	"jpeg",
	"gif",
}

var pageExtensions = []string{
	"wikitext",
}

var allowedExtensionsToBeImported = append(imageExtensions, pageExtensions...)

var MAP_UPLOAD_WARNING_MESSAGE = map[mw.UploadWarning]string{
	mw.UPLOAD_WARNING_SAME_FILE_NO_CHANGE: "File was not uploaded because the existing image is exactly the same.",
}
