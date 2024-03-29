package email

type Attachment struct {
	FileName string
	Content  []byte
}

type Attachments []Attachment
