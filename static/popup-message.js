// static/popup-message.js

// Get the popup element and its internal components
const popMessageElement = document.getElementById('popMessage');
const popIconElement = popMessageElement ? popMessageElement.querySelector('.pop-icon') : null;
const popTextElement = popMessageElement ? popMessageElement.querySelector('.pop-text') : null;
const popSpinnerElement = popMessageElement ? popMessageElement.querySelector('.pop-spinner') : null;

let popMessageTimeoutId = null;

/**
 * Displays a popup message.
 * @param {string} type - 'success', 'error', 'info', or 'loading'.
 * @param {string} message - The message text.
 * @param {number} [duration=3000] - Duration in milliseconds for non-loading messages.
 * @param {string|null} [redirectUrl=null] - URL to redirect to after popup disappears.
 */
function showPopMessage(type, message, duration = 3000, redirectUrl = null) {
    if (!popMessageElement || !popTextElement) {
        console.error("Popup message elements not found in DOM.");
        return;
    }

    // Clear any existing classes and reset for new message
    popMessageElement.className = 'pop-message'; // Reset to base class
    popMessageElement.classList.add(type);

    popTextElement.textContent = message;

    // Set icon and spinner based on type
    if (popIconElement) popIconElement.innerHTML = ''; // Clear existing icon

    switch (type) {
        case 'success':
            if (popIconElement) popIconElement.innerHTML = '<i class="fas fa-check-circle"></i>';
            if (popSpinnerElement) popSpinnerElement.style.display = 'none';
            break;
        case 'error':
            if (popIconElement) popIconElement.innerHTML = '<i class="fas fa-times-circle"></i>';
            if (popSpinnerElement) popSpinnerElement.style.display = 'none';
            break;
        case 'info':
            if (popIconElement) popIconElement.innerHTML = '<i class="fas fa-info-circle"></i>';
            if (popSpinnerElement) popSpinnerElement.style.display = 'none';
            break;
        case 'loading':
            // Spinner is already inside .pop-spinner
            if (popIconElement) popIconElement.innerHTML = ''; // No fixed icon, just spinner
            if (popSpinnerElement) popSpinnerElement.style.display = 'inline-block'; // Show spinner
            break;
        default:
            if (popIconElement) popIconElement.innerHTML = '<i class="fas fa-info-circle"></i>';
            if (popSpinnerElement) popSpinnerElement.style.display = 'none';
    }

    // Make the popup visible
    popMessageElement.classList.add('active');

    // Clear any previous timeout
    if (popMessageTimeoutId) {
        clearTimeout(popMessageTimeoutId);
    }

    // Auto-hide unless it's a 'loading' message
    if (type !== 'loading') {
        popMessageTimeoutId = setTimeout(() => {
            hidePopMessage();
            if (redirectUrl) {
                window.location.replace(redirectUrl);
            }
        }, duration);
    }
}

/**
 * Hides the popup message.
 */
function hidePopMessage() {
    if (!popMessageElement) return;

    clearTimeout(popMessageTimeoutId); // Clear any pending hide timeout
    popMessageElement.classList.remove('active');

    // Optional: clear content after transition to reset for next message
    setTimeout(() => {
        popMessageElement.className = 'pop-message'; // Reset classes
        if (popIconElement) popIconElement.innerHTML = '';
        if (popTextElement) popTextElement.textContent = '';
        if (popSpinnerElement) popSpinnerElement.style.display = 'none';
    }, 500); // Matches CSS transition duration
}