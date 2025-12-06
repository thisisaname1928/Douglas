let startTestTime
let maxDuration = 0
timer = document.getElementById("timer")

async function getCurrentServerTime() {
    res = await fetch("/api/getCurrentServerTime", { method: 'POST', body: { uuid: uuid } })

    return await res.text()
}

function delay(ms) {
    return new Promise(resolve => {
        setTimeout(resolve, ms);
    });
}

function setEndTime() { timer.innerHTML = `Đã hết thời gian` }

async function timerLoop() {
    duration = 0
    while (true) {
        ctStr = await getCurrentServerTime()

        if (await checkIfTestDone()) {
            timer.innerHTML = `Đã hết thời gian`
            if (!isDone) {
                await doneTest()
                location.reload()
            }
            continue
        }

        t = new Date(ctStr)
        duration = t - startTestTime
        console.log(duration)

        durationSec = maxDuration * 60 - Math.round(duration / 1000)

        if (durationSec < 0) {
            timer.innerHTML = `Đã hết thời gian`
            if (!isDone) {
                await doneTest()
                location.reload()
            }
            continue
        }

        currentTestTimeMin = Math.trunc(durationSec / 60)
        currentTestTimeSec = Math.round(((durationSec % 60)))

        timer.innerHTML = `Thời gian còn lại: ${String(currentTestTimeMin).padStart(2, '0')}:${String(currentTestTimeSec).padStart(2, '0')}`

        // set timeout
        await delay(1010)
    }
}

async function setUpTimer(timeStr, testDuration) {
    startTestTime = new Date(timeStr)

    maxDuration = testDuration

    if (maxDuration != 0)
        timerLoop()
}