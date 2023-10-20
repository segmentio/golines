package fixtures

import "fmt"

func testFunc2() {
	c := make(chan int)

	for {
		select {
		case <-c:
			switch a {
			case 1:
				return fmt.Errorf(
					"This is a really long line that can be broken up twice %s %s",
					fmt.Sprintf(
						"This is a really long sub-line that should be broken up more because %s %s",
						"xxxx",
						"yyyy",
					),
					fmt.Sprintf("A short one %d", 3),
				)
			case 2:
			}
		}

		break
	}
}
