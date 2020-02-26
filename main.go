//go:generate go run -tags generate gen.go

package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"

	"github.com/zserge/lorca"

	"github.com/will-rowe/herald/src/herald"
	"github.com/will-rowe/herald/src/server"
)

// dbLocation is where the db is stored TODO: allow user to change this
const dbLocation = "/tmp/herald/db"

// main is the app entrypoint
func main() {

	// setup the UI
	args := []string{}
	if runtime.GOOS == "linux" {
		args = append(args, "--class=Lorca")
	}
	ui, err := lorca.New("", "", 1200, 600, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	// UI is ready
	ui.Bind("start", func() {
		log.Println("UI is ready")
	})

	// Load HTML
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	go http.Serve(ln, http.FileServer(FS))
	ui.Load(fmt.Sprintf("http://%s", ln.Addr()))

	// create the HERALD
	var heraldObj *herald.Herald
	if heraldObj, err = herald.InitHerald(dbLocation); err != nil {
		ui.Eval(fmt.Sprintf("`console.log('failed to init herald: %v')`", err))
	}
	defer heraldObj.Destroy()

	// Bind HERALD methods to the UI
	ui.Bind("wipeStorage", heraldObj.WipeStorage)
	ui.Bind("getSampleCount", heraldObj.GetSampleCount)
	ui.Bind("createSample", heraldObj.CreateSample)
	ui.Bind("deleteSample", heraldObj.DeleteSample)
	ui.Bind("getSampleLabel", heraldObj.GetSampleLabel)
	ui.Bind("printSampleToString", heraldObj.PrintSampleToString)

	// Setup a JS function to populate all storage data fields in the app
	ui.Bind("renderPage", func() {

		// get the db location and number of samples in storage
		ui.Eval(fmt.Sprintf(`document.getElementById('staging_dbLocation').innerHTML = 'filepath: %v'`, heraldObj.GetDbPath()))
		ui.Eval(fmt.Sprintf(`document.getElementById('staging_sampleCount').innerText = '%d'`, heraldObj.GetSampleCount()))

		// check the network connection
		if server.NetworkActive() {
			ui.Eval(fmt.Sprintf(`document.getElementById('status_network').innerHTML = '<i class="far fa-check-circle" style="color: #35cebe;"></i>'`))
		} else {
			ui.Eval(fmt.Sprintf(`document.getElementById('status_network').innerHTML = '<i class="far fa-times-circle" style="color: red;"></i>'`))
		}

		// print a message
		msg := "app loaded successfully"
		ui.Eval(fmt.Sprintf(`document.getElementById('banner_messageBox').innerText = '%v'`, msg))
	})

	// Wait until the interrupt signal arrives or browser window is closed
	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
	case <-ui.Done():
	}
}
