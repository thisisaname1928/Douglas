const modeToggleButton = document.getElementById('modeToggleButton');
const body = document.body;
const filePathInput = document.getElementById('filePathInput')
const msg = document.getElementById('msg')
const questions = document.getElementById('questions')
const exportButton = document.getElementById('exportButton')

document.addEventListener('DOMContentLoaded', () => {
    const savedMode = localStorage.getItem('theme');
    if (savedMode === 'dark') {
        body.classList.add('dark-mode');
        modeToggleButton.textContent = 'Chế độ sáng';
    } else {
        modeToggleButton.textContent = 'Chế độ tối';
    }
})

modeToggleButton.addEventListener('click', () => {
    body.classList.toggle('dark-mode');

    if (body.classList.contains('dark-mode')) {
        modeToggleButton.textContent = 'Chế độ sáng';
        localStorage.setItem('theme', 'dark');
    } else {
        modeToggleButton.textContent = 'Chế độ tối';
        localStorage.setItem('theme', 'light');
    }
});

function transTNAnswer(A, B, C, D) {
    if (A) {
        return 'A'
    }

    if (B) {
        return 'B'
    }

    if (C) {
        return 'C'
    }

    if (D) {
        return 'D'
    }

    return 'Chưa xác định'
}

function transTNDSAnswer(A, B, C, D) {
    output = ""

    if (A) {
        output += 'Đ'
    } else {
        output += 'S'
    }

    if (B) {
        output += 'Đ'
    } else {
        output += 'S'
    }

    if (C) {
        output += 'Đ'
    } else {
        output += 'S'
    }

    if (D) {
        output += 'Đ'
    } else {
        output += 'S'
    }

    return output
}

function transTLNAnswer(a1, a2, a3, a4) {
    return a1 + a2 + a3 + a4
}

function fetchAPI(name, obj) {
    return fetch('/LivePreview/API/' + name, { method: 'POST', body: JSON.stringify(obj) })
}

function getTnAns(sheet) {
    for (i = 0; i < 4; i++) {
        if (sheet[i]) {
            return String.fromCharCode(65 + i)
        }
    }

    return "Chua co"
}

function prepareQuestions(json) {
    if (!json.status) {
        return;
    }

    ques = ""

    if (!json.questions || json.questions.length == 0) {
        questions.innerHTML = ''
        return
    }

    for (i = 0; i < json.questions.length; i++) {
        if (json.questions[i].type == 0x12) {
            ques += `<div class="question-card">
                    <div class="question-text">
                        Câu trắc nghiệm loại ${json.questions[i].stype} : ${json.questions[i].content}
                    </div>
                    <div class="options-list">
                        <div class="option-item">
                            <div class="option-letter">A.</div>
                            <div class="option-text">${json.questions[i].answers[0]}</div>
                        </div>
                        <div class="option-item">
                            <div class="option-letter">B.</div>
                            <div class="option-text">${json.questions[i].answers[1]}</div>
                        </div>
                        <div class="option-item">
                            <div class="option-letter">C.</div>
                            <div class="option-text">${json.questions[i].answers[2]}</div>
                        </div>
                        <div class="option-item">
                            <div class="option-letter">D.</div>
                            <div class="option-text">${json.questions[i].answers[3]}</div>
                        </div>
                        <div>Đáp án: ${transTNAnswer(json.questions[i].TNAnswers[0], json.questions[i].TNAnswers[1], json.questions[i].TNAnswers[2], json.questions[i].TNAnswers[3])}</div>
                    </div>
                </div > `
        } else if (json.questions[i].type == 0x13) {
            ques += `<div class="question-card">
                    <div class="question-text">
                        Câu trắc nghiệm trả lời ngắn loại ${json.questions[i].stype}: ${json.questions[i].content}
                    </div>
                    <div class="options-list">
                        <div class="option-item">
                            <div class="option-text">Đáp án: ${transTLNAnswer(json.questions[i].TLNAnswers[0], json.questions[i].TLNAnswers[1], json.questions[i].TLNAnswers[2], json.questions[i].TLNAnswers[3])}</div>
                        </div>
                    </div>
                </div> `
        } else if (json.questions[i].type == 0x15) {
            ques += `<div class="question-card">
                    <div class="question-text">
                        Câu trắc nghiệm đúng sai loại ${json.questions[i].stype} : ${json.questions[i].content}
                    </div>
                    <div class="options-list">
                        <div class="option-item">
                            <div class="option-letter">A.</div>
                            <div class="option-text">${json.questions[i].answers[0]}</div>
                        </div>
                        <div class="option-item">
                            <div class="option-letter">B.</div>
                            <div class="option-text">${json.questions[i].answers[1]}</div>
                        </div>
                        <div class="option-item">
                            <div class="option-letter">C.</div>
                            <div class="option-text">${json.questions[i].answers[2]}</div>
                        </div>
                        <div class="option-item">
                            <div class="option-letter">D.</div>
                            <div class="option-text">${json.questions[i].answers[3]}</div>
                        </div>
                        <div>Đáp án: ${transTNDSAnswer(json.questions[i].TNAnswers[0], json.questions[i].TNAnswers[1], json.questions[i].TNAnswers[2], json.questions[i].TNAnswers[3])}</div>
                    </div>
                </div> `}
    }
    questions.innerHTML = ques;
}

function fetchForQuestions() {
    return fetch('/LivePreview/API/genJson', { method: 'POST', body: JSON.stringify({ path: filePathInput.value }) }).then((r) => { r.json().then((json) => prepareQuestions(json)) })
}

async function fetch4Questions() {
    textContent = document.getElementById("textEditor").value

    res = await fetch("/API/quickPreview", { method: "POST", body: JSON.stringify({ content: textContent }) })
    obj = await res.json()

    prepareQuestions(obj)
}

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

async function livePreviewLoop() {
    while (true) {
        await sleep(700);
        fetch4Questions()
    }
}

function uploadFile(data, uuid) {
    var hdrs = new Headers()
    hdrs.append("uuid", uuid)
    hdrs.append("Content-Type", "application/octet-stream")
    return fetch("/Export/API/upload", { method: "POST", headers: hdrs, body: data })
}

async function getUUID() {
    response = await fetch("/Export/API/genUUID", { method: "POST" })
    res = await response.text()

    return await res
}


exportButton.addEventListener('click', async () => {
    textContent = document.getElementById('textEditor').value
    uuid = await getUUID()

    await uploadFile(textContent, uuid)

    window.location.href = window.location.href.replace("/quickEditor", "/Export/Config/UUID/" + uuid + "?exportType=useRawText")
})

document.getElementById('textEditor').addEventListener('input', () => {
    fetch4Questions()
})