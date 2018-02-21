package pastebin

type User struct {
	Name        string `xml:"user_name"`
	Format      string `xml:"user_format_short"`
	Expiration  string `xml:"user_expiration"`
	AvatarUrl   string `xml:"user_avatar_url"`
	Private     string `xml:"user_private"`
	Website     string `xml:"user_website"`
	Email       string `xml:"user_email"`
	Location    string `xml:"user_location"`
	AccountType string `xml:"user_account_type"`
}
