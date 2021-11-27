/*
 @Author      : Simon Chen
 @Email       : bafelem@gmail.com
 @datetime    : 2021-07-23 16:27:49
 @Description : Description
 @FileName    : main_test.go
*/

package finalartwork

import (
	_ "bufio"
	_ "errors"
	_ "fmt"
	_ "log"
	_ "os"
	_ "strings"
	"testing"
)

func TestRun(t *testing.T) {
	MakeEmail("hallo", "this is a test email for your selft")
}
