const openLPButton = document.getElementById('openLP')
const openExButton = document.getElementById('openEx')
const openStButton = document.getElementById('openSt')

openLPButton.addEventListener('click', () => {
    window.location.href = window.location.href.replace("/Home", "/LivePreview")
});

openStButton.addEventListener('click', () => {
    window.location.href = window.location.href.replace("/Home", "/StartTest")
});

openExButton.addEventListener('click', () => {
    window.location.href = window.location.href.replace("/Home", "/Export")
});
