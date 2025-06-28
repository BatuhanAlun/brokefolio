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


    popMessageElement.className = 'pop-message';
    popMessageElement.classList.add(type);

    popTextElement.textContent = message;


    if (popIconElement) popIconElement.innerHTML = '';

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
            if (popIconElement) popIconElement.innerHTML = '';
            if (popSpinnerElement) popSpinnerElement.style.display = 'inline-block';
            break;
        default:
            if (popIconElement) popIconElement.innerHTML = '<i class="fas fa-info-circle"></i>';
            if (popSpinnerElement) popSpinnerElement.style.display = 'none';
    }

 
    popMessageElement.classList.add('active');


    if (popMessageTimeoutId) {
        clearTimeout(popMessageTimeoutId);
    }


    if (type !== 'loading') {
        popMessageTimeoutId = setTimeout(() => {
            hidePopMessage();
            if (redirectUrl) {
                window.location.replace(redirectUrl);
            }
        }, duration);
    }
}

function hidePopMessage() {
    if (!popMessageElement) return;

    clearTimeout(popMessageTimeoutId);
    popMessageElement.classList.remove('active');

    setTimeout(() => {
        popMessageElement.className = 'pop-message';
        if (popIconElement) popIconElement.innerHTML = '';
        if (popTextElement) popTextElement.textContent = '';
        if (popSpinnerElement) popSpinnerElement.style.display = 'none';
    }, 500); 
}