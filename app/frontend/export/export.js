const modeToggleButton = document.getElementById('modeToggleButton');
const uploadTest = document.getElementById('uploadTest');
const fileInput = document.getElementById('fileInput');
const body = document.body;
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
    const response = await fetch("/Export/API/genUUID", { method: "POST" })
    res = await response.text()

    return res
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
        console.error("No file selected.");
        return;
    }

    try {
        const arrayBuffer = await readFile(file);

        uploadWrapper(arrayBuffer)

    } catch (error) {
        console.error("Error reading file:", error);
    }
}

function uploadWrapper(data) {
    // get unique uuid for test that being uploaded
    getUUID().then((res) => uploadFile(data, res))
}

function uploadFile(data, uuid) {
    console.log(uuid)
    var hdrs = new Headers()
    hdrs.append("uuid", uuid)
    hdrs.append("Content-Type", "application/octet-stream")
    fetch("/Export/API/upload", { method: "POST", headers: hdrs, body: data })
}

uploadTest.addEventListener('click', () => {
    fileInput.click()
});

fileInput.addEventListener('change', () => {
    handleFile(fileInput.files[0])
});
