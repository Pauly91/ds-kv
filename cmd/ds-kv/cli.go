package main

import (
	"bufio"
	"ds-kv/pkg/client"
	"fmt"
	"log"
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

	kv, e := client.NewKVClient("")
	if e != nil {
		log.Fatalf("error creating kv client: '%s'", e.Error())
	}

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
			if kv.Set(parts[1], parts[2]) != nil {
				fmt.Println("Error setting value")
			}

		case "get":
			if len(parts) < 2 {
				fmt.Println("Error: get requires key")
				continue
			}
			if val, e := kv.Get(parts[1]); e == nil {
				fmt.Printf("%s = %s\n", parts[1], val)
			} else {
				fmt.Printf("Key '%s' not found\n", parts[1])
			}

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
