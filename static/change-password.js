// static/change-password.js

document.addEventListener('DOMContentLoaded', () => {
    const changePasswordForm = document.getElementById('changePasswordForm');
    const currentPasswordField = document.getElementById('currentPassword');
    const newPasswordField = document.getElementById('newPassword');
    const confirmNewPasswordField = document.getElementById('confirmNewPassword');
    // REMOVED: const myStatusPopup = document.querySelector('modernstatus-popup'); // No longer needed

    if (!changePasswordForm || !currentPasswordField || !newPasswordField || !confirmNewPasswordField) {
        console.error("Change Password page elements not found. Cannot attach listeners.");
        return;
    }

    changePasswordForm.addEventListener('submit', async (event) => {
        event.preventDefault(); // Prevent default form submission

        // --- CHANGE: Use the global showPopMessage function ---
        showPopMessage("loading", "Şifre değiştiriliyor...");

        const currentPassword = currentPasswordField.value.trim();
        const newPassword = newPasswordField.value.trim();
        const confirmNewPassword = confirmNewPasswordField.value.trim();

        // --- Client-Side Validation ---
        if (newPassword !== confirmNewPassword) {
            showPopMessage("error", "Yeni şifreler eşleşmiyor!", 3000);
            newPasswordField.value = '';
            confirmNewPasswordField.value = '';
            newPasswordField.focus();
            return;
        }

        if (newPassword.length < 8) {
            showPopMessage("error", "Yeni şifreniz en az 8 karakter olmalıdır.", 3000);
            newPasswordField.value = '';
            confirmNewPasswordField.value = '';
            newPasswordField.focus();
            return;
        }

        if (currentPassword === newPassword) {
            showPopMessage("error", "Mevcut şifrenizle yeni şifreniz aynı olamaz.", 3000);
            newPasswordField.value = '';
            confirmNewPasswordField.value = '';
            newPasswordField.focus();
            return;
        }
        // --- End Client-Side Validation ---

        try {
            // Retrieve the JWT token from cookie just before making the request
            const authToken = getCookie('authToken'); // Assuming 'authToken' is your cookie name
            if (!authToken) {
                showPopMessage("error", "Yetkilendirme hatası: Lütfen tekrar giriş yapın.", 4000);
                setTimeout(() => { window.location.href = "/logout"; }, 2000);
                return;
            }


            const response = await fetch('/api/user/change-password', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${authToken}` // Include authorization token
                },
                body: JSON.stringify({
                    currentPassword: currentPassword,
                    newPassword: newPassword,
                })
            });

            const result = await response.json();

            if (response.ok) {
                // --- CHANGE: Use the global showPopMessage function with redirectUrl ---
                showPopMessage("success", result.message || "Şifreniz başarıyla değiştirildi!", 2000, "/profile");
                currentPasswordField.value = '';
                newPasswordField.value = '';
                confirmNewPasswordField.value = '';
            } else {
                let errorMessage = result.message || `Şifre değiştirme başarısız oldu: ${response.statusText}`;
                if (response.status === 401 || response.status === 403) {
                    errorMessage = "Oturum süreniz dolmuş veya yetkiniz yok. Lütfen tekrar giriş yapın.";
                    // --- CHANGE: Use the global showPopMessage function with redirectUrl ---
                    showPopMessage("error", errorMessage, 4000, "/logout");
                } else if (response.status === 400) {
                    showPopMessage("error", errorMessage, 4000);
                    currentPasswordField.value = '';
                    currentPasswordField.focus();
                } else {
                    showPopMessage("error", errorMessage, 4000);
                }
            }
        } catch (error) {
            console.error('Şifre değiştirme işlemi sırasında ağ hatası:', error);
            showPopMessage("error", "Şifre değiştirme işlemi sırasında bir ağ hatası oluştu.", 4000);
        }
    });

    // Helper function to get a cookie by name (if not already defined globally or in another script)
    function getCookie(name) {
        const value = `; ${document.cookie}`;
        const parts = value.split(`; ${name}=`);
        if (parts.length === 2) return parts.pop().split(';').shift();
        return null;
    }
});