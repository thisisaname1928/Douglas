const modeToggleButton = document.getElementById('modeToggleButton');
const msg = document.getElementById('msg')
const configBox = document.getElementById('configBox')
const exportButton = document.getElementById('exp')
const body = document.body;
tmp = window.location.href.split("/")
const UUID = tmp[tmp.length - 1]

exportButton.addEventListener('click', () => {
    expData = getConfig()
    if (!expData.status) {
        alert(expData.msg)
    }

    exportTest(expData)
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

let confObj = {};

async function initConf() {
    res = await fetch4ExportConfig()
    obj = await res.json()

    confObj = obj

    if (!obj.status) {
        msg.innerText = obj.msg
        return
    }

    configBox.innerHTML += `<div class="config-container">
                <label class="config-label">Tac gia:</label>
                <input id="author" type="text" class="config-input" placeholder="Ten nguoi ra de">
            </div>`

    configBox.innerHTML += `<div class="config-container">
                <label class="config-label">Mat khau cho de:</label>
                <input id="key" type="text" class="config-input" placeholder="De trong neu khong dung ma hoa">
            </div>`

    configBox.innerHTML += `<div class="config-container">
                <label class="config-label">Thoi gian lam bai (giay):</label>
                <input id="testDuration" type="text" class="config-input" placeholder="Dat 0 neu khong co thoi gian co dinh">
            </div>`

    for (i = 0; i < obj.stype.length; i++) {
        configBox.innerHTML += renderStypeConfig(obj.stype[i].stype, obj.stype[i].N)
    }
}

function getConfig() {
    author = document.getElementById("author").value
    key = document.getElementById("key").value
    testDuration = Math.abs(Number(document.getElementById("testDuration").value)) // abs to prevent user do somethings wrong=)))
    if (Number.isNaN(testDuration)) testDuration = 0;
    stype = []

    if (key.length > 16) { return { status: false, UUID: UUID, testDuration: testDuration, msg: "key qua dai", author: author, key: key, stype: stype } }


    for (i = 0; i < confObj.stype.length; i++) {
        numberOfQuesPerTest = Math.abs(Number(document.getElementById(`${confObj.stype[i].stype}.N`).value))
        pointPerQues = Math.abs(Number(document.getElementById(`${confObj.stype[i].stype}.Point`).value))

        if (Number.isNaN(numberOfQuesPerTest) || Number.isNaN(pointPerQues)) {
            return { status: false, UUID: UUID, testDuration: testDuration, msg: `du lieu nhap cho loai cau ${confObj.stype[i].stype} khong hop le!`, author: author, key: key, stype: stype }
        }

        if (numberOfQuesPerTest > confObj.stype[i].N) {
            return { status: false, UUID: UUID, testDuration: testDuration, msg: `so cau moi de cua loai ${confObj.stype[i].stype} lon hon so cau ton tai`, author: author, key: key, stype: stype }
        }

        stype.push({ stype: confObj.stype[i].stype, N: numberOfQuesPerTest, Point: pointPerQues })
    }

    return { status: true, UUID: UUID, testDuration: testDuration, msg: "ok", author: author, key: key, stype: stype }
}

async function exportTest(obj) {
    res = await fetch("/Export/API/export", { method: "POST", body: JSON.stringify(obj) })
    dat = res.json()
    if (dat.status == false) {
        alert(dat.msg)
        return
    }

    alert('Xuat de thanh cong!')
    const link = document.createElement("a")
    link.href = "/Export/Download/UUID/" + UUID
    link.download = 'exported.dou'
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
}

function fetch4ExportConfig() {
    res = fetch("/Export/API/getConfig", { method: "POST", body: JSON.stringify({ UUID: UUID }) })
    return res
}

