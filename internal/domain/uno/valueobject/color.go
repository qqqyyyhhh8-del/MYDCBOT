package valueobject

type Color string

const (
	ColorRed    Color = "Red"
	ColorBlue   Color = "Blue"
	ColorGreen  Color = "Green"
	ColorYellow Color = "Yellow"
	ColorWild   Color = "Wild"
)

func (c Color) IsValid() bool {
	switch c {
	case ColorRed, ColorBlue, ColorGreen, ColorYellow, ColorWild:
		return true
	}
	return false
}

func (c Color) String() string {
	return string(c)
}
