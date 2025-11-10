const testContent = document.getElementById("testContent")
const summitBtn = document.getElementById("summitTest")
const testResult = document.getElementById("testResult")
let isDone = false

async function checkIfTestDone() {
    res = await fetch("/api/getTestStatus", { method: "POST", body: JSON.stringify({ uuid: uuid }) })
    jsonRes = await res.json()

    return jsonRes.status
}

window.addEventListener('load', async function () {
    a = this.window.location.href.split("/")
    uuid = a[a.length - 1]

    isDone = await checkIfTestDone()
    await getTest(uuid)
    if (!isDone) {
        setUpTimer(glbTestsvr.test.startTime, glbTestsvr.test.testDuration)
    } else {
        res = await fetch("/api/getPoint", { method: "POST", body: JSON.stringify({ uuid: uuid }) })
        jsonRes = await res.json()

        document.getElementById("test").innerHTML = JSON.stringify(jsonRes)

        bgTime = new Date(glbTestsvr.test.startTime)
        edTime = new Date(glbTestsvr.test.endTime)
        duration = edTime - bgTime

        testResult.innerHTML = `<h1>Diem: ${jsonRes.point}</h1><p>So cau dung hoan toan: ${jsonRes.trueQuesCount}/${questions.length}<br><br>Ten hoc sinh: ${glbTestsvr.test.name}<br>Lop: ${glbTestsvr.test.class}
        <br><br>
        Thoi gian lam bai: ${Math.trunc(duration / 60000)} phut ${Math.round(((duration % 60000) / 1000))} giay
        <br>Thoi gian bat dau: ${bgTime.toLocaleTimeString("en-us", {
            hour: '2-digit',
            minute: '2-digit', hour12: false
        })}<br>Thoi gian ket thuc: ${edTime.toLocaleTimeString("en-us", {
            hour: '2-digit',
            minute: '2-digit', hour12: false
        })}<br><br>Ma bai lam: ${uuid}</p>`

        setEndTime()
    }



    // load question sheet data
    loadUpAnsSheet()
})

async function chooseTNOption(i, ans, shouldUpdate) {
    // check if update client side or both client and server
    if (shouldUpdate) {
        result = await updateAnswerSheet(i, [String.fromCharCode('A'.charCodeAt(0) + ans), '', '', ''])
        //result = await updateAnswerSheet2(i, ans, String.fromCharCode('A'.charCodeAt(0) + ans))

        if (!result) {
            return
        }
    }

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

async function chooseTNDSOption(i, ansI, value, shouldUpdate) {
    rightAns = document.getElementById(`QUES.${i}.TNDS.${ansI}.R`)
    wrongAns = document.getElementById(`QUES.${i}.TNDS.${ansI}.W`)

    if (value) {

        if (shouldUpdate) {
            result = await updateAnswerSheet2(i, ansI, "T")

            if (!result) {
                return
            }
        }

        rightAns.classList.replace("tnds-ans-r", "tnds-ans-r-highlighted")
        wrongAns.classList.replace("tnds-ans-w-highlighted", "tnds-ans-w")
    } else {
        if (shouldUpdate) {
            result = await updateAnswerSheet2(i, ansI, "F")

            if (!result) {
                return
            }
        }

        wrongAns.classList.replace("tnds-ans-w", "tnds-ans-w-highlighted")
        rightAns.classList.replace("tnds-ans-r-highlighted", "tnds-ans-r")
    }
}

async function chooseTLNAnswer(index, answerIndex) {
    inp = document.getElementById(`QUES.${index}.TLN.${answerIndex}`)

    data = inp.value
    // add space to it
    if (data == "") {
        data = " "
    }

    result = await updateAnswerSheet2(index, answerIndex, data)
    if (!result) {
        // reset value
        inp.value = ""
    }
}

let questions
let glbTestsvr

async function renderTest(testsvr) {
    questions = testsvr.test.questions
    glbTestsvr = testsvr

    // render questions
    for (i = 0; i < questions.length; i++) {
        if (questions[i].type == 0x12) { // TN
            testContent.innerHTML += `
<div class="question-card">
    <div class="question-text">
        Câu ${i + 1} (Trac nghiem): ${questions[i].content}
    </div>
    <div class="options-list">
        <div class="option-item" id="QUES.${i}.TN.0" onclick="chooseTNOption(${i}, 0, true)">
            <div class="option-letter">A.</div>
            <div class="option-text">${questions[i].answers[0]}</div>
        </div>
        <div class="option-item" id="QUES.${i}.TN.1" onclick="chooseTNOption(${i}, 1, true)">
            <div class="option-letter">B.</div>
            <div class="option-text">${questions[i].answers[1]}</div>
        </div>
        <div class="option-item" id="QUES.${i}.TN.2" onclick="chooseTNOption(${i}, 2, true)">
            <div class="option-letter">C.</div>
            <div class="option-text">${questions[i].answers[2]}</div>
        </div>
        <div class="option-item" id="QUES.${i}.TN.3" onclick="chooseTNOption(${i}, 3, true)">
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
            <input class="square-input" type="text" maxlength="1" oninput="chooseTLNAnswer(${i}, 0)" id="QUES.${i}.TLN.0">
            <input class="square-input" type="text" maxlength="1" oninput="chooseTLNAnswer(${i}, 1)" id="QUES.${i}.TLN.1">
            <input class="square-input" type="text" maxlength="1" oninput="chooseTLNAnswer(${i}, 2)" id="QUES.${i}.TLN.2">
            <input class="square-input" type="text" maxlength="1" oninput="chooseTLNAnswer(${i}, 3)" id="QUES.${i}.TLN.3">
        </div>
    </div>`
        } else if (questions[i].type == 0x15) {
            testContent.innerHTML += `
<div class="question-card">
    <div class="question-text">
        Câu ${i + 1} (Trac nghiem dung sai): ${questions[i].content}
    </div>
    <div class="options-list">
        <div class="option-item" id="QUES.${i}.TNDS.0">
            <button class="tnds-ans-r" id="QUES.${i}.TNDS.0.R" onclick="chooseTNDSOption(${i}, 0, true, true)" >D</button>
            <button class="tnds-ans-w" id="QUES.${i}.TNDS.0.W" onclick="chooseTNDSOption(${i}, 0, false, true)">S</button>
            <div class="option-letter">a) </div>
            <div class="option-text">${questions[i].answers[0]}</div>
        </div>
        <div class="option-item" id="QUES.${i}.TNDS.1">
            <button class="tnds-ans-r" id="QUES.${i}.TNDS.1.R" onclick="chooseTNDSOption(${i}, 1, true, true)">D</button>
            <button class="tnds-ans-w" id="QUES.${i}.TNDS.1.W" onclick="chooseTNDSOption(${i}, 1, false, true)">S</button>
            <div class="option-letter">b) </div>
            <div class="option-text">${questions[i].answers[1]}</div>
        </div>
        <div class="option-item" id="QUES.${i}.TNDS.2">
            <button class="tnds-ans-r" id="QUES.${i}.TNDS.2.R" onclick="chooseTNDSOption(${i}, 2, true, true)">D</button>
            <button class="tnds-ans-w" id="QUES.${i}.TNDS.2.W" onclick="chooseTNDSOption(${i}, 2, false, true)">S</button>
            <div class="option-letter">c) </div>
            <div class="option-text">${questions[i].answers[2]}</div>
        </div>
        <div class="option-item" id="QUES.${i}.TNDS.3">
            <button class="tnds-ans-r" id="QUES.${i}.TNDS.3.R" onclick="chooseTNDSOption(${i}, 3, true, true)">D</button>
            <button class="tnds-ans-w" id="QUES.${i}.TNDS.3.W" onclick="chooseTNDSOption(${i}, 3, false, true)">S</button>
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

    await renderTest(jsonRes)

    summitBtn.addEventListener("click", function () { doneTest(); location.reload() })
}

// load reload answer sheet
async function loadUpAnsSheet() {
    res = await fetch("/api/getCurrentAnsSheet", { method: "POST", body: JSON.stringify({ UUID: uuid }) })
    jsonRes = await res.json()

    if (!jsonRes.status) {
        alert(jsonRes.msg)
        return false
    }

    for (i = 0; i < questions.length; i++) {
        // TN

        if (questions[i].type == 18) {
            if (jsonRes.ansSheet[i][0] != "") {
                chooseTNOption(i, jsonRes.ansSheet[i][0].charCodeAt(0) - 'A'[0].charCodeAt(0), false)
            }
        } else if (questions[i].type == 21) {
            for (j = 0; j < 4; j++) {
                if (jsonRes.ansSheet[i][j] != "") {
                    chooseTNDSOption(i, j, jsonRes.ansSheet[i][j] == "T", false)
                }
            }
        } else if (questions[i].type == 19) {
            for (j = 0; j < 4; j++) {
                if (jsonRes.ansSheet[i][j] != "") {
                    // manually copy
                    inp = document.getElementById(`QUES.${i}.TLN.${j}`)
                    inp.value = jsonRes.ansSheet[i][j]
                }
            }
        }
    }
}

async function doneTest() {
    res = await fetch("/api/handleDoneTest", { method: "POST", body: JSON.stringify({ UUID: uuid }) })
    jsonRes = await res.json()

    if (!jsonRes.status) {
        alert(jsonRes.msg)
        return false
    }

    return true
}

// update single answer in answerSheet array
async function updateAnswerSheet2(index, answerIndex, data) {
    res = await fetch("/api/updateAnswer", { method: "POST", body: JSON.stringify({ UUID: uuid, index: index, answerIndex: answerIndex, data: data, shouldClear: false }) })
    jsonRes = await res.json()

    if (!jsonRes.status) {
        if (jsonRes.msg == "OUT_OF_TIME") {
            location.reload()
            return
        }
        if (jsonRes.msg != "TEST_SESSION_LOCKED") { alert(jsonRes.msg) }
        return false
    }

    return true
}

async function updateAnswerSheet(i, answers) {
    // add data field just to patch
    res = await fetch("/api/updateAnswer", { method: "POST", body: JSON.stringify({ UUID: uuid, index: i, answerSheet: answers, data: "?", shouldClear: true }) })
    jsonRes = await res.json()

    if (!jsonRes.status) {
        if (jsonRes.msg == "OUT_OF_TIME") {
            location.reload()
            return
        }
        if (jsonRes.msg != "TEST_SESSION_LOCKED") { alert(jsonRes.msg) }
        return false
    }

    return true
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
        } else if (questions[i].type == 0x15) {
            result.push(getTNDSAnswer(i))
        } else if (questions[i].type == 0x13) {
            result.push(getTLNAnswer(i))
        }
    }

    return result
}