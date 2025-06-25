// status-popup.js

class ModernStatusPopup extends HTMLElement {
    constructor() {
        super();
        this.shadow = this.attachShadow({ mode: 'open' });
        this.timeoutId = null;

        // Embed Font Awesome CSS directly for icons to work inside Shadow DOM.
        const fontAwesomeStyle = `
            @import url('https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0-beta3/css/all.min.css');
        `;

        this.shadow.innerHTML = `
            <style>
                ${fontAwesomeStyle} /* Import Font Awesome here */

                .status-popup {
                    display: none; /* Hidden by default */
                    position: fixed;
                    top: 50%;
                    left: 50%;
                    transform: translate(-50%, -50%);
                    background-color: #2E1A47; /* Brokefolio deep purple */
                    border: 1px solid #ccc; /* Default border */
                    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.4); /* Enhanced shadow */
                    z-index: 1001;
                    padding: 30px 40px; /* More padding */
                    border-radius: 15px;
                    text-align: center;
                    min-width: 280px; /* Ensures a decent size */
                    transition: opacity 0.3s ease-out, transform 0.3s ease-out;
                    opacity: 0;
                    pointer-events: none; /* Allows clicks to pass through when hidden */
                }

                .status-popup.visible {
                    display: block; /* Use block for actual display */
                    opacity: 1;
                    pointer-events: auto; /* Enable clicks when visible */
                }

                .status-overlay {
                    display: none; /* Initially hidden */
                    position: fixed;
                    top: 0;
                    left: 0;
                    width: 100%;
                    height: 100%;
                    background-color: rgba(0, 0, 0, 0.5); /* Darker overlay */
                    z-index: 1000;
                    transition: opacity 0.3s ease-out;
                    opacity: 0;
                    pointer-events: none;
                }

                .status-overlay.visible {
                    display: block;
                    opacity: 1;
                    pointer-events: auto;
                }

                .popup-content {
                    display: flex;
                    flex-direction: column;
                    align-items: center;
                }

                #icon {
                    font-size: 4.5em; /* Slightly adjusted size */
                    margin-bottom: 15px;
                    display: block;
                    line-height: 1;
                }

                .success-icon { color: #4CAF50; } /* Green */
                .error-icon { color: #E74C3C; } /* Red */
                .info-icon { color: #00BFFF; } /* Bright Blue */
                .loading-icon { color: #B76BC4; } /* Brokefolio accent color */

                #message {
                    margin-bottom: 10px; /* Slightly less margin */
                    font-size: 22px; /* Adjusted font size */
                    color: white;
                    font-family: 'Inter', sans-serif; /* Ensure font consistency */
                }
            </style>
            <div class="status-popup">
                <div class="popup-content">
                    <span id="icon"></span>
                    <p id="message"></p>
                </div>
            </div>
            <div class="status-overlay"></div>
        `;
        this.popupElement = this.shadow.querySelector('.status-popup');
        this.overlayElement = this.shadow.querySelector('.status-overlay');
        this.iconElement = this.shadow.querySelector('#icon');
        this.messageElement = this.shadow.querySelector('#message');
    }

    /**
     * Displays the status popup.
     * @param {string} type - 'success', 'error', 'info', or 'loading'.
     * @param {string} message - The message to display.
     * @param {number} [duration=3000] - Duration in milliseconds for auto-hide. Ignored for 'loading'.
     * @param {string|null} [redirectUrl=null] - URL to redirect to after auto-hide.
     */
    show(type, message, duration = 3000, redirectUrl = null) {
        clearTimeout(this.timeoutId); // Clear any previous timeout

        this.messageElement.textContent = message;
        this.popupElement.classList.add('visible');
        this.overlayElement.classList.add('visible');

        // Set icon based on type
        this.iconElement.className = ''; // Clear previous classes
        this.iconElement.innerHTML = ''; // Clear previous content

        switch (type) {
            case 'success':
                this.iconElement.className = 'fas fa-check-circle success-icon';
                break;
            case 'error':
                this.iconElement.className = 'fas fa-times-circle error-icon';
                break;
            case 'info':
                this.iconElement.className = 'fas fa-info-circle info-icon';
                break;
            case 'loading':
                this.iconElement.className = 'fas fa-spinner fa-spin loading-icon';
                break;
            default:
                this.iconElement.className = 'fas fa-info-circle info-icon'; // Default to info
        }

        // Auto-hide unless it's a 'loading' type
        if (type !== 'loading') {
            this.timeoutId = setTimeout(() => {
                this.hide();
                if (redirectUrl) {
                    window.location.replace(redirectUrl);
                }
            }, duration);
        }
    }

    hide() {
        clearTimeout(this.timeoutId); // Clear any pending hide timeout
        this.popupElement.classList.remove('visible');
        this.overlayElement.classList.remove('visible');
    }
}

customElements.define('modernstatus-popup', ModernStatusPopup);