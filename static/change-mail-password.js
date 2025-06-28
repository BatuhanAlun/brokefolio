document.addEventListener('DOMContentLoaded', () => {
    const changePasswordForm = document.getElementById('changePasswordForm');
    const newPasswordField = document.getElementById('newPassword');
    const confirmNewPasswordField = document.getElementById('confirmNewPassword');

    const userEmail = window.APP_DATA.userEmail;
    console.log(userEmail);

    if (!changePasswordForm || !newPasswordField || !confirmNewPasswordField) {
        console.error("Change Password page elements not found. Cannot attach listeners.");
        return;
    }

    changePasswordForm.addEventListener('submit', async (event) => {
        event.preventDefault(); // Prevent default form submission

        showPopMessage("loading", "Şifre değiştiriliyor...");

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



        try {
            const response = await fetch('/api/user/change-password-mail', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    newPassword: newPassword,
                    email: userEmail,
                })
            });
            const result = await response.json();

            if (response.ok) {
                showPopMessage("success", result.message || "Şifreniz başarıyla değiştirildi!", 2000, "/login");

                newPasswordField.value = '';
                confirmNewPasswordField.value = '';
            } else {
                let errorMessage = result.message || `Şifre değiştirme başarısız oldu: ${response.statusText}`;
                if (response.status === 400) {
                    showPopMessage("error", errorMessage, 4000);
                    newPasswordField.focus();
                } else {
                    showPopMessage("error", errorMessage, 4000);
                }
            }
        } catch (error) {
            console.error('Şifre değiştirme işlemi sırasında ağ hatası:', error);
            showPopMessage("error", "Şifre değiştirme işlemi sırasında bir ağ hatası oluştu.", 4000);
        }
    });
});