const AIModeButton = document.getElementById('AIModeButton')
const createAITest = document.getElementById('createAITest')
const blackLayer = document.getElementById('blackLayer')
const closeButton = document.getElementById('closeBtn')

AIModeButton.addEventListener('click', () => { AIMode() })
createAITest.addEventListener('click', async () => {
    numberOfTNQues = parseInt(document.getElementById('numOfTNQues').value)
    numberOfTNDSQues = parseInt(document.getElementById('numOfTNDSQues').value)
    numberOfTLNQues = parseInt(document.getElementById('numOfTLNQues').value)
    testContent = document.getElementById('testContent').value;

    hideElement('AIPopup')
    showElement('loadingPopup')
    showElement("auroraLayer")

    resObj = await fetch4AITest(numberOfTNQues, numberOfTNDSQues, numberOfTLNQues, testContent)

    if (!resObj.status) {
        if (resObj.msg == "ERR_GEN_AI_NOT_AVAILABLE") {
            alert("Tính năng chỉ có thể dùng cho phiên bản PRO")
        } else {
            alert(`Lỗi nội bộ: ${resObj.msg}`)
        }
        hideElement('loadingPopup')
        hideElement('auroraLayer')
        return
    }

    document.getElementById('textEditor').value = resObj.content

    fetch4Questions()

    hideElement('loadingPopup')
    hideElement('auroraLayer')
})

blackLayer.addEventListener('click', () => {
    hideElement('AIPopup')
    hideElement('auroraLayer')
})

closeButton.addEventListener('click', () => {
    hideElement('AIPopup')
    hideElement('auroraLayer')
})

async function fetch4AITest(numberOfTNQues, numberOfTNDSQues, numberOfTLNQues, testContent) {
    res = await fetch("/API/genAI", {
        method: "POST", body: JSON.stringify({
            numberOfQuesTN: numberOfTNQues,
            numberOfQuesTNDS: numberOfTNDSQues,
            numberOfQuesTLN: numberOfTLNQues,
            content: testContent
        })
    })

    return await res.json()
}

function showElement(e) {
    document.getElementById(e).classList.remove('hidden-element')
}

function hideElement(e) {
    document.getElementById(e).classList.add('hidden-element')
}

async function AIMode() {
    showElement('AIPopup')
}