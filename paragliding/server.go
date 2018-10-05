package paragliding

import (
	"github.com/marni/goigc"
)

func StartServer() {
	s := "http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc"
	track, err := igc.ParseLocation(s)
}
