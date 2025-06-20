// status-popup-component.js

class StatusPopup extends HTMLElement {
    constructor() {
        super();
        this.shadow = this.attachShadow({ mode: 'open' });
        this.timeoutId = null;
        this.shadow.innerHTML = `
            <style>
                .status-popup {
                    display: none; /* Initially hidden */
                    position: fixed;
                    top: 50%;
                    left: 50%;
                    transform: translate(-50%, -50%);
                    background-color: #2E1A47;
                    border: 1px solid #ccc;
                    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
                    z-index: 1001;
                    padding: 20px;
                    border-radius: 15px;
                    text-align: center;
                }

                .status-overlay {
                    display: none;
                    position: fixed;
                    top: 0;
                    left: 0;
                    width: 100%;
                    height: 100%;
                    background-color: rgba(0, 0, 0, 0.3);
                    z-index: 1000;
                }

                .popup-content {
                    display: flex;
                    flex-direction: column;
                    align-items: center;

                }

                #icon {
                    font-size: 5em; /* Increased font size */
                    margin-bottom: 10px;
                    display: block; /* Ensure proper spacing */
                    line-height: 1; /* Adjust line height for better vertical centering */
                }

                .success-icon {
                    color: green;
                }

                .error-icon {
                    color: red;
                }

                #message {
                    margin-bottom: 15px;
                    font-size: 24px;
                    color: white;
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

    show(isSuccess, message, duration = 3000, redirectUrl = null) {
        clearTimeout(this.timeoutId);
        this.messageElement.textContent = message;
        this.popupElement.style.display = 'block';
        this.overlayElement.style.display = 'block';

        if (isSuccess) {
            this.iconElement.textContent = '✔';
            this.iconElement.className = 'success-icon';
        } else {
            this.iconElement.textContent = '✖';
            this.iconElement.className = 'error-icon';
        }

        this.timeoutId = setTimeout(() => {
            this.hide();
            if (redirectUrl) {
                window.location.replace(redirectUrl);
            }
        }, duration);
    }

    hide() {
        this.popupElement.style.display = 'none';
        this.overlayElement.style.display = 'none';
    }
}

customElements.define('status-popup', StatusPopup);