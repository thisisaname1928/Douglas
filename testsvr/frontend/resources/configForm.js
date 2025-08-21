const startTestButton = document.getElementById("startTest")
const nameInput = document.getElementById("nameInput")
const classInput = document.getElementById("classInput")
const testContent = document.getElementById("testContent")
const configForm = document.getElementById("configForm")

startTestButton.addEventListener("click", () => {
    sendConfigForm(nameInput.value, classInput.value)
})

function chooseTNOption(i, ans) {
    if (ans > 3) {
        return;
    }
    const item = document.getElementById(`QUES.${i}.TN.${ans}`)
    if (item.classList.contains("option-item-highlighted")) { // no need to reset
        return;
    }


    item.classList.replace("option-item", "option-item-highlighted")

    for (j = 0; j < 4; j++) {
        if (j != ans) {
            document.getElementById(`QUES.${i}.TN.${j}`).classList.replace("option-item-highlighted", "option-item")
        }
    }
}

async function renderTest(testsvr) {
    questions = testsvr.test.questions

    // delete config form
    configForm.innerHTML = "";

    // render questions
    for (i = 0; i < questions.length; i++) {
        if (questions[i].type == 0x12) { // TN
            testContent.innerHTML += `
<div class="question-card">
    <div class="question-text">
        Câu ${i + 1} (Trac nghiem): ${questions[i].content}
    </div>
    <div class="options-list">
        <div class="option-item" id="QUES.${i}.TN.0" onclick="chooseTNOption(${i}, 0)">
            <div class="option-letter">A.</div>
            <div class="option-text">${questions[i].answers[0]}</div>
        </div>
        <div class="option-item" id="QUES.${i}.TN.1" onclick="chooseTNOption(${i}, 1)">
            <div class="option-letter">B.</div>
            <div class="option-text">${questions[i].answers[1]}</div>
        </div>
        <div class="option-item" id="QUES.${i}.TN.2" onclick="chooseTNOption(${i}, 2)">
            <div class="option-letter">C.</div>
            <div class="option-text">${questions[i].answers[2]}</div>
        </div>
        <div class="option-item" id="QUES.${i}.TN.3" onclick="chooseTNOption(${i}, 3)">
            <div class="option-letter">D.</div>
            <div class="option-text">${questions[i].answers[3]}</div>
        </div>
    </div>
</div>`
        }
    }
}

async function getTest(uuid) {
    res = await fetch("/api/getTest", { method: "POST", body: JSON.stringify({ uuid: uuid }) })
    jsonRes = await res.json()

    if (!jsonRes.status) {
        alert(jsonRes.msg)
        return
    }

    renderTest(jsonRes)
}

async function sendConfigForm(name, className) {
    res = await fetch("/api/startTest", { method: "POST", body: JSON.stringify({ name: name, className: className }) })
    jsonRes = await res.json()

    if (!jsonRes.status) {
        alert(jsonRes.msg)
    }

    getTest(jsonRes.uuid)
}