package config

var imageExtensions = []string{
	"png",
	"jpg",
	"jpeg",
	"gif",
}

var pageExtensions = []string{
	"wikitext",
	"txt",
}

var allowedExtensionsToBeImported = append(imageExtensions, pageExtensions...)
