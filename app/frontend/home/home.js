const openLPButton = document.getElementById('openLP')
const openExButton = document.getElementById('openEx')
const openStButton = document.getElementById('openSt')
const openQEButton = document.getElementById('openQE')

openLPButton.addEventListener('click', () => {
    window.location.href = window.location.href.replace("/Home", "/LivePreview")
});

openStButton.addEventListener('click', () => {
    window.location.href = window.location.href.replace("/Home", "/StartTest")
});

openExButton.addEventListener('click', () => {
    window.location.href = window.location.href.replace("/Home", "/Export")
});

openQEButton.addEventListener('click', () => {
    window.location.href = window.location.href.replace("/Home", "/quickEditor")
});

async function getAppVersion() {
    res = await fetch("/getVersion", { method: "GET" })
    resObj = await res.json()

    return resObj
}

check4NewVersion()

async function check4NewVersion() {
    try {
        res = await fetch("https://raw.githubusercontent.com/thisisaname1928/Douglas/refs/heads/master/appVersion.json", { method: "GET" })
        resOb = await res.json()

        let appVersion = await getAppVersion()

        if (appVersion.versionInt != resOb.versionInt) {
            document.getElementById("versionNotify").innerHTML = `Đã có phiên bản mới hơn: Douglas_${resOb.versionStr}!`
        }
    }
    catch { }
}
