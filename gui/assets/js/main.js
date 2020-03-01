////////////////////////////////////////////////////////////////////
// HELPERS
// removeOptions will clear a select box
function removeOptions(selectbox) {
    var i
    for (i = selectbox.options.length - 1; i >= 0; i--) {
        selectbox.remove(i)
    }
}

// handleForm will prevent default form action
function handleForm(event) {
    event.preventDefault()
}

////////////////////////////////////////////////////////////////////
// MESSAGES
// printErrorMsg
function printErrorMsg(msg) {
    console.log('error: ', msg)
    $('body').overhang({
        type: 'error',
        message: msg,
        overlay: true,
        closeConfirm: true
    })
}

// printSuccessMsg
function printSuccessMsg(msg) {
    console.log(msg)
    $('body').overhang({
        custom: true,
        primary: '#35cebe',
        accent: '#25beae',
        message: msg
    })
}

// updatStatus will print status and time
function updateStatus(id, active) {
    var today = new Date()
    var time =
        today.getHours() +
        ':' +
        (today.getMinutes() < 10 ? '0' : '') +
        today.getMinutes()
    if (active === 'true') {
        document.getElementById(id).innerHTML =
            '<i class="far fa-check-circle" style="color: #35cebe;"></i>'
        document.getElementById(id + '_timestamp').innerText =
            'online since: ' + time
    } else {
        document.getElementById(id).innerHTML =
            '<i class="far fa-times-circle" style="color: red;"></i>'
        document.getElementById(id + '_timestamp').innerText =
            'offline since: ' + time
    }
}

// update network status
window.addEventListener('online', () => updateStatus('status_network', 'true'))
window.addEventListener('offline', () =>
    updateStatus('status_network', 'false')
)

////////////////////////////////////////////////////////////////////
// BUTTONS
// get the buttons that control the app
const refreshPage = document.getElementById('refreshPage')
const wipeDatabase = document.getElementById('wipeDatabase')

// add an event listener to the refreshPage button
refreshPage.addEventListener('click', async() => {
    console.log('refreshing the app')

    pageRefresh()
    printSuccessMsg('refreshed the app')
})

// add an event listener to wipeDatabase button
wipeDatabase.addEventListener('click', async() => {
    console.log('wiping storage')

    // TODO: add a confirm prompt

    // call the Go wipeStorage method
    try {
        await wipeStorage()
    } catch (e) {
        console.log(e)
        printErrorMsg(e)
        return
    }

    // reset the page and report success
    fullPageRender()
    printSuccessMsg('storage wiped')
})

////////////////////////////////////////////////////////////////////
// MODALS
// modal closing
var modalClosers = document.getElementsByClassName('modal-close')
for (let i = 0; i < modalClosers.length; i++) {
    modalClosers[i].addEventListener('click', function() {
        modalClosers[i].closest('.modal').style.display = 'none'
    })
}
window.onclick = function(event) {
    if (event.target.className == 'modal') {
        event.target.style.display = 'none'
    }
}

// createExperimentModal
// open on button click
const createExperimentModal = document.getElementById('createExperimentModal')
const createExperimentModalOpen = document.getElementById(
    'createExperimentModalOpen'
)
createExperimentModalOpen.addEventListener('click', function() {
    createExperimentModal.style.display = 'block'
})

// addSampleModal
// open on button click (Go handles the disabled flag)
const addSampleModal = document.getElementById('addSampleModal')
const addSampleModalOpen = document.getElementById('addSampleModalOpen')
addSampleModalOpen.addEventListener('click', function() {
    addSampleModal.style.display = 'block'
})

// sampleDetailsModal
const sampleDetailsModal = document.getElementById('sampleDetailsModal')

// getSampleJSONdump returns a stringified protobuf dump of a sample from the storage
const getSampleJSONdump = async function(sampleLabel) {
    var sampleJSONdump = `${await window.printSampleToJSONstring(sampleLabel)}`
    return sampleJSONdump
}

////////////////////////////////////////////////////////////////////
// CREATE EXPERIMENT FORM
// get the form
const createExperimentForm = document.getElementById('createExperimentForm')
createExperimentForm.addEventListener('submit', handleForm)

// add listener to the output location text box so we can validate it
var outputLocation = document.getElementById('formLabel_outputLocation')
outputLocation.addEventListener('change', async() => {
    try {
        await checkDir(outputLocation.value)
    } catch (e) {
        printErrorMsg(e)
        return
    }
})

// add an event listener to the createExperimentForm submit button
createExperimentForm.addEventListener('submit', async() => {
    console.log('creating experiment')

    var elements = createExperimentForm.elements

    // check for data directories if the user is creating an experiment with existing data
    if (elements['formLabel_sequenced'].value == 'true') {
        // TODO: check this properly, at the moment just looking for experiment name dir and fast5 sub dir
        var dirName =
            elements['formLabel_outputLocation'].value +
            '/' +
            elements['formLabel_experimentName'].value
        var fast5_dirName = dirName + '/fast5_pass'
        try {
            await checkDir(dirName)
        } catch (e) {
            printErrorMsg(e)
            return
        }
        try {
            await checkDir(fast5_dirName)
        } catch (e) {
            printErrorMsg(e)
            return
        }
    }

    // create an experiment and add it to the store
    try {
        await createExperiment(
            elements['formLabel_experimentName'].value,
            elements['formLabel_outputLocation'].value
        )
    } catch (e) {
        printErrorMsg(e)
        return
    }

    // reset the form, refresh the page, close the modal and report success
    createExperimentForm.reset()
    pageRefresh()
    createExperimentModal.style.display = 'none'
    printSuccessMsg('experiment created')
})

////////////////////////////////////////////////////////////////////
// SAMPLE SUBMISSION FORM
// get the form
const addSampleForm = document.getElementById('addSampleForm')
addSampleForm.addEventListener('submit', handleForm)

// updateExperimentDropDown is a function to update the experiments in the sample submission form
const expDropDown = document.getElementById('formLabel_sampleExperiment')
const updateExperimentDropDown = async() => {
    // get the current experiment count so that we can iterate over the experiments
    var expCount = `${await window.getExperimentCount()}`

    // if there are no experiments, leave the default blank option
    if (expCount === '0') {
        return
    }

    // wipe the current options
    removeOptions(expDropDown)

    // add each name to the drop down
    for (var i = 0; i < expCount; i++) {
        var expName = `${await window.getExperimentName(i)}`
        var newOpt = document.createElement('option')
        newOpt.text = expName
        expDropDown.options.add(newOpt)
    }
}

// add listener to the form tags so that more form appears
document
    .getElementById('formLabel_rampartTag')
    .addEventListener('change', async() => {
        var options = document.getElementById('addSampleForm_rampartOptions')
        if (options.style.display === 'block') {
            options.style.display = 'none'
        } else {
            options.style.display = 'block'
        }
    })

// add an event listener to the addSampleForm submit button
addSampleForm.addEventListener('submit', async() => {
    console.log('adding sample to storage')

    var elements = addSampleForm.elements

    // grab the tags from the form
    var tags = []
    for (var i = 0, element;
        (element = elements[i++]);) {
        if (element.type === 'checkbox' && element.checked) {
            tags.push(element.value)
        }
    }

    // create a sample and add it to the storage
    try {
        // TODO: try reading form straight into protobuf and then send a serialised stream to Go
        await createSample(
            elements['formLabel_sampleLabel'].value,
            elements['formLabel_sampleExperiment'].value,
            parseInt(elements['formLabel_sampleBarcode'].value, 10),
            elements['formLabel_sampleComment'].value,
            tags
        )
    } catch (e) {
        printErrorMsg(e)
        return
    }

    // update the table
    $('#sampleTable')
        .DataTable()
        .row.add([
            elements['formLabel_sampleLabel'].value,
            elements['formLabel_sampleExperiment'].value
        ])
        .draw(true)

    // reset the form, refresh the page and report success
    addSampleForm.reset()
    pageRefresh()
    printSuccessMsg('sample added')
})

////////////////////////////////////////////////////////////////////
// TABLES
// set up the table
var table = $('#sampleTable').DataTable({
    columnDefs: [{
        targets: 2,
        data: null,
        searchable: true,
        orderable: true,
        defaultContent: '<button class="button button-outline">Manage</button>'
    }]
})

// set up the manage button
$('#sampleTable tbody').on('click', 'button', function() {
    var row = table.row($(this).parents('tr'))
    var data = row.data()
    var sampleLabel = data[0]
    document.getElementById('sampleModal_samplename').innerHTML = sampleLabel
    getSampleJSONdump(sampleLabel).then(sampleProtobufDump => {
        // get the sample protobuf dump
        document.getElementById('sampleModal_content').innerHTML =
            '<pre>' + sampleProtobufDump + '</pre>'

        // display modal
        document.getElementById('sampleDetailsModal').style.display = 'block'

        // set up delete button
        document
            .getElementById('sampleModal_delete')
            .addEventListener('click', async() => {
                console.log('deleting sample')

                // delete from the db
                try {
                    await deleteSample(sampleLabel)
                } catch (e) {
                    console.log(e)
                    printErrorMsg(e)
                    return
                }

                // remove from the table
                row.remove().draw(true)

                // reset the runtime info and report success
                await pageRefresh()
                document.getElementById('sampleDetailsModal').style.display = 'none'
                printSuccessMsg('sample deleted')
            })
    })
})

// buildTable will get the database keys via Go and then populate the table
const buildTable = async() => {
    console.log('building table from sample labels in storage')

    // wipe any existing table
    table.clear().draw(true)

    // get the current sample number so that we can iterate over the samples
    var sampleCount = `${await window.getSampleCount()}`

    // process each sample label
    for (var i = 0; i < sampleCount; i++) {
        var sampleLabel = `${await window.getSampleLabel(i)}`
            //var sampleCreation = `${await window.getSampleCreation(i)}`
        var sampleExperiment = `${await window.getSampleExperiment(i)}`

        // create the table entry
        table.row.add([sampleLabel, sampleExperiment]).draw(true)
    }
}

////////////////////////////////////////////////////////////////////
// PAGE RENDERING
// setup the time stamps
const dateStamp = document.getElementById('dateStamp')
const refreshStamp = document.getElementById('refreshStamp')

// printTimeStamps will add the date and the refresh time to the app
function printTimeStamps() {
    var today = new Date()
    var options = {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
    }
    dateStamp.innerHTML = today.toLocaleDateString('en-US', options)
    refreshStamp.innerHTML =
        'last refreshed: ' +
        today.getHours() +
        ':' +
        (today.getMinutes() < 10 ? '0' : '') +
        today.getMinutes()
}

// set up the empty pie chart
var pieCanvas = document.getElementById('pieChart')
var pieData = {
    labels: ['Active Announcements', 'Tagged Samples', 'Untagged Samples'],
    datasets: [{
        label: 'entry point',
        data: [0, 0, 0],
        backgroundColor: ['#35cebe', '#a0a0a0', '#dfdfdf'],
        hoverBackgroundColor: ['#25beae', '#999999', '#cccccc']
    }]
}
var pieOptions = {
    responsive: true,
    segmentShowStroke: false,
    legend: false
}
var myPieChart = new Chart(pieCanvas, {
    type: 'pie',
    data: pieData,
    options: pieOptions
})

// updatePieChart will refresh the pie chart with current data
const updatePieChart = async() => {
    // get counts
    var untaggedSampleCount = `${await window.getUntaggedSampleCount()}`
    var taggedSampleCount = `${await window.getTaggedSampleCount()}`
    var activeAnnouncementCount = `${await window.getAnnouncedSampleCount()}`

    // update the chart data
    myPieChart.data.datasets[0].data[0] = activeAnnouncementCount
    myPieChart.data.datasets[0].data[1] = taggedSampleCount
    myPieChart.data.datasets[0].data[2] = untaggedSampleCount

    // update the chart
    myPieChart.update()
}

// pageRefresh will refresh the Herald runtime info in Go and then freshen up the page (does not rebuild the table)
const pageRefresh = async() => {
    console.log('refreshing runtime info and re-rendering the page')

    // reload the Go Herald instance and repopulate the page data
    try {
        await loadRuntimeInfo()
    } catch (e) {
        printErrorMsg(e)
        return
    }

    // update the experiment drop down
    await updateExperimentDropDown()

    // update the pie chart
    await updatePieChart()

    // check the minKNOW status TODO: do this in Go via routine
    var minknowStatus = `${await checkAPIstatus()}`
    updateStatus('status_minknow', minknowStatus)

    // print a new timestamp
    printTimeStamps()
}

// fullPageRender will insert various bits of runtime info from JS and Go into the app
// TODO: this is virtually the same as pageRefresh - so just combine and add flag?
const fullPageRender = async() => {
    console.log('starting Go Herald instance and rendering the page')

    // load the Go Herald instance and populate the page data
    try {
        await loadRuntimeInfo()
    } catch (e) {
        printErrorMsg(e)
        return
    }

    // update the experiment drop down
    await updateExperimentDropDown()

    // print the pie chart
    await updatePieChart()

    // check the minKNOW status TODO: do this in Go via routine
    var minknowStatus = `${await checkAPIstatus()}`
    updateStatus('status_minknow', minknowStatus)

    // print the table
    await buildTable()

    // print a new timestamp
    printTimeStamps()
}