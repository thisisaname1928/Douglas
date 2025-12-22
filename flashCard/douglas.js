const debugArea = document.getElementById("debugArea")

function println(str) {
    debugArea.innerHTML += str + "<br>"
}

async function getFile(src) {
    return (await fetch(src))
}

function wordPathtoRID(path) {
    for (i = 0; i < path.length; i++) {
        if (path[i] == '\\')
            path[i] = '/'
    }

    return '/' + path
}


let rIDTable = new Map()

async function normalizeDouImgPath() {
    const imgTags = document.getElementsByTagName('img')

    for (i = 0; i < imgTags.length; i++) {
        currentTagSrc = imgTags[i].getAttribute("src")
        if (currentTagSrc.startsWith("media\\") || currentTagSrc.startsWith("media/")) {
            imgTags[i].setAttribute("src", rIDTable.get(wordPathtoRID(currentTagSrc)))
        }
    }
}


async function getDouNeccessaryFile(blobReader) {
    const zipReader = new zip.ZipReader(blobReader)

    entries = await zipReader.getEntries()

    data = ""
    info = ""

    for (i = 0; i < entries.length; i++) {
        if (!entries[i].directory) {
            if (entries[i].filename == "info.json") {
                info = entries[i].getData(new zip.TextWriter())
            } else if (entries[i].filename == "data.json") {
                data = entries[i].getData(new zip.TextWriter())
            } else if (entries[i].filename.startsWith("/media")) {
                const blob = await entries[i].getData(new zip.BlobWriter())

                // add to blob map
                url = URL.createObjectURL(blob)
                rIDTable.set(entries[i].filename, url)
            }
        }
    }

    return [info, data]
}

async function getQuestionFromDou(blobReader) {
    let res = await getDouNeccessaryFile(blobReader)
    let info = await res[0]
    let data = await res[1]
    let infoJson
    let dataJson

    try {
        infoJson = await JSON.parse(info)
        dataJson = await JSON.parse(data)
        console.log(infoJson)
        console.log(dataJson)
    } catch (e) {
        throw e
    }

    if (!infoJson.revision || !dataJson) {
        throw "ERR_NOT_VALID_FILE"
    }

    questionsRaw = []

    // merge all quesion

    for (i = 0; i < dataJson.questions.length; i++) {
        for (j = 0; j < dataJson.questions[i].questions.length; j++) {
            dataJson.questions[i].questions[j]
            questionsRaw.push(shuffleQuesAnsArray(dataJson.questions[i].questions[j]))
        }
    }

    questions = shuffleQuesArray(questionsRaw)

    return questions
}

function showElement(id) {
    document.getElementById(id).classList.remove("hidden-element")
}

function hideElement(id) {
    document.getElementById(id).classList.add("hidden-element")
}

let currentQues
let currentQuesIdx = 0
let n = 0
let curBlobReader
let numOfRightQues = 0
let curBlob

function checkAns() {

    if (!currentQues.type) {
        return false
    }

    switch (currentQues.type) {
        case 0x12:
            ans = getTNAnswer(0)

            if (currentQues.TNAnswers[ans]) {
                item = document.getElementById(`QUES.0.TN.${ans}`)
                item.innerHTML += "<i class=\"material-icons icon\">check</i>"
                numOfRightQues++
                setRightQuesNum(numOfRightQues)
            } else {
                item = document.getElementById(`QUES.0.TN.${ans}`)
                item.innerHTML += "<i class=\"material-icons icon\">close</i>"

                for (counter = 0; counter < 4; counter++) {
                    if (currentQues.TNAnswers[counter]) {
                        item = document.getElementById(`QUES.0.TN.${counter}`)
                        item.innerHTML += "<i class=\"material-icons icon\">check</i>"
                        break
                    }
                }
            }
            break
        case 0x13:
            break
        case 0x15:
            ans = getTNDSAnswer(0)
            nr = 0
            for (i = 0; i < 4; i++) {
                if (ans[i] == currentQues.TNAnswers[i]) {
                    nr++
                    item = document.getElementById(`QUES.0.TNDS.${i}`)
                    item.innerHTML += "<i class=\"material-icons icon\">check</i>"
                } else {
                    item = document.getElementById(`QUES.0.TNDS.${i}`)
                    item.innerHTML += "<i class=\"material-icons icon\">close</i>"
                }

            }

            if (nr == 4) numOfRightQues++
            break
    }

    hideElement("checkAnswer")
    showElement("nextQues")

    currentQuesIdx++

    setRightQuesNum(numOfRightQues)
}



function setQuesIdx(i, n) {
    document.getElementById("quesIdx").innerText = `Câu ${i}/${n}`
}

function setRightQuesNum(n) {
    document.getElementById("quesRightNum").innerText = `Đúng: ${n}`
}

function redoTest() {
    hideElement("endTest")

    doTest(curBlob)
}

async function doTest(blob) {
    document.getElementById("testContent").innerHTML = '<br>Chờ xíu, đang tải!'

    const blobReader = new zip.BlobReader(await blob)
    curBlob = blob
    curBlobReader = blobReader

    try {
        questions = await getQuestionFromDou(blobReader)
    }
    catch {
        setDoTestErr("Đọc đề lỗi, hãy chắc chắn rằng file đề đúng định dạng .dou")
        document.getElementById("testContent").innerHTML = ''
        return
    }

    hideElement('startTest')

    n = questions.length

    numOfRightQues = 0

    showElement("quesIdx")
    showElement("quesRightNum")
    showElement("checkAnswer")

    currentQuesIdx = 0

    await nextQuestion(currentQuesIdx)
}

async function nextQuestion(i) {

    if (currentQuesIdx >= n) {
        document.getElementById("testContent").innerHTML = ''
        hideElement("nextQues")

        showElement("endTest")
        document.getElementById("doneQuesNumL").innerHTML = n
        document.getElementById("rightQuesNumL").innerHTML = numOfRightQues

        return
    }

    renderTest([await questions[i]])
    setQuesIdx(i + 1, n)
    setRightQuesNum(numOfRightQues)
    currentQues = questions[i]

    normalizeDouImgPath()
}

async function renderTest(questions) {
    testContent = document.getElementById("testContent")
    testContent.innerHTML = ''


    // render questions
    for (i = 0; i < questions.length; i++) {
        if (questions[i].type == 0x12) { // TN
            testContent.innerHTML += `
<div class="question-card" id="QUES.${i}.TN">
    <div class="question-text">
        Câu ${currentQuesIdx + 1} (Trắc nghiệm): ${questions[i].content}
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
    <div class="question-card" id="QUES.${i}.TLN">
        <div class="question-text">
            Câu ${currentQuesIdx + 1} (Trắc nghiệm trả lời ngắn): ${questions[i].content}
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
<div class="question-card" id="QUES.${i}.TNDS">
    <div class="question-text" >
        Câu ${currentQuesIdx + 1} (Trắc nghiệm đúng sai): ${questions[i].content}
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

// just grab code from testsvr
async function chooseTNOption(i, ans, shouldUpdate) {
    // if (shouldUpdate) {
    //     result = await updateAnswerSheet(i, [String.fromCharCode('A'.charCodeAt(0) + ans), '', '', ''])

    //     if (!result) {
    //         return
    //     }
    // }

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
        rightAns.classList.replace("tnds-ans-r", "tnds-ans-r-highlighted")
        wrongAns.classList.replace("tnds-ans-w-highlighted", "tnds-ans-w")
    } else {
        wrongAns.classList.replace("tnds-ans-w", "tnds-ans-w-highlighted")
        rightAns.classList.replace("tnds-ans-r-highlighted", "tnds-ans-r")
    }
}

// fisher-yates shuffling
function shuffleQuesArray(array) {
    let currentIndex = array.length, randomIndex;
    while (currentIndex != 0) {

        randomIndex = Math.floor(Math.random() * currentIndex);
        currentIndex--;

        [array[currentIndex], array[randomIndex]] = [
            array[randomIndex], array[currentIndex]
        ];
    }

    return array;
}

// fisher-yates shuffling
function shuffleQuesAnsArray(array) {
    let currentIndex = array.answers.length, randomIndex;
    while (currentIndex != 0) {

        randomIndex = Math.floor(Math.random() * currentIndex);
        currentIndex--;

        [array.answers[currentIndex], array.answers[randomIndex]] = [
            array.answers[randomIndex], array.answers[currentIndex]
        ];
        [array.TNAnswers[currentIndex], array.TNAnswers[randomIndex]] = [
            array.TNAnswers[randomIndex], array.TNAnswers[currentIndex]
        ];

    }

    return array;
}

function getTNAnswer(th) {
    for (i = 0; i < 4; i++) {
        item = document.getElementById(`QUES.${th}.TN.${i}`)
        if (item.classList.contains("option-item-highlighted")) {
            return i
        }
    }

    return ''
}

function getTNDSAnswer(th) {
    // prevent smth bad
    ans = [false, false, false, false]
    for (i = 0; i < 4; i++) {
        item = document.getElementById(`QUES.${th}.TNDS.${i}.R`)
        if (item.classList.contains("tnds-ans-r-highlighted")) {
            ans[i] = true
        }

        item = document.getElementById(`QUES.${th}.TNDS.${i}.W`)
        if (item.classList.contains("tnds-ans-w-highlighted")) {
            ans[i] = false
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

    // remove space
    for (i = 0; i < 4; i++) {
        if (ans[0] == ' ') {
            ans[0] = ''
        }
    }

    return ans
}

document.getElementById("fileInput").addEventListener('change', (e) => {
    const f = e.target.files[0]

    if (!f) {
        return
    }

    const reader = new FileReader()

    reader.onload = async function (e) {
        const arrayBuffer = e.target.result;

        blob = new Blob([arrayBuffer])
        doTest(await blob)
    }

    reader.readAsArrayBuffer(f)
})

document.getElementById("fileInputBtn").addEventListener("click", () => {
    document.getElementById("fileInput").click()
})

function setDoTestErr(msg) {
    document.getElementById('doTestStatus').innerHTML = msg
}

document.addEventListener('DOMContentLoaded', async () => {
    const params = new URLSearchParams(window.location.search);

    link = params.get("useLink")

    if (!link || link == "") {
        return
    }

    // fetch 4 file
    try {
        request = await fetch(link)
        data = await request.arrayBuffer()

        blob = new Blob([data])
    }
    catch {
        setDoTestErr("Link truy cập tới bài làm bị hỏng!")
        return
    }

    doTest(blob)
})