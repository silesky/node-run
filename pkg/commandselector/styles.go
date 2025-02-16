package commandselector

type AnsiColors struct {
	magenta  string
	charcoal string
	blue     string
	purple   string
	white    string
	yellow   string
}

var (
	Colors = AnsiColors{
		white:    "15",
		magenta:  "205",
		charcoal: "236",
		blue:     "20",
		purple:   "55",
		yellow:   "226",
	}
)
