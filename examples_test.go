package timezonermax_test

import (
	"fmt"
	"time"

	timezoner "github.com/Jhonatan-Code-dev/timezonermax"
)

func ExampleConvert() {
	utcTime := time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC)
	limaTime, err := timezoner.Convert(utcTime, "America/Lima")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(limaTime.Format("15:04"))
	// Output:
	// 07:00
}

func ExampleAt() {
	utcTime := time.Date(2026, 6, 1, 15, 0, 0, 0, time.UTC)

	formatted, err := timezoner.At(utcTime).
		In("America/Lima").
		Format("15:04 MST")

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(formatted)
	// Output:
	// 10:00 -05
}

func ExampleDifference() {
	at := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	diff, err := timezoner.Difference("Europe/London", "America/Lima", at)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(diff)
	// Output:
	// 5h0m0s
}
