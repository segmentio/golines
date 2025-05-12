package fixtures

import (
	"fmt"
	"time"
)

func _() {
	veryLongString := fmt.Sprintf("The current time is %v and %d, %d, %d, %d, %d, %d", time.Now(), 1, 2, 3, 4, 5, 6)
	anotherLongString := fmt.Sprintf("The current time is %v and %d, %d, %d, %d, %d, %d, %d, %d, %d", time.Now(), 1, 2, 3, 4, 5, 6, 7, 8, 9)

	extremelyLongSlice := []string{"This is the first very long element in a slice", "This is the second very long element that contains a lot of text to make the line exceed 120 characters", "And this is the third one"}

	someVeryLongMap := map[string]string{"ThisIsAVeryLongKeyThatWillMakeThisLineExceed120Characters": "AndThisIsAnEquallyLongValueThatWillContributeToMakingThisLineVeryVeryLongIndeed"}

	_ = veryLongString
	_ = anotherLongString
	_ = extremelyLongSlice
	_ = someVeryLongMap
}
