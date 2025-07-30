function try_login() {
    fetch(`${window.API_URL}/isLogged`, {credentials: "include" })
        .then(response => {
            if (response.status !== 200) {
                window.location = "login.html"
            }
        })
}
try_login()
setInterval(try_login, 5000)

const browser = document.getElementById("browser")
let wd = "/"

ls(wd)
function ls(directory) {
    wd = directory
    fetch(`${window.API_URL}/fileserver/ls?path=${wd}`, {
            method: "GET",
            credentials: "include"
        })
        .then(response => response.json())
        .then(data => {
            let up = ""
            if (wd != "/") {
                let parent = wd.split("/").slice(0, -1).join("/")
                if (parent === "") parent = "/"
                up = `<tr onClick="ls('${parent}')">
                <td>..</td>
                <td></td>
                <td></td>
                </tr>`
            }
            browser.innerHTML = up + data.data.map(file => {
                if (file.isDirectory) {
                    return `<tr onClick="ls('${file.path}')">
                    <td>${file.name}</td>
                    <td>${file.modifiedOn}</td>
                    <td>${file.size}</td>
                    </tr>`
                }
                else {
                    return `<tr onClick="readFile('${file.path}')">
                    <td>${file.name}</td>
                    <td>${file.modifiedOn}</td>
                    <td>${file.size}</td>
                    </tr>`
                }
            })
        })
}

