const modeToggleButton = document.getElementById('modeToggleButton');
const msg = document.getElementById('msg')
const configBox = document.getElementById('configBox')
const exportButton = document.getElementById('exp')
const body = document.body;
tmp = window.location.href.split("/")
let UUID = tmp[tmp.length - 1]
const parsedUrl = new URL(location.href)
const parseExportType = parsedUrl.searchParams.get("exportType")


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

    if (!parseExportType) {
        configBox.innerHTML = 'Lỗi không xác định, vui lòng xuất lại đề!'
        return
    }

    res = ""
    for (i = 0; i < UUID.length; i++) {
        if (UUID[i] == '?') {
            break
        }

        res += UUID[i]
    }

    UUID = res

    exportButton.addEventListener('click', () => {
        expData = getConfig()
        if (!expData.status) {
            alert(expData.msg)
        }

        exportTest(expData)
    })

    initConf()
})

function renderStypeConfig(quesType, maxQues) {
    return `<div class="config-container">
                
                <label class="config-label">Câu loại ${quesType}, số câu ${maxQues}:</label>
                <input id="${quesType}.N" type="number" class="config-input" placeholder="Số câu mỗi đề">
                <input id="${quesType}.Point" type="number" class="config-input" placeholder="Số điểm mỗi câu">
                
            </div>`
}

let confObj = {};

async function initConf() {
    res = await fetch4ExportConfig()
    obj = await res.json()

    confObj = obj

    if (!obj.status) {
        if (obj.msg == "invalid docx file") {
            msg.innerText = "File được xuất không hợp lệ!"
        } else {
            msg.innerText = obj.msg
        }
        return
    }

    configBox.innerHTML += `<div class="config-container">
                <label class="config-label">Tác giả:</label>
                <input id="author" type="text" class="config-input" placeholder="Tên người ra đề">
            </div>`

    configBox.innerHTML += `<div class="config-container">
                <label class="config-label">Mật khẩu cho đề:</label>
                <input id="key" type="text" class="config-input" placeholder="Để trống nếu không dùng mã hóa">
            </div>`

    configBox.innerHTML += `<div class="config-container">
                <label class="config-label">Thời gian làm bài (phút):</label>
                <input id="testDuration" type="number" class="config-input" placeholder="Đặt 0 nếu không có thời gian cố định">
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

    if (key.length > 16) { return { status: false, UUID: UUID, testDuration: testDuration, msg: "Mật khẩu quá dài", author: author, key: key, stype: stype } }


    for (i = 0; i < confObj.stype.length; i++) {
        numberOfQuesPerTest = Math.abs(Number(document.getElementById(`${confObj.stype[i].stype}.N`).value))
        pointPerQues = Math.abs(Number(document.getElementById(`${confObj.stype[i].stype}.Point`).value))

        if (Number.isNaN(numberOfQuesPerTest) || Number.isNaN(pointPerQues)) {
            return { status: false, UUID: UUID, testDuration: testDuration, msg: `Dữ liệu nhập cho loại câu ${confObj.stype[i].stype} không hợp lệ!`, author: author, key: key, stype: stype }
        }

        if (numberOfQuesPerTest > confObj.stype[i].N) {
            return { status: false, UUID: UUID, testDuration: testDuration, msg: `Số câu mỗi đề của loại ${confObj.stype[i].stype} lớn hơn số câu tồn tại`, author: author, key: key, stype: stype }
        }

        stype.push({ stype: confObj.stype[i].stype, N: numberOfQuesPerTest, Point: pointPerQues })
    }

    return { status: true, UUID: UUID, testDuration: testDuration, msg: "ok", author: author, key: key, stype: stype, exportType: parseExportType }
}

async function exportTest(obj) {
    res = await fetch("/Export/API/export", { method: "POST", body: JSON.stringify(obj) })
    dat = res.json()
    if (dat.status == false) {
        alert(dat.msg)
        return
    }

    alert('Xuất đề thành công!')
    const link = document.createElement("a")
    link.href = "/Export/Download/UUID/" + UUID + ".dou"
    link.download = 'exported.dou'
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
}

function fetch4ExportConfig() {
    res = fetch("/Export/API/getConfig", { method: "POST", body: JSON.stringify({ UUID: UUID, exportType: parseExportType }) })
    return res
}

