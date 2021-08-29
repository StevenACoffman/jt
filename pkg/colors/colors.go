package colors

import "strconv"

// ANSIColor is a color (0-15) as defined by the ANSI Standard.
type ANSIColor int

const (
	ANSIBlack         ANSIColor = iota // "#000000"
	ANSIRed                            // "#800000"
	ANSIGreen                          // "#008000"
	ANSIYellow                         // "#808000"
	ANSIBlue                           // "#000080"
	ANSIMagenta                        // "#800080"
	ANSICyan                           // "#008080"
	ANSIWhite                          // "#808080"
	ANSIBrightBlack                    // "#c0c0c0"
	ANSIBrightRed                      // "#ff0000"
	ANSIBrightGreen                    // "#00ff00"
	ANSIBrightYellow                   // "#ffff00"
	ANSIBrightBlue                     // "#ff0000"
	ANSIBrightMagenta                  // "#ff00ff"
	ANSIBrightCyan                     // "#00ffff"
	ANSIBrightWhite                    // "#ffffff"
)

func (c ANSIColor) String() string {
	return strconv.Itoa(int(c))
}

func (c ANSIColor) Hex() string {
	return [...]string{
		"#000000", // ANSIBlack
		"#800000", // ANSIRed
		"#008000", // ANSIGreen
		"#808000", // ANSIYellow
		"#000080", // ANSIBlue
		"#800080", // ANSIMagenta
		"#008080", // ANSICyan
		"#808080", // ANSIWhite
		"#c0c0c0", // ANSIBrightBlack
		"#ff0000", // ANSIBrightRed
		"#00ff00", // ANSIBrightGreen
		"#ffff00", // ANSIBrightYellow
		"#ff0000", // ANSIBrightBlue
		"#ff00ff", // ANSIBrightMagenta
		"#00ffff", // ANSIBrightCyan
		"#ffffff", // ANSIBrightWhite
	}[c]
}
