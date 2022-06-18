package main

import (
	"fmt"
	queue "github.com/emirpasic/gods/queues/linkedlistqueue"
	"httk/cli/terminal"
)

func main() {

	queue.New()

	/*if ok {
		switch v := val.(type) {
		case *test:
			fmt.Println("============:  ", v.id, v.name)
		case int:

		}

	}*/

	err := terminal.OpenShell()
	if err != nil {
		fmt.Println(err)
	}

}
