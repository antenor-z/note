function login() {
    const username = document.getElementById("username").value
    const password = document.getElementById("password").value
    console.log(password)
    fetch("http://localhost:5000/login",
        {
            method: "POST",
            body: JSON.stringify({ username: username, password: password }),
            credentials: 'include'
        })
        .then(data => {
            if (data.status !== 200) {
                return
            }
            document.location = "/"
        })
}