package icon

type Validator interface {
	IsValidMimeType(mimeType string) (bool, error)
}
