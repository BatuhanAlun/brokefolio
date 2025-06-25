const resetForm = document.getElementById("resetForm");
const myStatusPopup = document.querySelector('status-popup');


resetForm.addEventListener('submit', async (event) => {
    event.preventDefault();

    const formData = new FormData(resetForm);
    const data = Object.fromEntries(formData);


    try{
        const response = await fetch('http://127.0.0.1:8080/api/forgetPassword',{
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',

            },
            body:JSON.stringify(data)
        });

        const result = await response.json();


        if (response.ok) {
            myStatusPopup.show("success",result.message,2000);


            setTimeout(function(){
                window.location.replace("http://127.0.0.1:5500/login.html");
            },3000)
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
            myStatusPopup.show("error", errorMessage, 2000)
        }

    }catch(error){
        console.error('Error:',error);
        myStatusPopup.show("error",'Something Went Wrong!',2000);
    }
});