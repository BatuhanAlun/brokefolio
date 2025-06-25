// static/register.js

// Ensure popup-message.js is loaded BEFORE this script in your HTML
// for 'showPopMessage' function to be available globally.

document.addEventListener('DOMContentLoaded', () => {
    const registerForm = document.getElementById('registerForm');
    const nameInput = document.getElementById('name');
    const surnameInput = document.getElementById('surname');
    const usernameInput = document.getElementById('username');
    const emailInput = document.getElementById('email');
    const passwordInput = document.getElementById('password');
    const confirmPasswordInput = document.getElementById('confirm_password');
    const confirmPassSpan = document.getElementById('confirmPass'); // For password match message
    // Removed: const myStatusPopup = document.querySelector('modernstatus-popup'); // This line is removed

    // New elements for avatar upload
    const avatarUploadInput = document.getElementById('avatarUploadInput');
    const avatarPreview = document.getElementById('avatarPreview');

    // --- Avatar Preview Logic ---
    if (avatarUploadInput && avatarPreview) {
        avatarUploadInput.addEventListener('change', function() {
            const file = this.files[0];
            if (file) {
                // Basic file type validation (optional, backend should do robust validation)
                const allowedTypes = ['image/jpeg', 'image/png', 'image/gif'];
                if (!allowedTypes.includes(file.type)) {
                    // Changed: Using showPopMessage
                    showPopMessage("error", "Geçersiz dosya türü. Sadece JPG, PNG veya GIF yükleyin.", 3000);
                    avatarUploadInput.value = ''; // Clear file input
                    avatarPreview.src = "/static/default-avatar.png"; // Reset preview (corrected default path)
                    return;
                }

                // Basic file size validation (optional, backend should do robust validation)
                const maxSize = 5 * 1024 * 1024; // 5MB
                if (file.size > maxSize) {
                    // Changed: Using showPopMessage
                    showPopMessage("error", "Dosya boyutu çok büyük. Maksimum 5MB.", 3000);
                    avatarUploadInput.value = ''; // Clear file input
                    avatarPreview.src = "/static/default-avatar.png"; // Reset preview
                    return;
                }

                const reader = new FileReader();
                reader.onload = function(e) {
                    avatarPreview.src = e.target.result;
                };
                reader.readAsDataURL(file); // Read file as data URL for preview
            } else {
                avatarPreview.src = "/static/default-avatar.png"; // Reset to default if no file selected
            }
        });
    }

    // --- Form Submission Logic ---
    registerForm.addEventListener('submit', async (e) => {
        e.preventDefault(); // Prevent default form submission

        // Changed: Using showPopMessage
        showPopMessage("loading", "Kayıt olunuyor..."); // Show loading message

        // Client-side password validation (from your HTML script, but also good to have here)
        if (passwordInput.value !== confirmPasswordInput.value) {
            // Changed: Using showPopMessage
            showPopMessage("error", "Şifreler eşleşmiyor!", 3000);
            confirmPassSpan.style.color = 'red';
            confirmPassSpan.innerHTML = 'Şifreler Eşleşmiyor'; // Update message
            return;
        }
        if (passwordInput.value.length < 8) {
            // Changed: Using showPopMessage
            showPopMessage("error", "Şifre en az 8 karakter olmalıdır.", 3000);
            return;
        }

        // Create FormData object to handle both text inputs and file uploads
        const formData = new FormData();
        formData.append('name', nameInput.value.trim());
        formData.append('surname', surnameInput.value.trim());
        formData.append('username', usernameInput.value.trim());
        formData.append('email', emailInput.value.trim());
        formData.append('password', passwordInput.value.trim());

        // Append the avatar file if one is selected
        if (avatarUploadInput.files.length > 0) {
            formData.append('avatar', avatarUploadInput.files[0]);
        }

        try {
            const response = await fetch('/api/register', { // Your existing registration API endpoint
                method: 'POST',
                // IMPORTANT: Do NOT set 'Content-Type': 'application/json' when using FormData.
                // The browser sets 'Content-Type': 'multipart/form-data' automatically,
                // along with the correct boundary.
                body: formData // Send the FormData object directly
            });

            const data = await response.json(); // Assuming your backend returns JSON

            if (response.ok) { // HTTP status 200-299 indicates success
                // Changed: Using showPopMessage
                showPopMessage("success", data.message || "Kayıt başarılı! Giriş yapabilirsiniz.", 2000, "/login");
            } else {
                // Backend sent an error response (e.g., 400 Bad Request, 409 Conflict)
                const errorMessage = data.message || "Kayıt işlemi başarısız oldu.";
                // Changed: Using showPopMessage
                showPopMessage("error", errorMessage, 4000);
            }
        } catch (error) {
            console.error('Kayıt işlemi sırasında ağ hatası:', error);
            // Changed: Using showPopMessage
            showPopMessage("error", "Kayıt işlemi sırasında bir ağ hatası oluştu.", 4000);
        }
    });

});