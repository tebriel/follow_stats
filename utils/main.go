package utils

import (
	"bufio"
	// "encoding/json"
	"log"
	"os"
)

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		//'scanner.Text()' represents the test case, do something with it
		// stringArray := strings.Split(scanner.Text(), " ")
		// fmt.Println(strings.Join(find_cycle(stringArray), " "))
	}
}
