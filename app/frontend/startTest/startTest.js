const modeToggleButton = document.getElementById('modeToggleButton');
const body = document.body;
const openUploadButton = document.getElementById('uploadFile')
const fileInput = document.getElementById('fileInput');
const configPopUp = document.getElementById('popup')
const blackLayer = document.getElementById('blackLayer')
const doneButton = document.getElementById('doneButton')

async function hideConfigPopup() {
    configPopUp.classList.add("hidden-element")

    blackLayer.classList.add("hidden-element")
}

async function showConfigPopup() {
    configPopUp.classList.remove("hidden-element")
    blackLayer.classList.remove("hidden-element")
}

//showConfigPopup()

openUploadButton.addEventListener('click', async () => {
    fileInput.click()
});

fileInput.addEventListener('change', async () => {
    handleFile(fileInput.files[0])
    showConfigPopup()
});

function readFile(file) {
    return new Promise((resolve, reject) => {
        const reader = new FileReader();

        reader.onload = (event) => {
            resolve(event.target.result);
        };

        reader.onerror = (error) => {
            reject(error);
        };

        reader.readAsArrayBuffer(file);
    });
}

async function handleFile(file) {
    if (!file) {
        return;
    }

    try {
        data = await readFile(file)
        await uploadFile(data)
    } catch (error) {
        return;
    }
}

async function loadTest() {
    res = await fetch("/StartTest/API/load", { method: "POST", body: JSON.stringify({ name: document.getElementById("name").value, key: document.getElementById("key").value }) })
    jsonRes = await res.json()

    if (!jsonRes.status) {
        alert(jsonRes.msg)
        return false
    }

    return true
}

async function uploadFile(data) {
    var hdrs = new Headers()
    hdrs.append("Content-Type", "application/octet-stream")
    return await fetch("/StartTest/API/upload", { method: "POST", headers: hdrs, body: data })
}


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


doneButton.addEventListener('click', async () => {
    await loadTest()
    hideConfigPopup()
})

blackLayer.addEventListener('click', () => {
    // only hide when 
    if (!blackLayer.classList.contains("hidden-element"))
        hideConfigPopup()
})

function delay(ms) {
    return new Promise(resolve => {
        setTimeout(resolve, ms);
    });
}

