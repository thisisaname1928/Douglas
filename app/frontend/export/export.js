const modeToggleButton = document.getElementById('modeToggleButton');
const uploadTest = document.getElementById('uploadTest');
const fileInput = document.getElementById('fileInput');
const body = document.body;
let currentUUID = "NONE"

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

async function getUUID() {
    response = await fetch("/Export/API/genUUID", { method: "POST" })
    res = await response.text()

    return await res
}



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
        UUID = await uploadWrapper(data)
        window.location.href = window.location.href.replace("/Export", "/Export/Config/UUID/" + UUID + "?exportType=useDocx")

    } catch (error) {
        return;
    }
}

async function uploadWrapper(data) {
    curUUID = "NONE"
    curUUID = await getUUID()
    console.log(curUUID)
    await uploadFile(data, curUUID)

    return curUUID
}

function uploadFile(data, uuid) {
    var hdrs = new Headers()
    hdrs.append("uuid", uuid)
    hdrs.append("Content-Type", "application/octet-stream")
    return fetch("/Export/API/upload", { method: "POST", headers: hdrs, body: data })
}

uploadTest.addEventListener('click', () => {
    fileInput.click()
});

fileInput.addEventListener('change', () => {
    handleFile(fileInput.files[0])
});
