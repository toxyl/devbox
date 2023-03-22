package core

import (
	"fmt"
	"os"
	"strings"

	"github.com/toxyl/devbox/bashcompletion"
	"github.com/toxyl/glog"
)

type command struct {
	Name       string
	Usage      string
	Example    string
	Help       string
	Args       ArgInfoList
	argsMax    int
	argsMin    int
	argsOffset int
	hidden     bool
	run        func(arg ...string) error
}

func (ac *command) helpTextUsage() string {
	u := reToken.ReplaceAllStringFunc(ac.Usage, func(s string) string {
		return glog.Auto(s)
	})
	return glog.Bold() + APP_NAME + glog.Reset() + " " + glog.Underline() + glog.Auto(ac.Name) + glog.Reset() + " " + u
}

func (ac *command) helpTextUsagePadded(maxLen, maxLenUsage int) string {
	u := reToken.ReplaceAllStringFunc(ac.Usage, func(s string) string {
		return glog.Auto(s)
	})
	return glog.Bold() + APP_NAME + glog.Reset() + " " + glog.Underline() + glog.PadRight(glog.Auto(ac.Name), maxLen, ' ') + glog.Reset() + " " + glog.PadRight(u, maxLenUsage+1, ' ') + glog.WrapDarkGreen("// "+glog.StripANSI(strings.Split(ac.Help, "\n")[0]))
}

func (ac *command) helpTextExample() string {
	return APP_NAME + " " + ac.Name + " " + ac.Example
}

func (ac *command) helpText() string {
	return fmt.Sprintf(
		"%s\n%s\n\n%s\n\n",
		glog.Bold()+"Usage:   "+glog.Reset()+"sudo "+ac.helpTextUsage(),
		glog.Bold()+"Example: "+glog.Reset()+"sudo "+ac.helpTextExample(),
		ac.Help,
	)
}

func (ac *command) is(op string) bool {
	return ac.Name == op
}

func (ac *command) Run() error {
	offsetArgs := ac.argsOffset
	maxArgs := ac.argsMax
	minArgs := ac.argsMin
	numArgs := len(os.Args) - offsetArgs
	hasArgs := numArgs > 0
	hasTooFewArgs := minArgs > -1 && numArgs < minArgs
	hasTooManyArgs := maxArgs > -1 && numArgs > maxArgs

	if !hasArgs && minArgs > 0 {
		fmt.Print(ac.helpText())
		os.Exit(EXIT_OK)
	}

	if hasTooFewArgs {
		fmt.Print(ac.helpText())
		fmt.Println(glog.WrapOrange("Not enough arguments!") + " Got: " + glog.Auto(os.Args[1:]))
		os.Exit(EXIT_MISSING_ARGS)
	}

	if hasTooManyArgs {
		fmt.Print(ac.helpText())
		fmt.Println(glog.WrapOrange("Too many arguments!") + " Got: " + glog.Auto(os.Args[1:]))
		os.Exit(EXIT_TOO_MANY_ARGS)
	}

	args := os.Args[offsetArgs:]

	return ac.run(args...)
}

type ArgInfo struct {
	Optional       bool
	Variadic       bool
	Type           int
	Name           string
	Example        string
	OptionalSuffix string
}

func (ai *ArgInfo) Completion() bashcompletion.Completion {
	switch ai.Type {
	case ARG_TYPE_DIR:
		return bashcompletion.Completion{
			Variadic:    ai.Variadic,
			Completions: "($(compgen -d -- $cur))",
			Compopt:     []string{"plusdirs"},
			Shopt:       []string{},
		}
	case ARG_TYPE_TARBALL:
		return bashcompletion.Completion{
			Variadic:    ai.Variadic,
			Completions: "($(compgen -f -X '!(*.tar.gz|*.tar.xz)' -- $cur))",
			Compopt:     []string{"plusdirs"},
			Shopt:       []string{"extglob"},
		}
	case ARG_TYPE_FILE:
		return bashcompletion.Completion{
			Variadic:    ai.Variadic,
			Completions: "($(compgen -f -- $cur))",
			Compopt:     []string{"plusdirs"},
			Shopt:       []string{},
		}
	case ARG_TYPE_DEVBOX_LIST:
		return bashcompletion.Completion{
			Variadic:    ai.Variadic,
			Completions: "($(compgen -W '$(for dir in " + GetDevboxDir() + "/*/; do basename \"$dir\"; done)' -- $cur))",
			Compopt:     []string{},
			Shopt:       []string{},
		}

	case ARG_TYPE_WORKSPACE_LIST:
		return bashcompletion.Completion{
			Variadic:    ai.Variadic,
			Completions: "($(compgen -W '$(for dir in " + GetWorkspaceDir() + "/*/; do basename \"$dir\"; done)' -- $cur))",
			Compopt:     []string{},
			Shopt:       []string{},
		}
	}
	return bashcompletion.Completion{
		Variadic:    ai.Variadic,
		Completions: "($(compgen -f -- $cur))",
		Compopt:     []string{"plusdirs"},
		Shopt:       []string{},
	}
}

func (ai *ArgInfo) Strings() []string {
	suffix := ""
	if ai.OptionalSuffix != "" {
		if ai.Optional {
			suffix = fmt.Sprintf("{%s}", ai.OptionalSuffix)
		} else {
			suffix = fmt.Sprintf("<%s>", ai.OptionalSuffix)
		}
	}

	if ai.Variadic {
		// we need to return 4 results, where the first can be optional
		first := fmt.Sprintf("%s_1%s", ai.Name, suffix)
		if ai.Optional {
			first = fmt.Sprintf("<%s>", first)
		} else {
			first = fmt.Sprintf("[%s]", first)
		}

		// reset suffix if needed
		if ai.OptionalSuffix != "" {
			suffix = fmt.Sprintf("{%s}", ai.OptionalSuffix)
		}

		return []string{first, fmt.Sprintf("<%s_2>", ai.Name), "..", fmt.Sprintf("<%s_n%s>", ai.Name, suffix)}
	}

	if ai.Optional {
		// single optional arg
		return []string{fmt.Sprintf("<%s%s>", ai.Name, suffix)}
	}

	// single required arg
	return []string{fmt.Sprintf("[%s%s]", ai.Name, suffix)}
}

func (ai *ArgInfo) String() string {
	return strings.Join(ai.Strings(), " ")
}

func (ai *ArgInfo) StringANSI() string {
	res := []string{}
	for _, s := range ai.Strings() {
		res = append(res, glog.Auto(s))
	}
	return strings.Join(res, " ")
}

type ArgInfoList []ArgInfo

func (ail *ArgInfoList) Example() string {
	res := []string{}
	for _, ai := range *ail {
		res = append(res, ai.Example)
	}
	return strings.Join(res, " ")
}

func (ail *ArgInfoList) String() string {
	res := []string{}
	for _, ai := range *ail {
		res = append(res, ai.String())
	}
	return strings.Join(res, " ")
}

func (ail *ArgInfoList) StringANSI() string {
	res := []string{}
	for _, ai := range *ail {
		res = append(res, ai.StringANSI())
	}
	return strings.Join(res, " ")
}

func (ail *ArgInfoList) NumArgs() (min, max int) {
	min = 0
	max = len(*ail)

	for _, ai := range *ail {
		if ai.Variadic {
			max = -1
			if !ai.Optional {
				min++
			}
			continue
		}
		if !ai.Optional {
			min++
		}
	}
	return
}

func newCommand(
	op string,
	help string,
	args ArgInfoList,
	run func(arg ...string) error,
	hidden bool,
) *command {
	argsMin, argsMax := args.NumArgs()
	ac := &command{
		Name:       op,
		Usage:      args.String(),
		Example:    args.Example(),
		Args:       args,
		Help:       help,
		argsMin:    argsMin,
		argsMax:    argsMax,
		argsOffset: 2,
		run:        run,
		hidden:     hidden,
	}

	return ac
}

func RegisterCommand(
	op string,
	help string,
	args ArgInfoList,
	run func(arg ...string) error,
) {
	cmdReg = append(cmdReg, newCommand(
		op,
		help,
		args,
		run,
		false,
	))
}

func RegisterHiddenCommand(
	op string,
	args ArgInfoList,
	run func(arg ...string) error,
) {
	cmdReg = append(cmdReg, newCommand(
		op,
		"(internal use)",
		args,
		run,
		true,
	))
}

func FindCommand(op string) *command {
	for _, cmd := range cmdReg {
		if cmd.is(op) {
			return cmd
		}
	}
	return nil
}

func ListCommands() {
	maxLenName := GetMaxCommandNameLength()
	maxLenUsage := GetMaxCommandUsageLength()
	for _, cmd := range cmdReg {
		if cmd.hidden {
			continue
		}
		fmt.Println("  sudo " + cmd.helpTextUsagePadded(maxLenName, maxLenUsage))
	}
}

func GetMaxCommandNameLength() int {
	maxLen := 0
	for _, cmd := range cmdReg {
		if cmd.hidden {
			continue
		}
		maxLen = glog.Max(maxLen, len(cmd.Name))
	}

	return maxLen
}

func GetMaxCommandUsageLength() int {
	maxLen := 0
	for _, cmd := range cmdReg {
		if cmd.hidden {
			continue
		}
		maxLen = glog.Max(maxLen, len(cmd.Usage))
	}

	return maxLen
}

func GetCommandNames() []string {
	res := []string{}
	for _, cmd := range cmdReg {
		if cmd.hidden {
			continue
		}
		res = append(res, cmd.Name)
	}
	return res
}

func GetCommandData() map[string][]bashcompletion.Completion {
	res := map[string][]bashcompletion.Completion{}
	for _, cmd := range cmdReg {
		if cmd.hidden {
			continue
		}
		res[cmd.Name] = []bashcompletion.Completion{}
		for _, a := range cmd.Args {
			res[cmd.Name] = append(res[cmd.Name], a.Completion())
		}
	}

	return res
}
