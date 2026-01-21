// Theme handling - follows system setting, allows temporary toggle
function getSystemTheme() {
    return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
}

function applyTheme(theme) {
    document.documentElement.setAttribute('data-theme', theme);
    document.querySelector('.theme-icon').textContent = theme === 'dark' ? '🌙' : '☀️';
}

function toggleTheme() {
    const current = document.documentElement.getAttribute('data-theme');
    const next = current === 'dark' ? 'light' : 'dark';
    applyTheme(next);
}

// Apply system theme on load
document.documentElement.setAttribute('data-theme', getSystemTheme());

// Listen for system theme changes
window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
    applyTheme(e.matches ? 'dark' : 'light');
});

// Typing effect for hero
const heroText = "I build things.";
let heroIndex = 0;
function typeHero() {
    if (heroIndex < heroText.length) {
        document.getElementById('hero-text').textContent += heroText.charAt(heroIndex);
        heroIndex++;
        setTimeout(typeHero, 80);
    }
}

// Focus input
function focusInput() {
    document.getElementById('cmd-input').focus();
}

// Handle clear command via HTMX
document.body.addEventListener('htmx:beforeSwap', function(evt) {
    if (evt.detail.xhr.getResponseHeader('HX-Retarget') === '#terminal-output') {
        document.getElementById('terminal-output').innerHTML = '';
        evt.detail.shouldSwap = false;
    }
});

// Initialize on DOM ready
document.addEventListener('DOMContentLoaded', function() {
    applyTheme(getSystemTheme());
    setTimeout(typeHero, 500);
    setTimeout(focusInput, 100);
});
