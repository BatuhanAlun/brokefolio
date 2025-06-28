document.addEventListener('DOMContentLoaded', () => {
    const registerForm = document.getElementById('registerForm');
    const nameInput = document.getElementById('name');
    const surnameInput = document.getElementById('surname');
    const usernameInput = document.getElementById('username');
    const emailInput = document.getElementById('email');
    const passwordInput = document.getElementById('password');
    const confirmPasswordInput = document.getElementById('confirm_password');
    const confirmPassSpan = document.getElementById('confirmPass');

    const avatarUploadInput = document.getElementById('avatarUploadInput');
    const avatarPreview = document.getElementById('avatarPreview');


    if (avatarUploadInput && avatarPreview) {
        avatarUploadInput.addEventListener('change', function() {
            const file = this.files[0];
            if (file) {

                const allowedTypes = ['image/jpeg', 'image/png', 'image/gif'];
                if (!allowedTypes.includes(file.type)) {

                    showPopMessage("error", "Geçersiz dosya türü. Sadece JPG, PNG veya GIF yükleyin.", 3000);
                    avatarUploadInput.value = '';
                    avatarPreview.src = "/static/default-avatar.png";
                    return;
                }


                const maxSize = 5 * 1024 * 1024; // 5MB
                if (file.size > maxSize) {
                    showPopMessage("error", "Dosya boyutu çok büyük. Maksimum 5MB.", 3000);
                    avatarUploadInput.value = '';
                    avatarPreview.src = "/static/default-avatar.png";
                    return;
                }

                const reader = new FileReader();
                reader.onload = function(e) {
                    avatarPreview.src = e.target.result;
                };
                reader.readAsDataURL(file);
            } else {
                avatarPreview.src = "/static/default-avatar.png";
            }
        });
    }

    // --- Form Submission Logic ---
    registerForm.addEventListener('submit', async (e) => {
        e.preventDefault(); // Prevent default form submission


        showPopMessage("loading", "Kayıt olunuyor..."); 

        if (passwordInput.value !== confirmPasswordInput.value) {

            showPopMessage("error", "Şifreler eşleşmiyor!", 3000);
            confirmPassSpan.style.color = 'red';
            confirmPassSpan.innerHTML = 'Şifreler Eşleşmiyor';
            return;
        }
        if (passwordInput.value.length < 8) {
            showPopMessage("error", "Şifre en az 8 karakter olmalıdır.", 3000);
            return;
        }

        const formData = new FormData();
        formData.append('name', nameInput.value.trim());
        formData.append('surname', surnameInput.value.trim());
        formData.append('username', usernameInput.value.trim());
        formData.append('email', emailInput.value.trim());
        formData.append('password', passwordInput.value.trim());


        if (avatarUploadInput.files.length > 0) {
            formData.append('avatar', avatarUploadInput.files[0]);
        }

        try {
            const response = await fetch('/api/register', {
                method: 'POST',
                body: formData
            });

            const data = await response.json();

            if (response.ok) {
                showPopMessage("success", data.message || "Kayıt başarılı! Giriş yapabilirsiniz.", 2000, "/login");
            } else {

                const errorMessage = data.message || "Kayıt işlemi başarısız oldu.";
                showPopMessage("error", errorMessage, 4000);
            }
        } catch (error) {
            console.error('Kayıt işlemi sırasında ağ hatası:', error);
            showPopMessage("error", "Kayıt işlemi sırasında bir ağ hatası oluştu.", 4000);
        }
    });

});