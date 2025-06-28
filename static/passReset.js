const resetForm = document.getElementById("resetForm");
const myStatusPopup = document.querySelector('status-popup');


resetForm.addEventListener('submit', async (event) => {
    event.preventDefault();

    const formData = new FormData(resetForm);
    const data = Object.fromEntries(formData);


    try{
        const response = await fetch('http://localhost:8080/api/forgetPassword',{
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',

            },
            body:JSON.stringify(data)
        });

        const result = await response.json();


        if (response.ok) {
            showPopMessage("success", data.message || "İstek Gönderildi!", 2000, "/login");
        }else {
            let errorMessage = `Registration failed with status: ${response.status}`;
            try{
                const errorResult = await response.json;
                if (errorResult && errorResult.message) {
                    errorMessage = errorResult.message;
                }
            }catch(jsonError){
                console.error('Failed to parse error JSON', jsonError)
            }
            showPopMessage("error", data.message || "Hata!", 2000);
        }

    }catch(error){
        console.error('Error:',error);
        showPopMessage("error", data.message || "Bağlantı Hatası!", 2000);
    }
});