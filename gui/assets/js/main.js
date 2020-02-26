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

// getProtobufDump returns a stringified protobuf dump of a sample from the database
const getProtobufDump = async function(sampleLabel) {
    var sampleProtobufDump = `${await window.printSampleToString(sampleLabel)}`
    return sampleProtobufDump
}

////////////////////////////////////////////////////////////////////
// BUTTONS
// get the buttons that control the app
const refreshPage = document.getElementById('refreshPage')
const wipeDatabase = document.getElementById('wipeDatabase')

// add an event listener to the refreshPage button
refreshPage.addEventListener('click', async() => {
    console.log('refreshing the app')

    fullPageRender()
    printSuccessMsg('refreshed the app')
})

// add an event listener to wipeDatabase button
wipeDatabase.addEventListener('click', async() => {
    console.log('wiping database')

    // TODO: add a confirm prompt

    try {
        await wipeStorage()
    } catch (e) {
        console.log(e)
        printErrorMsg(e)
        return
    }

    // reset the page and report success
    fullPageRender(true)
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

    // reset the form, re-render the page and report success
    addSampleForm.reset()
    fullPageRender()
    printSuccessMsg('sample added')
})

////////////////////////////////////////////////////////////////////
// TABLES
// clearTable will delete any table content already on the page
const clearTable = async() => {
    $('#sampleTable')
        .DataTable()
        .destroy()
    $('#sampleTableContent').empty()
}

// buildTable will get the database keys via Go and then populate the table
const buildTable = async() => {
    console.log('building table from the database')

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
        var data = table.row($(this).parents('tr')).data()
        var sampleLabel = data[0]
        document.getElementById('sampleModal_samplename').innerHTML = sampleLabel
        getProtobufDump(sampleLabel).then(sampleProtobufDump => {
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

                    try {
                        await deleteSample(sampleLabel)
                    } catch (e) {
                        console.log(e)
                        printErrorMsg(e)
                        return
                    }
                    //table.row($(this).parents('tr')).delete();

                    // reset the page and report success
                    fullPageRender(true)
                    document.getElementById('sampleModal').style.display = 'none'
                    printSuccessMsg('sample deleted')
                })
        })
    })

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

// printChart will print the pie chart
const printChart = async() => {
    // get total samples
    var sampleCount = `${await window.getSampleCount()}`

    // TODO: get total tagged samples
    //var taggedSamples = 1

    // TODO: get current announcements

    // tmp holder for the data (will get this from Go)
    var pieData = [{
            value: 300,
            color: '#35cebe',
            highlight: '#25beae',
            label: 'Active Announcements'
        },
        {
            value: 50,
            color: '#a0a0a0',
            highlight: '#999999',
            label: 'Tagged Samples'
        },
        {
            value: 220,
            color: '#dfdfdf',
            highlight: '#cccccc',
            label: 'Untagged Samples'
        }
    ]

    // print it
    var chart = document.getElementById('pie-chart').getContext('2d')
    window.myPie = new Chart(chart).Pie(pieData, {
        responsive: true,
        segmentShowStroke: false
    })
}

// fullPageRender will insert various bits of runtime info from JS and Go into the app
const fullPageRender = async tableFlag => {
    // clear and re-build the table if requested
    if (tableFlag === undefined) {
        tableFlag = false
    }
    if (tableFlag === true) {
        await clearTable()
        await buildTable()
    }

    // render the data from Go
    console.log('collecting data via Go and rendering the page')
    await window.renderPage()

    // print the pie chart
    printChart()

    // print a new timestamp
    printTimeStamps()
}