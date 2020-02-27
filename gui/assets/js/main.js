////////////////////////////////////////////////////////////////////
// MESSAGES
// printErrorMsg
function printErrorMsg(msg) {
    console.log('error: ', msg)
    $('body').overhang({
        type: 'error',
        message: 'error: ' + msg,
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

////////////////////////////////////////////////////////////////////
// SAMPLE MODAL
// Get the required elements
const sampleModalClose = document.getElementsByClassName('modal-close')[0]
const sampleModal = document.getElementById('sampleModal')

// When the user clicks on <span> (x), close the modal
sampleModalClose.addEventListener('click', async() => {
    sampleModal.style.display = 'none'
})

// When the user clicks anywhere outside of the modal, close it
window.onclick = function(event) {
    if (event.target == sampleModal) {
        sampleModal.style.display = 'none'
    }
}

// getSampleJSONdump returns a stringified protobuf dump of a sample from the database
const getSampleJSONdump = async function(sampleLabel) {
    var sampleJSONdump = `${await window.printSampleToJSONstring(sampleLabel)}`
    return sampleJSONdump
}

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
    console.log('wiping database')

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
    printSuccessMsg('database wiped')
})

////////////////////////////////////////////////////////////////////
// FORMS
// get the forms
const addSampleForm = document.getElementById('addSampleForm')

// prevent default form action on the addSampleForm
function handleForm(event) {
    event.preventDefault()
}
addSampleForm.addEventListener('submit', handleForm)

// add an event listener to the addSampleForm
addSampleForm.addEventListener('submit', async() => {
    console.log('adding sample to database')

    var elements = addSampleForm.elements

    // grab the tags from the form
    var tags = []
    for (var i = 0, element;
        (element = elements[i++]);) {
        if (element.type === 'checkbox' && element.checked) {
            tags.push(element.value)
        }
    }

    // create a sample and add it to the database
    try {
        // TODO: try reading form straight into protobuf and then send a serialised stream to Go
        await createSample(
            elements['formLabel_sampleLabel'].value,
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
        .row.add([elements['formLabel_sampleLabel'].value])
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
        targets: 1,
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
        document.getElementById('sampleModal').style.display = 'block'

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
                document.getElementById('sampleModal').style.display = 'none'
                printSuccessMsg('sample deleted')
            })
    })
})

// buildTable will get the database keys via Go and then populate the table
const buildTable = async() => {
    console.log('building table from the database')

    // wipe any existing table
    table.clear().draw(true)

    // get the current sample number so that we can iterate over the samples
    var sampleCount = `${await window.getSampleCount()}`

    // process each sample label
    for (var i = 0; i < sampleCount; i++) {
        var sampleLabel = `${await window.getSampleLabel(i)}`

        // create the table entry
        table.row.add([sampleLabel]).draw(true)
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
    var retVal = `${await window.loadRuntimeInfo()}`
    if (retVal !== '') {
        printErrorMsg(retVal)
    }

    // update the pie chart
    await updatePieChart()

    // print a new timestamp
    printTimeStamps()
}

// fullPageRender will insert various bits of runtime info from JS and Go into the app
const fullPageRender = async() => {
    console.log('starting Go Herald instance and rendering the page')

    // load the Go Herald instance and populate the page data
    var retVal = `${await window.loadRuntimeInfo()}`
    if (retVal !== '') {
        printErrorMsg(retVal)
    }

    // print the pie chart
    await updatePieChart()

    // print the table
    await buildTable()

    // print a new timestamp
    printTimeStamps()
}