const testContent = document.getElementById("testContent")

window.addEventListener('load', function () {
    a = this.window.location.href.split("/")
    uuid = a[a.length - 1]
    getTest(uuid)
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

function chooseTNDSOption(i, ansI, value) {
    rightAns = document.getElementById(`QUES.${i}.TNDS.${ansI}.R`)
    wrongAns = document.getElementById(`QUES.${i}.TNDS.${ansI}.W`)

    if (value) {
        rightAns.classList.replace("tnds-ans-r", "tnds-ans-r-highlighted")
        wrongAns.classList.replace("tnds-ans-w-highlighted", "tnds-ans-w")
    } else {
        wrongAns.classList.replace("tnds-ans-w", "tnds-ans-w-highlighted")
        rightAns.classList.replace("tnds-ans-r-highlighted", "tnds-ans-r")
    }
}

let questions

async function renderTest(testsvr) {
    questions = testsvr.test.questions

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
        } else if (questions[i].type == 0x13) {
            testContent.innerHTML += `    
    <div class="question-card">
        <div class="question-text">
            Câu ${i + 1} (Trac nghiem tra loi ngan): ${questions[i].content}
        </div>
        <div class="TLN-input-container">
            <input class="square-input" type="text" maxlength="1" id="QUES.${i}.TLN.0">
            <input class="square-input" type="text" maxlength="1" id="QUES.${i}.TLN.1">
            <input class="square-input" type="text" maxlength="1" id="QUES.${i}.TLN.2">
            <input class="square-input" type="text" maxlength="1" id="QUES.${i}.TLN.3">
        </div>
    </div>`
        } else if (questions[i].type == 0x14) {
            testContent.innerHTML += `
<div class="question-card">
    <div class="question-text">
        Câu ${i + 1} (Trac nghiem dung sai): ${questions[i].content}
    </div>
    <div class="options-list">
        <div class="option-item" id="QUES.${i}.TNDS.0">
            <button class="tnds-ans-r" id="QUES.${i}.TNDS.0.R" onclick="chooseTNDSOption(${i}, 0, true)" >D</button>
            <button class="tnds-ans-w" id="QUES.${i}.TNDS.0.W" onclick="chooseTNDSOption(${i}, 0, false)">S</button>
            <div class="option-letter">a) </div>
            <div class="option-text">${questions[i].answers[0]}</div>
        </div>
        <div class="option-item" id="QUES.${i}.TNDS.1">
            <button class="tnds-ans-r" id="QUES.${i}.TNDS.1.R" onclick="chooseTNDSOption(${i}, 1, true)">D</button>
            <button class="tnds-ans-w" id="QUES.${i}.TNDS.1.W" onclick="chooseTNDSOption(${i}, 1, false)">S</button>
            <div class="option-letter">b) </div>
            <div class="option-text">${questions[i].answers[1]}</div>
        </div>
        <div class="option-item" id="QUES.${i}.TNDS.2">
            <button class="tnds-ans-r" id="QUES.${i}.TNDS.2.R" onclick="chooseTNDSOption(${i}, 2, true)">D</button>
            <button class="tnds-ans-w" id="QUES.${i}.TNDS.2.W" onclick="chooseTNDSOption(${i}, 2, false)">S</button>
            <div class="option-letter">c) </div>
            <div class="option-text">${questions[i].answers[2]}</div>
        </div>
        <div class="option-item" id="QUES.${i}.TNDS.3">
            <button class="tnds-ans-r" id="QUES.${i}.TNDS.3.R" onclick="chooseTNDSOption(${i}, 3, true)">D</button>
            <button class="tnds-ans-w" id="QUES.${i}.TNDS.3.W" onclick="chooseTNDSOption(${i}, 3, false)">S</button>
            <div class="option-letter">d) </div>
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

function getTNAnswer(th) {
    for (i = 0; i < 4; i++) {
        item = document.getElementById(`QUES.${th}.TN.${i}`)
        if (item.classList.contains("option-item-highlighted")) {
            if (i == 0) {
                return 'A'
            }

            if (i == 1) {
                return 'B'
            }

            if (i == 2) {
                return 'C'
            }

            if (i == 3) {
                return 'D'
            }
        }
    }

    return ''
}

function getTNDSAnswer(th) {
    // prevent smth bad
    ans = ['', '', '', '']
    for (i = 0; i < 4; i++) {
        item = document.getElementById(`QUES.${th}.TNDS.${i}.R`)
        if (item.classList.contains("tnds-ans-r-highlighted")) {
            ans[i] = 'r'
        }

        item = document.getElementById(`QUES.${th}.TNDS.${i}.W`)
        if (item.classList.contains("tnds-ans-w-highlighted")) {
            ans[i] = 'w'
        }
    }

    return ans
}

function getTLNAnswer(th) {
    ans = ['', '', '', '']

    ans[0] = document.getElementById(`QUES.${th}.TLN.0`).value
    ans[1] = document.getElementById(`QUES.${th}.TLN.1`).value
    ans[2] = document.getElementById(`QUES.${th}.TLN.2`).value
    ans[3] = document.getElementById(`QUES.${th}.TLN.3`).value

    return ans
}

function getAnswer() {
    result = []
    for (let i = 0; i < questions.length; i++) {
        if (questions[i].type == 0x12) {
            result.push(getTNAnswer(i))
        } else if (questions[i].type == 0x14) {
            result.push(getTNDSAnswer(i))
        } else if (questions[i].type == 0x13) {
            result.push(getTLNAnswer(i))
        }
    }

    return result
}