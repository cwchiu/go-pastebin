package pastebin

type Paste struct {
	Key         string `xml:"paste_key"`
	Date        string `xml:"paste_date"`
	Title       string `xml:"paste_title"`
	Private     string `xml:"paste_private"` // 0=public 1=unlisted 2=private
	ExpireDate  string `xml:"paste_expire_date"`
	FormatLong  string `xml:"paste_format_long"`
	FormatShort string `xml:"paste_format_short"`
	Url         string `xml:"paste_url"`
	Hits        string `xml:"paste_hits"`
	Code        string
	MemberOnly  bool
}

type PasteList struct {
	Pastes []Paste `xml:"paste"`
}

func CreateNormalPaste(title, code string) *Paste {
	return &Paste{
		Code:       code,
		Title:      title,
		MemberOnly: false,
		Private:    "2",
	}
}
