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
	"time"

	"github.com/zserge/lorca"

	"github.com/will-rowe/herald/src/helpers"
	"github.com/will-rowe/herald/src/herald"
	"github.com/will-rowe/herald/src/services"
)

// dbLocation is where the db is stored - it is set at compile time to be platform specific
var dbLocation string

// getSampleServiceTagsHTML collects the registered services for the
// provided record type and returns the HTML block to display them
// to the user.
func getServiceTagsHTML(recordType string) string {
	serviceTagsHTML := "<label>Service requests</label>"
	for serviceName, service := range services.ServiceRegister {
		if recordType == service.GetRecordType().String() {
			serviceTagsHTML += fmt.Sprintf("<input type=\"checkbox\" id=\"formLabel_%v\" value=\"%v\"><label class=\"label-inline\" for=\"formLabel_%v\"> - %v</label><div class=\"clearfix\"></div>", serviceName, serviceName, serviceName, serviceName)
		}
	}
	return serviceTagsHTML
}

// getServiceStatusHTML checks the status of all registered
// services and returns the HTML block to display the
// service names and statuses.
func getServiceStatusHTML() string {
	currentTime := time.Now()
	serviceStatusHTML := ""
	for serviceName, service := range services.ServiceRegister {
		if service.CheckAccess() == true {
			serviceStatusHTML += fmt.Sprintf("<div class=\"mt-1\"><div class=\"float-left\"><i class=\"far fa-check-circle\" style=\"color: #35cebe;\"></i></div><div class=\"float-left ml-1\"><p class=\"m-0\"><strong>%s</strong> <span class=\"text-muted\">service</span></p><p class=\"text-small text-muted\">checked at %d:%d</p></div><div class=\"clearfix\"></div></div>", serviceName, currentTime.Hour(), currentTime.Minute())
		} else {
			serviceStatusHTML += fmt.Sprintf("<div class=\"mt-1\"><div class=\"float-left\"><i class=\"far fa-times-circle\" style=\"color: red;\"></i></div><div class=\"float-left ml-1\"><p class=\"m-0\"><strong>%s</strong> <span class=\"text-muted\">service</span></p><p class=\"text-small text-muted\">checked at %d:%d</p></div><div class=\"clearfix\"></div></div>", serviceName, currentTime.Hour(), currentTime.Minute())
		}
	}
	return serviceStatusHTML
}

// main is the app entrypoint
func main() {

	// setup lorca
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
		ui.Eval(fmt.Sprintf(`console.log('failed to init herald: %v')`, err))
	}
	defer heraldObj.Destroy()

	// Bind HERALD methods to the UI
	// buttons
	ui.Bind("addRun", heraldObj.AddRun)
	ui.Bind("createSample", heraldObj.CreateSample)
	ui.Bind("deleteSample", heraldObj.DeleteSample)
	ui.Bind("announceSamples", heraldObj.AnnounceSamples)
	ui.Bind("wipeStorage", heraldObj.WipeStorage)
	ui.Bind("getUser", heraldObj.GetUser)
	ui.Bind("editConfig", heraldObj.EditConfig)
	// counters
	ui.Bind("getRunCount", heraldObj.GetRunCount)
	ui.Bind("getSampleCount", heraldObj.GetSampleCount)
	ui.Bind("getUntaggedCount", heraldObj.GetUntaggedCount)
	ui.Bind("getTaggedIncompleteCount", heraldObj.GetTaggedIncompleteCount)
	ui.Bind("getTaggedCompleteCount", heraldObj.GetTaggedCompleteCount)
	ui.Bind("getAnnouncementQueueSize", heraldObj.GetAnnouncementQueueSize)
	ui.Bind("getAnnouncementCount", heraldObj.GetAnnouncementCount)
	// table / modals / forms
	ui.Bind("getRunName", heraldObj.GetLabel)
	ui.Bind("getSampleLabel", heraldObj.GetSampleLabel)
	ui.Bind("getSampleCreation", heraldObj.GetSampleCreation)
	ui.Bind("getSampleRun", heraldObj.GetSampleRun)
	ui.Bind("printSampleToJSONstring", heraldObj.PrintSampleToJSONstring)
	ui.Bind("printConfigToJSONstring", heraldObj.PrintConfigToJSONstring)

	// Bind helper functions to the UI
	ui.Bind("checkDirExists", helpers.CheckDirExists)
	ui.Bind("getServiceStatusHTML", getServiceStatusHTML)
	ui.Bind("getPrimerSchemes", heraldObj.GetPrimerSchemes)

	// Setup a JS function to init the HERALD and populate all storage data fields in the app
	ui.Bind("loadRuntimeInfo", func() error {

		// load all samples from the storage and populate runtime info
		if err := heraldObj.GetRuntimeInfo(); err != nil {
			return err
		}

		// print the db location and number of runs and samples in storage etc.
		ui.Eval(fmt.Sprintf(`document.getElementById('staging_runCount').innerText = '%d'`, heraldObj.GetRunCount()))
		ui.Eval(fmt.Sprintf(`document.getElementById('staging_runRequests').innerHTML = '%d with completed service requests'`, heraldObj.GetTaggedCompleteCount("runs")))
		ui.Eval(fmt.Sprintf(`document.getElementById('staging_sampleCount').innerText = '%d'`, heraldObj.GetSampleCount()))
		ui.Eval(fmt.Sprintf(`document.getElementById('staging_sampleRequests').innerText = '%d with completed service requests'`, heraldObj.GetTaggedCompleteCount("samples")))
		ui.Eval(fmt.Sprintf(`document.getElementById('stagingAnnouncementQueueCount').innerText = '%d'`, heraldObj.GetAnnouncementQueueSize()))
		ui.Eval(fmt.Sprintf(`document.getElementById('stagingAnnouncementCount').innerText = '%d announcements made'`, heraldObj.GetAnnouncementCount()))

		// enable the add sample button if there are runs to use
		if heraldObj.GetRunCount() == 0 {
			ui.Eval(`document.getElementById('addSampleModalOpen').disabled = true`)
		} else {
			ui.Eval(`document.getElementById('addSampleModalOpen').disabled = false`)
		}

		// update the add run form with the available services for tagging
		ui.Eval(fmt.Sprintf(`document.getElementById('runTags').innerHTML = '%v'`, getServiceTagsHTML("run")))

		// update the add sample form with the available services for tagging
		ui.Eval(fmt.Sprintf(`document.getElementById('sampleTags').innerHTML = '%v'`, getServiceTagsHTML("sample")))

		// enable the announce button if there are tagged service requests for runs/samples in the queue
		if heraldObj.GetAnnouncementQueueSize() == 0 {
			ui.Eval(`document.getElementById('stagingAnnounce').disabled = true`)
		} else {
			ui.Eval(`document.getElementById('stagingAnnounce').disabled = false`)
		}

		// check the network connection and update the service tags
		if helpers.NetworkActive() {
			ui.Eval(`document.getElementById('status_network').innerHTML = '<i class="far fa-check-circle" style="color: #35cebe;"></i>'`)
		} else {
			ui.Eval(`document.getElementById('status_network').innerHTML = '<i class="far fa-times-circle" style="color: red;"></i>'`)
		}
		ui.Eval(fmt.Sprintf(`document.getElementById('serviceStatus').innerHTML = '%s'`, getServiceStatusHTML()))

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
