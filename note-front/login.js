if (response.status === 401) {
    const username = window.prompt("username", "")
    const password = window.prompt("password", "")
    fetch("http://localhost:5000/login",
        {
            method: "POST",
            body: JSON.stringify({ username: username, password: password }),
            credentials: 'include'
        })
        .then(data => {
            if (data.status !== 200) {
                try_login()
            }
            fetchNotes()
        })
}
})