function logout() {
    fetch(`${window.API_URL}/logout`,
        {
            method: "POST",
            credentials: 'include'
        })
        .then(data => {
            document.location = "/login"
        })
}