// static/profile.js

// Ensure popup-message.js is loaded BEFORE this script in your HTML
// for 'showPopMessage' function to be available globally.

function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
    return null;
}

document.addEventListener('DOMContentLoaded', () => {
    // --- Get DOM Elements ---
    // Buttons
    const changePasswordBtn = document.getElementById('changePasswordBtn');
    const editProfileBtn = document.getElementById('editProfileBtn');
    const saveProfileBtn = document.getElementById('saveProfileBtn');
    const cancelEditBtn = document.getElementById('cancelEditBtn');
    const deleteAccountBtn = document.getElementById('deleteAccountBtn');

    // Display elements
    const profileUsernameDisplay = document.getElementById('profileUsernameDisplay'); // Top username
    const profileEmailDisplay = document.getElementById('profileEmailDisplay'); // Top email
    const profileNameDisplay = document.getElementById('profileNameDisplay');
    const profileSurnameDisplay = document.getElementById('profileSurnameDisplay');
    const profileUsernameEditableDisplay = document.getElementById('profileUsernameEditableDisplay'); // User info grid username
    const profileEmailEditableDisplay = document.getElementById('profileEmailEditableDisplay'); // User info grid email

    // Input elements
    const profileNameInput = document.getElementById('profileNameInput');
    const profileSurnameInput = document.getElementById('profileSurnameInput');
    const profileUsernameInput = document.getElementById('profileUsernameInput');
    const profileEmailInput = document.getElementById('profileEmailInput');

    // Avatar related elements
    const avatarImage = document.getElementById('avatarImage');
    const avatarEditOverlay = document.getElementById('avatarEditOverlay');
    const avatarUploadInput = document.getElementById('avatarUploadInput');
    const profileAvatarSection = document.querySelector('.profile-avatar-section'); // Parent for editing class

    // Modal & Popup
    const confirmationModal = document.getElementById('confirmationModal');
    const cancelDeleteBtn = document.getElementById('cancelDeleteBtn');
    const confirmDeleteBtn = document.getElementById('confirmDeleteBtn');
    const authToken = getCookie('authToken');

    // --- State Variable ---
    let isEditing = false; // Tracks whether the profile is currently in edit mode
    let selectedFile = null; // To store the selected new avatar file

    // --- Utility Function to Toggle Edit Mode ---
    function toggleEditMode(enable) {
        isEditing = enable;

        // Toggle visibility of display spans vs. input fields
        const infoItems = document.querySelectorAll('.user-info-grid .info-item');
        infoItems.forEach(item => {
            const displaySpan = item.querySelector('span');
            const editInput = item.querySelector('input.edit-input');

            if (displaySpan && editInput) {
                if (enable) {
                    displaySpan.style.display = 'none';
                    editInput.style.display = 'block'; // Make input visible
                    editInput.value = displaySpan.textContent; // Pre-fill input with current display value
                } else {
                    displaySpan.style.display = 'inline'; // Or 'block' depending on desired layout
                    editInput.style.display = 'none';
                }
            }
        });

        // Toggle visibility of action buttons
        if (enable) {
            editProfileBtn.style.display = 'none';
            saveProfileBtn.style.display = 'flex'; // Show Save button
            cancelEditBtn.style.display = 'flex';   // Show Cancel button
            changePasswordBtn.style.display = 'none'; // Hide change password button during edit
            deleteAccountBtn.style.display = 'none'; // Hide delete button during edit
            profileAvatarSection.classList.add('editing'); // Enable avatar overlay
        } else {
            editProfileBtn.style.display = 'flex'; // Show Edit button
            saveProfileBtn.style.display = 'none';
            cancelEditBtn.style.display = 'none';
            changePasswordBtn.style.display = 'flex'; // Show other buttons again
            deleteAccountBtn.style.display = 'flex';
            profileAvatarSection.classList.remove('editing'); // Disable avatar overlay
            // Reset avatar preview if not saved
            if (selectedFile) {
                avatarImage.src = avatarImage.dataset.originalSrc || avatarImage.src; // Revert to original if cancelled
                selectedFile = null;
            }
        }
    }

    // --- Event Listeners ---

    // Change Password Button (Existing Logic)
    if (changePasswordBtn) {
        changePasswordBtn.addEventListener('click', () => {
            showPopMessage("info", "Şifre Değiştirme Sayfasına Yönlendiriliyorsunuz", 1500);
            setTimeout(() => {
                window.location.href = "/change-password";
            }, 1500);
        });
    } else {
        console.error("Error: changePasswordBtn element not found!");
    }

    // Edit Profile Button (New Logic)
    if (editProfileBtn) {
        editProfileBtn.addEventListener('click', () => {
            toggleEditMode(true); // Enable edit mode
            // Store original avatar source for potential revert on cancel
            avatarImage.dataset.originalSrc = avatarImage.src;
        });
    } else {
        console.error("Error: editProfileBtn element not found!");
    }

    // Cancel Edit Button (New Logic)
    if (cancelEditBtn) {
        cancelEditBtn.addEventListener('click', () => {
            // Revert inputs to original values (optional, but good UX)
            profileNameInput.value = profileNameDisplay.textContent;
            profileSurnameInput.value = profileSurnameDisplay.textContent;
            profileUsernameInput.value = profileUsernameEditableDisplay.textContent;
            profileEmailInput.value = profileEmailEditableDisplay.textContent;

            toggleEditMode(false); // Disable edit mode
            showPopMessage("info", "Profil düzenleme iptal edildi.", 1500);
        });
    } else {
        console.error("Error: cancelEditBtn element not found!");
    }

    // NEW: Handle avatar image selection
    if (avatarUploadInput && avatarImage) {
        avatarUploadInput.addEventListener('change', (event) => {
            const file = event.target.files[0];
            if (file) {
                selectedFile = file; // Store the file
                const reader = new FileReader();
                reader.onload = (e) => {
                    avatarImage.src = e.target.result; // Display preview
                };
                reader.readAsDataURL(file);
            }
        });
    }

    // Save Profile Button (Modified Logic for Avatar Upload)
    if (saveProfileBtn) {
        saveProfileBtn.addEventListener('click', async () => {
            showPopMessage("loading", "Profil güncelleniyor...");

            const formData = new FormData();
            formData.append('name', profileNameInput.value.trim());
            formData.append('surname', profileSurnameInput.value.trim());
            formData.append('username', profileUsernameInput.value.trim());
            formData.append('email', profileEmailInput.value.trim());

            if (selectedFile) {
                formData.append('avatar', selectedFile); // Append the new avatar file
            }

            try {


                const response = await fetch('/api/user/update-profile', {
                    method: 'PUT',
                    headers: {
                        // For FormData, 'Content-Type' is usually set automatically by the browser
                        // Do NOT manually set 'Content-Type': 'application/json' when using FormData
                        'Authorization': `Bearer ${authToken}`
                    },
                    body: formData // Send FormData directly
                });

                if (response.ok) {
                    const responseData = await response.json(); // Assuming backend sends back updated data or new avatar URL

                    // Update displayed text fields
                    profileNameDisplay.textContent = formData.get('name');
                    profileSurnameDisplay.textContent = formData.get('surname');
                    profileUsernameDisplay.textContent = `@${formData.get('username')}`;
                    profileUsernameEditableDisplay.textContent = formData.get('username');
                    profileEmailDisplay.textContent = formData.get('email');
                    profileEmailEditableDisplay.textContent = formData.get('email');

                    // Update avatar image if a new one was uploaded and a new URL is returned
                    if (responseData.avatarURL) {
                        avatarImage.src = responseData.avatarURL;
                    } else if (selectedFile) {
                         // If no new URL, but a file was selected, use the preview. This is a fallback.
                         // Ideally, backend confirms the new URL.
                         // For a Go backend, you might need to re-fetch profile or update `{{.AvatarURL}}`
                         // with the new path after successful upload.
                    } else {
                        // If no new avatar, and no file was selected, ensure initials are updated if name/surname changed
                        const initials = (formData.get('name').charAt(0) || '') + (formData.get('surname').charAt(0) || '');
                        if (initials) {
                             avatarImage.src = `https://via.placeholder.com/200/B76BC4/FFFFFF?text=${initials.toUpperCase()}`;
                        }
                    }

                    selectedFile = null; // Clear the selected file after successful upload
                    toggleEditMode(false);
                    showPopMessage("success", "Profil başarıyla güncellendi!", 2000);
                } else {
                    const errorData = await response.json();
                    showPopMessage("error", `Profil güncellenirken hata oluştu: ${errorData.message || response.statusText}`, 4000);
                    // Revert avatar preview on error
                    avatarImage.src = avatarImage.dataset.originalSrc || avatarImage.src;
                    selectedFile = null;
                }
            } catch (error) {
                console.error('Profil güncelleme hatası:', error);
                showPopMessage("error", `Profil güncellenirken bir ağ hatası oluştu: ${error.message || 'Bilinmeyen Hata'}`, 4000);
                // Revert avatar preview on network error
                avatarImage.src = avatarImage.dataset.originalSrc || avatarImage.src;
                selectedFile = null;
            }
        });
    } else {
        console.error("Error: saveProfileBtn element not found!");
    }

    if (deleteAccountBtn && confirmationModal) {
        deleteAccountBtn.addEventListener('click', () => {
            console.log("Delete button clicked! Attempting to show modal.");
            confirmationModal.classList.add('active');
            document.body.style.overflow = 'hidden';
        });
    } else {
        console.error("Error: deleteAccountBtn or confirmationModal element not found!");
    }

    if (cancelDeleteBtn && confirmationModal) {
        cancelDeleteBtn.addEventListener('click', () => {
            confirmationModal.classList.remove('active');
            document.body.style.overflow = '';
        });
    } else {
        console.error("Error: cancelDeleteBtn or confirmationModal element not found!");
    }

    if (confirmDeleteBtn && confirmationModal) {
        confirmDeleteBtn.addEventListener('click', async () => {
            confirmationModal.classList.remove('active');
            document.body.style.overflow = '';
            showPopMessage("loading", "Hesap siliniyor...");

            try {
                const response = await fetch('/api/user/delete-account', {
                    method: 'DELETE',
                    headers: {
                        'Content-Type': 'application/json',
                    }
                });

                if (response.ok) {
                    showPopMessage('success', 'Hesabınız başarıyla silindi.', 2000, "/logout");
                } else {
                    const errorData = await response.json();
                    showPopMessage("error", `Hesap silinirken bir hata ile karşılaşıldı: ${errorData.message || response.statusText}`, 4000);
                }
            } catch (error) {
                console.error("Hesap silme hatası", error);
                showPopMessage("error", `Hesap silinirken bir ağ hatası ile karşılaşıldı: ${error.message || 'Bilinmeyen Hata'}`, 4000);
            }
        });
    } else {
        console.error("Error: confirmDeleteBtn or confirmationModal element not found!");
    }
});