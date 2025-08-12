const startTestButton = document.getElementById("startTest")
const nameInput = document.getElementById("nameInput")
const classInput = document.getElementById("classInput")

startTestButton.addEventListener("click", () => {
    sendConfigForm(nameInput.value, classInput.value)
})

async function sendConfigForm(name, className) {
    res = await fetch("/api/startTest", { method: "POST", body: JSON.stringify({ name: name, className: className }) })
    jsonRes = await res.json()

    if (!jsonRes.status) {
        alert(jsonRes.msg)
    }
}