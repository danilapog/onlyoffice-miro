package document

const (
	Desktop  EditorMode = "desktop"
	Embedded EditorMode = "embedded"
	Mobile   EditorMode = "mobile"
)

type EditorMode string

type DocumentConfigurer interface {
	ID() string
	FolderID() string
	Title() string
	URL() string
	ModifiedAt() string
}

type UserConfigurer interface {
	ID() string
	Name() string
	Language() string
}

type Permissions struct {
	Comment                 bool `json:"comment"`
	Copy                    bool `json:"copy"`
	DeleteCommentAuthorOnly bool `json:"deleteCommentAuthorOnly"`
	Download                bool `json:"download"`
	Edit                    bool `json:"edit"`
	EditCommentAuthorOnly   bool `json:"editCommentAuthorOnly"`
	FillForms               bool `json:"fillForms"`
	ModifyContentControl    bool `json:"modifyContentControl"`
	ModifyFilter            bool `json:"modifyFilter"`
	Print                   bool `json:"print"`
	Review                  bool `json:"review"`
}

type Document struct {
	Key         string      `json:"key"`
	Title       string      `json:"title"`
	URL         string      `json:"url"`
	FileType    string      `json:"fileType"`
	Permissions Permissions `json:"permissions"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Customization struct {
	Goback struct {
		RequestClose bool `json:"requestClose"`
	} `json:"goback"`
	Plugins       bool `json:"plugins"`
	HideRightMenu bool `json:"hideRightMenu"`
}

type Editor struct {
	User        User   `json:"user"`
	CallbackURL string `json:"callbackUrl"`
	Lang        string `json:"lang"`
}

type Config struct {
	Document     Document `json:"document"`
	DocumentType string   `json:"documentType"`
	Editor       Editor   `json:"editorConfig"`
	Type         string   `json:"type"`
	Token        string   `json:"token"`
}
