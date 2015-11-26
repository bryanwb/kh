package hand

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/dullgiulio/pingo"
	"io/ioutil"
	"os"
	"path"
)

var Logger *log.Logger
var Version = "0.0.1"
var Verbose bool
var HandHome string

type Response struct {
	Stdout string
	Stderr string
}

type Hand struct {
	Home    string
	Fingers map[string]string
}

func MakeHand(home string) (*Hand, error) {
	h := new(Hand)
	h.Home = home
	h.Fingers = make(map[string]string)
	err := h.FindFingers()
	return h, err
}

func (h *Hand) FindFingers() error {
	// find all directories w/ executables of same name
	_, err := os.Stat(h.Home)
	if os.IsNotExist(err) {
		return err
	}
	di, _ := ioutil.ReadDir(h.Home)
	for i := range di {
		fingerPath := di[i].Name()
		name := path.Base(fingerPath)
		fullFingerPath := path.Join(h.Home, fingerPath, name)
		if pathHasFinger(fullFingerPath) {
			h.Fingers[name] = fullFingerPath
		}
	}
	return err
}

func (h *Hand) MakeFinger(name string, args []string) (*Finger, error) {
	f := new(Finger)
	f.Path = path.Join(h.Home, name, name)
	f.plugin = pingo.NewPlugin("tcp", f.Path)
	f.args = args
	return f, nil
}

func (h *Hand) ExecuteFinger(name string, args []string) error {
	f, err := h.MakeFinger(name, args)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	if err = f.Execute(); err != nil {
		fmt.Errorf(err.Error())
	}
	return nil
}

type Finger struct {
	Path   string
	plugin *pingo.Plugin
	args   []string
}

func (f *Finger) Execute() error {
	f.plugin.Start()
	defer f.plugin.Stop()
	helpFlags := []string{"-h", "--help", "help"}
	if ContainsAny(helpFlags, f.args) {

	} else {
		resp := new(Response)
		if err := f.plugin.Call("Finger.Execute", f.args, &resp); err != nil {
			log.Print(err)
			return err
		} else {
			fmt.Println(resp.Stdout)
		}
	}
	return nil
}

func (h *Hand) FingerNames() []string {
	names := make([]string, 0)
	for k, _ := range h.Fingers {
		names = append(names, k)
	}
	return names
}

func pathHasFinger(fingerPath string) bool {
	if fi, err := os.Stat(fingerPath); err == nil && fi.Mode()&0111 != 0 {
		Logger.Debugf("Found finger at %s", fingerPath)
		return true
	}
	return false
}

// func main() {
// 	// p := pingo.NewPlugin("tcp", "/Users/hitman/.hand/hello-world/hello-world")
// 	// p.Start()
// 	// defer p.Stop()

// 	// var resp string = ""

// 	// if err := p.Call("MyPlugin.SayHello", "Go Developer", &resp); err != nil {
// 	// 	log.Print(err)
// 	// } else {
// 	// 	log.Print(resp)
// 	// }
// 	log.Out = os.Stderr
// 	formatter := &logrus.TextFormatter{}
// 	formatter.ForceColors = true
// 	log.Formatter = formatter
// 	log.Level = logrus.InfoLevel

// 	var Verbose bool
// 	usr, _ := user.Current()
// 	HandHome = path.Join(usr.HomeDir, "/.hand")

// 	fingerMsg := "No available fingers"
// 	if fingers, err := FindFingers(); err == nil {
// 		fingerMsg = fmt.Sprintf("Available fingers are:\n %s\n",
// 			strings.Join(fingers.Names(), "\n"))
// 	}
// 	var handCmd = &cobra.Command{
// 		Use:   "hand",
// 		Short: "hand is the home your shell scripts always wanted",
// 		Long:  `The Hand of the King (hand) is a tool for organizing and executing "shell" scripts`,
// 		Run: func(cmd *cobra.Command, args []string) {
// 			fmt.Println(`Please pull one of hand's fingers to get started. Type
// "hand help" for more information`)
// 		},
// 	}
// 	handCmd.PersistentFlags().StringVarP(&HandHome, "hand-home", "H", HandHome,
// 		"Home directory for hand")
// 	handCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false,
// 		"verbose output")
// 	handCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
// 		if Verbose == true {
// 			log.Level = logrus.DebugLevel
// 		}
// 	}

// 	var versionCmd = &cobra.Command{
// 		Use:   "version",
// 		Short: "Print the version number of sellsword",
// 		Long:  `All software has versions. This is Sellsword's`,
// 		Run: func(cmd *cobra.Command, args []string) {
// 			fmt.Printf("hand version %s\n", Version)
// 		},
// 	}
// 	handCmd.AddCommand(versionCmd)

// 	var listCmd = &cobra.Command{
// 		Use:   "list",
// 		Short: "Lists available fingers (plugins)",
// 		Long:  `Lists available fingers`,
// 		Run: func(cmd *cobra.Command, args []string) {
// 			fmt.Printf("hand version %s\n", Version)
// 		},
// 	}
// 	handCmd.AddCommand(listCmd)

// 	handCmd.Execute()
// }
