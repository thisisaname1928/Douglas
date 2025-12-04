const startTestButton = document.getElementById("startTest")
const nameInput = document.getElementById("nameInput")
const classInput = document.getElementById("classInput")
const configForm = document.getElementById("configForm")

startTestButton.addEventListener("click", () => {
    sendConfigForm(nameInput.value, classInput.value)
})

document.addEventListener("DOMContentLoaded", async () => {
    testName = await getTestName()
    document.getElementById("testNameHdr").innerText = `Làm bài kiểm tra "${testName}"`
})

function getTest(uuid) {
    window.location.href = "/taketest/" + uuid
}

async function sendConfigForm(name, className) {
    res = await fetch("/api/startTest", { method: "POST", body: JSON.stringify({ name: name, className: className }) })
    jsonRes = await res.json()

    if (!jsonRes.status) {
        alert(jsonRes.msg)
    }

    getTest(jsonRes.uuid)
}

async function getTestName() {
    res = await fetch("/api/getTestName", { method: "GET" })
    return await res.text()
}