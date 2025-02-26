function login() {
    const username = document.getElementById("username").value
    const password = document.getElementById("password").value
    document.getElementById("username").value = ""
    document.getElementById("password").value = ""
    document.getElementById('username').focus()
    fetch(`${window.API_URL}/login`,
        {
            method: "POST",
            body: JSON.stringify({ username: username, password: password }),
            credentials: 'include'
        })
        .then(data => {
            if (data.status !== 200) {
                document.getElementById('authfail').style.display = "block"
                return
            }
            document.location = "/"
        })
}
document.addEventListener("keydown", (event) => {
    if (event.key === "Enter") {
        login()
    }
})
window.addEventListener("load", () => {
    document.getElementById('username').focus()
})
