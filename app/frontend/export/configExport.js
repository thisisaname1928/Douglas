const modeToggleButton = document.getElementById('modeToggleButton');
const msg = document.getElementById('msg')
const configBox = document.getElementById('configBox')
const body = document.body;
tmp = window.location.href.split("/")
const UUID = tmp[tmp.length - 1]


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

document.addEventListener('DOMContentLoaded', () => {
    const savedMode = localStorage.getItem('theme');
    if (savedMode === 'dark') {
        body.classList.add('dark-mode');
        modeToggleButton.textContent = 'Chế độ sáng';
    } else {
        modeToggleButton.textContent = 'Chế độ tối';
    }

    initConf()
})

function renderStypeConfig(quesType, maxQues) {
    return `<div class="config-container">
                <label class="config-label">Cau loai ${quesType}, so cau ${maxQues};</label>
                <input id="${quesType}.N" type="text" class="config-input" placeholder="So cau cho mot de">
                <input id="${quesType}.Point" type="text" class="config-input" placeholder="So diem moi cau">
            </div>`
}

async function initConf() {
    res = await fetch4ExportConfig()
    obj = await res.json()

    console.log(obj)

    if (!obj.status) {
        msg.innerText = obj.msg
        return
    }


    for (i = 0; i < obj.stype.length; i++) {
        configBox.innerHTML += renderStypeConfig(obj.stype[i].stype, obj.stype[i].N)
    }
}

function fetch4ExportConfig() {
    res = fetch("/Export/API/getConfig", { method: "POST", body: JSON.stringify({ UUID: UUID }) })
    return res
}

