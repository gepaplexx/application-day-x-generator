package util

import (
	"github.com/DrSmithFr/go-console/color"
	"github.com/DrSmithFr/go-console/formatter"
	"github.com/DrSmithFr/go-console/output"

	"fmt"
	"strings"

	"atomicgo.dev/cursor"
)

var FmtGreenbold = formatter.NewOutputFormatterStyle(color.Green, color.Null, []string{color.Bold})
var FmtRedbold = formatter.NewOutputFormatterStyle(color.Red, color.Null, []string{color.Bold})
var FmtCyanbold = formatter.NewOutputFormatterStyle(color.Cyan, color.Null, []string{color.Bold})

func PrintDescription(conf string) {
	out := output.NewCliOutput(true, nil)
	str := strings.Repeat("#", 80)
	out.Println(str)
	out.Println("# <b>Name:</b><comment> Day-X-Operations Generator</comment>")
	out.Println("# <b>Description:</b><comment> TODO</comment>")
	out.Println("# <b>Author:</b><comment> gattma,fhochleitner</comment>")
	out.Println("# <b>Version:</b><comment> v1.0</comment>")
	out.Println("# <b>Documentation:</b><comment> https://gepardec.atlassian.net/wiki/spaces/G/pages/2393276417/Day-2-Operations</comment>")
	configLine := fmt.Sprintf("# <b>Configuration:</b><comment> %s</comment>", conf)
	out.Println(configLine)
	out.Println(str)
}

func PrintActionHeader(title string) {
	fmt.Println()
	out := output.NewCliOutput(true, nil)
	n := LINELENGTH - len(title)
	div := n / 2
	out.Println("<fg=yellow>" + strings.Repeat("#", LINELENGTH) + "</>")
	out.Println(strings.Repeat(" ", div) + title)
	out.Println("<fg=yellow>" + strings.Repeat("#", LINELENGTH) + "</>")
}

func PrintAction(action string) {
	fmt.Print(FmtCyanbold.Apply(action))
}

func PrintSuccess() {
	cursor.StartOfLine()
	cursor.Right(LINELENGTH - 2)
	fmt.Printf(FmtGreenbold.Apply("OK\n"))
}

func PrintFailure() {
	cursor.StartOfLine()
	cursor.Right(LINELENGTH - 6)
	fmt.Printf(FmtRedbold.Apply("FAILED\n"))
}

func WaitToContinue() {
	// https://golangbyexample.com/how-to-pause-a-go-program-until-enter-key-is-pressed/
	fmt.Println("Press enter to continue...")
	fmt.Scanln()
}
