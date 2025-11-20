let uuid = "NONE"
const justAButton = document.getElementById("justAButton")
const exportCsv = document.getElementById("exportCsv")

document.addEventListener('DOMContentLoaded', () => {
    uuid = window.location.pathname.replace("/StartTest.TestInfo/uuid/", "")

    updateTestInfo()
})

exportCsv.addEventListener('click', async () => {
    csvDat = await exportFileCSV()

    const fileContent = csvDat.csvData;
    const filename = "diem_thi.csv";
    const mimeType = "text/csv";

    const blob = new Blob([fileContent], { type: mimeType });

    const downloadLink = document.createElement("a")
    downloadLink.download = filename
    downloadLink.href = window.URL.createObjectURL(blob);

    downloadLink.click()
    window.URL.revokeObjectURL(downloadLink.href);
})


async function getTestInfo() {
    res = await fetch("/StartTest/API/getTestInfo", { method: 'POST', body: JSON.stringify({ uuid: uuid }) })
    return await res.json()
}

async function startATest() {
    res = await fetch("/StartTest/API/startATest", { method: 'POST', body: JSON.stringify({ uuid: uuid }) })
    return await res.json()
}

async function stopATest() {
    res = await fetch("/StartTest/API/stopATest", { method: 'POST', body: JSON.stringify({ uuid: uuid }) })
    return await res.json()
}

async function getTestIp() {
    res = await fetch("/StartTest/API/getTestIp", { method: 'POST', body: JSON.stringify({ uuid: uuid }) })

    return await res.text()
}
async function exportFileCSV() {
    res = await fetch("/StartTest/API/exportCsv", { method: 'POST', body: JSON.stringify({ uuid: uuid }) })

    return await res.json()
}


async function getCandinateList() {
    res = await fetch("/StartTest/API/getCandinateList", { method: 'POST', body: JSON.stringify({ uuid: uuid }) })
    resJ = await res.json()
    return await resJ.candinates
}

function checkCanMark(mark, done) {
    if (done) { return mark } else { return "chua nop bai" }
}

async function updateTestCandinate() {
    candinates = await getCandinateList()

    candinateBox = document.getElementById("candinateBox")

    for (i = 0; i < candinates.length; i++) {
        candinateBox.innerHTML += `<div class="test-card" style="cursor: pointer;" onclick='viewTest("${candinates[i].uuid}")'>
                    <b>Ten: ${candinates[i].name}, Lop: ${candinates[i].class}</b><br>
                    Diem: ${checkCanMark(candinates[i].mark, candinates[i].isDone)}
                </div>`
    }
}

async function updateTestInfo() {
    testInfo = await getTestInfo()

    document.getElementById("testName").innerHTML = `${testInfo.name}`
    document.getElementById("candinate").innerHTML = `So luot lam bai: ${testInfo.numberOfCandinate}`
    document.getElementById("testUUID").innerHTML = `Ma de thi: ${testInfo.uuid}`

    updateTestCandinate()

    if (testInfo.isStarted) {
        ip = await getTestIp()
        document.getElementById("testName").innerHTML += ` <a style="text-decoration:none ;color: lightgreen; cursor: pointer;" onclick="window.open('http:///${ip}', '_blank');">${ip}</a>`
        document.getElementById("testName").style = "color: lightgreen;"
        document.getElementById("testStatus").innerHTML = "Trang thai bai kiem tra: dang duoc mo"
        justAButton.innerHTML = "Dung bai kiem tra"

        justAButton.addEventListener('click', async () => {
            stopATest()
            location.reload()
        })
    } else {
        document.getElementById("testStatus").innerHTML = "Trang thai bai kiem tra: chua duoc mo"
        justAButton.addEventListener('click', async () => {
            startATest()
            location.reload()
        })
    }
}