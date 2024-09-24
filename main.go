package main

import (
	"Umka-1/satellite"
	"fmt"
	"time"
)

func main() {
	name := "UMKA-1 (RS40S)"
	line1 := "1 57172U 23091G   24263.53334166  .00009425  00000-0  59089-3 0  9999"
	line2 := "2 57172  97.6018 314.6827 0017222 154.9337 205.2732 15.09427738 67710"

	sat := satellite.New(line1, line2, name)

	coords, err := sat.Calculate(time.Now().UTC().Add(time.Hour * 24 * 7))
	if err != nil {
		panic("ошибка в Calculate")
	}

	fmt.Println(coords)

	// line1 = "1 57172U 23091G   24263.53334166  .00009425  00000-0  59089-3 0  9998"
	// line2 = "2 57172  97.6018 314.6827 0017222 154.9337 205.2732 15.09427738 67711"

	// sat.UpdateTLE(line1, line2)

	// fmt.Println(sat)
}
