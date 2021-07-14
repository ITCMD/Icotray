package acticon

import "regexp"

type Identifier string

func (identifier Identifier) IsValid() (bool, error) {
	return regexp.Match(`^[\w\-_]{3,}$`, []byte(identifier))
}
