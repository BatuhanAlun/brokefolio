const registerForm = document.getElementById("registerForm");
const messageDiv = document.getElementById("messageDiv")

const myStatusPopup = document.querySelector('status-popup');


registerForm.addEventListener('submit', async (event) => {
    event.preventDefault();

    const formData = new FormData(registerForm);
    const data = Object.fromEntries(formData);



    try{
        const response = await fetch('http://127.0.0.1:8080/api/register',{
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',

            },
            body:JSON.stringify(data)
        });

        if (response.ok) {
            const result = await response.json();
            myStatusPopup.show(true, 'Registered Succesfully', 2000);
            setTimeout(function() {
            window.location.replace("http://127.0.0.1:5500/login.html");
            }, 3000);
            
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
            myStatusPopup.show(false, errorMessage, 2000)
        }

    }catch(error){
        console.error('Error:', error);
        myStatusPopup.show(false, 'Something Went Wrong', 2000)
    }
});
