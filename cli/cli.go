package main

import (
	"dds/cli/terminal"
	"fmt"
	// queue "github.com/emirpasic/gods/queues/linkedlistqueue"
)

func main() {

	//queue.New()

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
