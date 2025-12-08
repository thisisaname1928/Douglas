const updateStatus = document.getElementById('updateStatus')

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

        res = await update()

        if (!res.status) {
            alert(res.msg)
        }
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