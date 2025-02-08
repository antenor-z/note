function logout() {
    fetch("http://localhost:5000/logout",
        {
            method: "POST",
            credentials: 'include'
        })
        .then(data => {
            document.location = "/login"
        })
}