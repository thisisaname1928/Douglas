let uuid = "NONE"
const justAButton = document.getElementById("justAButton")

document.addEventListener('DOMContentLoaded', () => {
    uuid = window.location.pathname.replace("/StartTest.TestInfo/uuid/", "")

    updateTestInfo()
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

async function updateTestInfo() {
    testInfo = await getTestInfo()

    document.getElementById("testName").innerHTML = `${testInfo.name}`
    document.getElementById("candinate").innerHTML = `So luot lam bai: ${testInfo.numberOfCandinate}`
    document.getElementById("testUUID").innerHTML = `Ma de thi: ${testInfo.uuid}`

    if (testInfo.isStarted) {
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