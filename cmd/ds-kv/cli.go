package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Key-Value Store CLI")
	fmt.Println("Commands:")
	fmt.Println("  set <key> <value>")
	fmt.Println("  get <key>")
	fmt.Println("  delete <key>")
	fmt.Println("  exit")

	

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		parts := strings.SplitN(input, " ", 3)

		if len(parts) == 0 {
			continue
		}

		command := strings.ToLower(parts[0])

		switch command {
		case "exit":
			fmt.Println("Goodbye!")
			return

		case "set":
			if len(parts) < 3 {
				fmt.Println("Error: set requires key and value")
				continue
			}

			fmt.Printf("Set %s = %s\n", parts[1], parts[2])

		case "get":
			if len(parts) < 2 {
				fmt.Println("Error: get requires key")
				continue
			}
			/* if val, exists := kvs.Get(parts[1]); exists {
				fmt.Printf("%s = %s\n", parts[1], val)
			} else {
				fmt.Printf("Key '%s' not found\n", parts[1])
			} */

		case "delete":
			if len(parts) < 2 {
				fmt.Println("Error: delete requires key")
				continue
			}

			fmt.Printf("Deleted %s\n", parts[1])

		default:
			fmt.Println("Invalid command. Valid commands: set, get, delete, exit")
		}
	}
}
