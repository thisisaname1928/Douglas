const modeToggleButton = document.getElementById('modeToggleButton');
const body = document.body;

document.addEventListener('DOMContentLoaded', () => {
    const savedMode = localStorage.getItem('theme');
    if (savedMode === 'dark') {
        body.classList.add('dark-mode');
        modeToggleButton.innerHTML = "dark_mode"
    } else {
        modeToggleButton.innerHTML = "light_mode"
    }
})

modeToggleButton.addEventListener('click', () => {
    body.classList.toggle('dark-mode');

    if (body.classList.contains('dark-mode')) {
        localStorage.setItem('theme', 'dark');
        modeToggleButton.innerHTML = "dark_mode"
    } else {
        localStorage.setItem('theme', 'light');
        modeToggleButton.innerHTML = "light_mode"
    }
});