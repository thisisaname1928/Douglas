const updateStatus = document.getElementById('updateStatus')
const downloadStatus = document.getElementById('downloadStatus')
const statusMsg = document.getElementById('status')

document.addEventListener('DOMContentLoaded', async () => {
    updateStatus.innerHTML = `<i
                    class="material-icons icon material-spinner">refresh</i>Đang kiểm tra bảng cập nhật`

    res = await check4NewVersion()

    if (res == "") {
        updateStatus.innerHTML = `<i
                    class="material-icons icon">check</i>Đã là bản cập nhật mới nhất`
    } else {
        updateStatus.innerHTML = `<i
                    class="material-icons icon">check</i>Có bản cập nhật mới: Douglas_${res}`

        downloadStatus.innerHTML = `<i
                    class="material-icons icon material-spinner">refresh</i>Đang tải xuống bản cập nhật`

        res = await update()

        downloadStatus.innerHTML = `<i
                    class="material-icons icon">check</i>Đang tải xuống bản cập nhật`

        if (!res.status) {
            statusMsg.innerHTML = `<i
                    class="material-icons icon">close</i>Cập nhật thất bại: ${res.msg}`
            return
        }

        statusMsg.innerHTML = `<i
                    class="material-icons icon">check</i>Đã xong`
    }
})

async function getAppVersion() {
    res = await fetch("/getVersion", { method: "GET" })
    resObj = await res.json()

    return resObj
}

async function update() {
    res = await fetch("/downloadUpdate")
    return await res.json()
}

async function check4NewVersion() {
    try {
        res = await fetch("https://raw.githubusercontent.com/thisisaname1928/Douglas/refs/heads/master/appVersion.json", { method: "GET" })
        resOb = await res.json()

        appVersion = await getAppVersion()

        if (appVersion.versionInt != resOb.versionInt) {
            return resOb.versionStr
        } else {
            return ""
        }
    }
    catch {
        return ""
    }
}