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
	"github.com/will-rowe/herald/src/minknow"
	"github.com/will-rowe/herald/src/services"
)

// dbLocation is where the db is stored - it is set at compile time to be platform specific
var dbLocation string

// getTagsHTML returns the HTML needed to display all available services for sample tagging
func getTagsHTML() string {
	ServiceTagsHTML := "<label>Tags</label>"
	for serviceName := range services.ServiceRegister {
		if serviceName == "sequence" || serviceName == "basecall" {
			continue
		}
		ServiceTagsHTML += fmt.Sprintf("<input type=\"checkbox\" id=\"formLabel_%v\" value=\"%v\"><label class=\"label-inline\" for=\"formLabel_%v\">%v</label><div class=\"clearfix\"></div>", serviceName, serviceName, serviceName, serviceName)
	}
	return ServiceTagsHTML
}

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

	// get the available processes for tagging
	ServiceTagsHTML := getTagsHTML()

	// create the HERALD
	var heraldObj *herald.Herald
	if heraldObj, err = herald.InitHerald(dbLocation); err != nil {
		ui.Eval(fmt.Sprintf(`console.log('failed to init herald: %v')`, err))
	}
	defer heraldObj.Destroy()

	// Bind HERALD methods to the UI
	// buttons
	ui.Bind("createExperiment", heraldObj.CreateExperiment)
	ui.Bind("createSample", heraldObj.CreateSample)
	ui.Bind("deleteSample", heraldObj.DeleteSample)
	ui.Bind("announceSamples", heraldObj.AnnounceSamples)
	ui.Bind("wipeStorage", heraldObj.WipeStorage)
	// counters
	ui.Bind("getExperimentCount", heraldObj.GetExperimentCount)
	ui.Bind("getSampleCount", heraldObj.GetSampleCount)
	ui.Bind("getUntaggedSampleCount", heraldObj.GetUntaggedRecordCount)
	ui.Bind("getTaggedSampleCount", heraldObj.GetTaggedRecordCount)
	ui.Bind("getAnnouncementCount", heraldObj.GetAnnouncementCount)
	// table / modals / forms
	ui.Bind("getExperimentName", heraldObj.GetLabel)
	ui.Bind("getSampleLabel", heraldObj.GetSampleLabel)
	ui.Bind("getSampleCreation", heraldObj.GetSampleCreation)
	ui.Bind("getSampleExperiment", heraldObj.GetSampleExperiment)
	ui.Bind("printSampleToJSONstring", heraldObj.PrintSampleToJSONstring)

	// Bind helper functions to the UI
	ui.Bind("checkDirExists", helpers.CheckDirExists)
	ui.Bind("checkAPIstatus", minknow.CheckAPIstatus)

	// Setup a JS function to init the HERALD and populate all storage data fields in the app
	ui.Bind("loadRuntimeInfo", func() error {

		// load all samples from the storage and populate runtime info
		if err := heraldObj.GetRuntimeInfo(); err != nil {
			return err
		}

		// print the db location and number of experiments and samples in storage etc.
		ui.Eval(fmt.Sprintf(`document.getElementById('staging_dbLocation').innerHTML = 'filepath: %v'`, heraldObj.GetDbPath()))
		ui.Eval(fmt.Sprintf(`document.getElementById('staging_experimentCount').innerText = '%d'`, heraldObj.GetExperimentCount()))
		ui.Eval(fmt.Sprintf(`document.getElementById('staging_sampleCount').innerText = '%d'`, heraldObj.GetSampleCount()))
		ui.Eval(fmt.Sprintf(`document.getElementById('staging_taggedCount').innerText = '%d'`, heraldObj.GetTaggedRecordCount()))
		ui.Eval(fmt.Sprintf(`document.getElementById('staging_processCount').innerText = '%d untagged'`, heraldObj.GetUntaggedRecordCount()))

		// enable the add sample button if there are experiments to use
		if heraldObj.GetExperimentCount() == 0 {
			ui.Eval(`document.getElementById('addSampleModalOpen').disabled = true`)
		} else {
			ui.Eval(`document.getElementById('addSampleModalOpen').disabled = false`)
		}

		// update the add sample form with the available processes for tagging
		ui.Eval(fmt.Sprintf(`document.getElementById('sampleTags').innerHTML = '%v'`, ServiceTagsHTML))

		// enable the announce button if there are tagged samples
		if heraldObj.GetTaggedRecordCount() == 0 {
			ui.Eval(`document.getElementById('staging_announce').disabled = true`)
		} else {
			ui.Eval(`document.getElementById('staging_announce').disabled = false`)
		}

		// check the network connection
		if helpers.NetworkActive() {
			ui.Eval(`document.getElementById('status_network').innerHTML = '<i class="far fa-check-circle" style="color: #35cebe;"></i>'`)
		} else {
			ui.Eval(`document.getElementById('status_network').innerHTML = '<i class="far fa-times-circle" style="color: red;"></i>'`)
		}

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

	// alert if a new release is available
	updateAvailable, releaseVersion, releaseLink, err := helpers.CheckLatestRelease()
	if err != nil {
		ui.Eval(fmt.Sprintf(`printErrorMsg('%v')`, err))
	}
	if updateAvailable {
		ui.Eval(fmt.Sprintf(`printUpdateMsg('a new version is available (%v), download now?', '%v')`, releaseVersion, releaseLink))
	}

	// Wait until the interrupt signal arrives or browser window is closed
	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
	case <-ui.Done():
	}
}
