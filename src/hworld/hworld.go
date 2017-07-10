package hworld

import (
	"fmt"
	"os"
)

func main() {
	fmt.Printf("hello, world\n")
	fmt.Fprintf(os.Stdout, "Testing\n")
}
