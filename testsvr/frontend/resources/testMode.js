const htmlElement = document.getElementsByTagName('html')[0]
TESTUUID = 'NONE'

function disableScroll() {
    document.body.style.overflow = 'hidden';

    document.documentElement.style.overflow = 'hidden';
}

async function warn() {
    res = await fetch("/api/warn", { method: "POST", body: JSON.stringify({ uuid: TESTUUID }) })
}

async function setUpTestMode(uuid) {
    TESTUUID = uuid

    htmlElement.requestFullscreen().catch((err) => {
        console.log(err)
        return false
    })

    document.addEventListener('fullscreenchange', (e) => {
        if (!document.fullscreenElement) {
            warn()
            showElement('popup-ask')
        }
    })

    document.addEventListener('keydown', (e) => {
        if (e.metaKey || e.ctrlKey) {
            if (e.key == 'v' || e.key == 'V' || e.key == 'C' || e.key == 'c') {
                e.preventDefault()
            }
        }
    })

    document.addEventListener('visibilitychange', (e) => {
        if (document.hidden) {
            warn()
        }
    })

    window.addEventListener('blur', (e) => {
        warn()
    });

    return true
}
