async function logout() {
    res = await confirmation("This will end your session. Exit now?")
    if (res !== "exit") {
        return
    }
    fetch(`${window.API_URL}/logout`,
        {
            method: "POST",
            credentials: 'include'
        })
        .then(data => {
            document.location = "/login"
        })
}