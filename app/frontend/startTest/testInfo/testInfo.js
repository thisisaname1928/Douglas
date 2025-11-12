let uuid = "NONE"


document.addEventListener('DOMContentLoaded', () => {
    uuid = window.location.pathname.replace("/StartTest.TestInfo/uuid/", "")

    updateTestInfo()
})

async function getTestInfo() {
    res = await fetch("/StartTest/API/getTestInfo", { method: 'POST', body: JSON.stringify({ uuid: uuid }) })
    return await res.json()
}

async function updateTestInfo() {
    testInfo = await getTestInfo()

    document.getElementById("testName").innerHTML = `${testInfo.name}`
    document.getElementById("candinate").innerHTML = `So luot lam bai: ${testInfo.numberOfCandinate}`
    document.getElementById("testUUID").innerHTML = `Ma de thi: ${testInfo.uuid}`
}