const modeToggleButton = document.getElementById('modeToggleButton');
const body = document.body;
const filePathInput = document.getElementById('filePathInput')
const msg = document.getElementById('msg')
const questions = document.getElementById('questions')

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

function fetchAPI(name, obj) {
    return fetch('/LivePreview/API/' + name, { method: 'POST', body: JSON.stringify(obj) })
}

function prepareQuestions(json) {
    if (!json.status) {
        msg.innerText = json.error
        return;
    }
    msg.innerText = ""

    ques = ""
    for (i = 0; i < json.questions.length; i++) {
        if (json.questions[i].type == 0x12) {
            ques += `<div class="question-card">
                    <div class="question-text">
                        Câu loại ${json.questions[i].stype} : ${json.questions[i].content}
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
                        <div>Đáp án: ${json.questions[i].TNAnswers}</div>
                    </div>
                </div>`
        } else {
            ques += `<div class="question-card">
                    <div class="question-text">
                        Câu loại ${json.questions[i].stype}: ${json.questions[i].content}
                    </div>
                    <div class="options-list">
                        <div class="option-item">
                            <div class="option-text">Đáp án: ${json.questions[i].TLNAnswers}</div>
                        </div>
                    </div>
                </div>`
        }
    }
    questions.innerHTML = ques;
}

function fetchForQuestions() {
    return fetch('/LivePreview/API/genJson', { method: 'POST', body: JSON.stringify({ path: filePathInput.value }) }).then((r) => { r.json().then((json) => prepareQuestions(json)) })
}

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

livePreviewLoop()
async function livePreviewLoop() {
    while (true) {
        await sleep(500);
        fetchForQuestions()
    }
}