package parsers

import (
	"mime"
	"path"
)

func init() {
	mime.AddExtensionType(".md", "text/markdown")
	mime.AddExtensionType(".markdown", "text/markdown")
}

// GetParserFromFileName : Does what it says on the tin.
func GetParserFromFileName(filename string) SpecParser {
	mimeType := mime.TypeByExtension(path.Ext(filename))
	mimeMediaType, _, _ := mime.ParseMediaType(mimeType)

	switch mimeMediaType {
	case "text/markdown":
		return new(Markdown)
	case "text/plain":
		return new(PlainTextParser)
	default:
		return nil
	}
}
