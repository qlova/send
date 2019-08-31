package send

import "io"

//Attachment can be addded to messages.
type Attachment struct {
	string
	reader io.Reader
}
