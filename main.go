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

	"github.com/will-rowe/herald/src/helpers"
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

	// create the HERALD
	var heraldObj *herald.Herald
	if heraldObj, err = herald.InitHerald(dbLocation); err != nil {
		ui.Eval(fmt.Sprintf("`console.log('failed to init herald: %v')`", err))
	}
	defer heraldObj.Destroy()

	// Bind HERALD methods to the UI
	ui.Bind("wipeStorage", heraldObj.WipeStorage)
	ui.Bind("getSampleCount", heraldObj.GetSampleCount)
	ui.Bind("getUntaggedSampleCount", heraldObj.GetUntaggedSampleCount)
	ui.Bind("getTaggedSampleCount", heraldObj.GetTaggedSampleCount)
	ui.Bind("getAnnouncedSampleCount", heraldObj.GetAnnouncedSampleCount)
	ui.Bind("createSample", heraldObj.CreateSample)
	ui.Bind("deleteSample", heraldObj.DeleteSample)
	ui.Bind("getSampleLabel", heraldObj.GetSampleLabel)
	ui.Bind("printSampleToJSONstring", heraldObj.PrintSampleToJSONstring)
	ui.Bind("createExperiment", heraldObj.CreateExperiment)

	// Bind helper functions to the UI
	ui.Bind("checkDir", helpers.CheckDir)

	// Setup a JS function to init the HERALD and populate all storage data fields in the app
	ui.Bind("loadRuntimeInfo", func() error {

		// load all samples from the storage and populate runtime info
		if err := heraldObj.GetRuntimeInfo(); err != nil {
			return err
		}

		// get the db location and number of experiments and samples in storage
		ui.Eval(fmt.Sprintf(`document.getElementById('staging_experimentCount').innerText = '%d'`, heraldObj.GetExperimentCount()))
		ui.Eval(fmt.Sprintf(`document.getElementById('staging_dbLocation').innerHTML = 'filepath: %v'`, heraldObj.GetDbPath()))
		ui.Eval(fmt.Sprintf(`document.getElementById('staging_sampleCount').innerText = '%d'`, heraldObj.GetSampleCount()))
		ui.Eval(fmt.Sprintf(`document.getElementById('staging_taggedCount').innerText = '%d'`, heraldObj.GetTaggedSampleCount()))

		// check the network connection
		if server.NetworkActive() {
			ui.Eval(fmt.Sprintf(`document.getElementById('status_network').innerHTML = '<i class="far fa-check-circle" style="color: #35cebe;"></i>'`))
		} else {
			ui.Eval(fmt.Sprintf(`document.getElementById('status_network').innerHTML = '<i class="far fa-times-circle" style="color: red;"></i>'`))
		}

		// print a message
		return nil
	})

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

	// Wait until the interrupt signal arrives or browser window is closed
	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
	case <-ui.Done():
	}
}
