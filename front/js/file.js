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

console.log(confirmation("aaaa"))
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
                up = `<tr>
                <td style="cursor: pointer;" onClick="ls('${parent}')">⬆️..</td>
                <td></td>
                <td></td>
                </tr>`
            }
            browser.innerHTML = up + data.data.map(file => {
                if (file.isDirectory) {
                    return `<tr>
                    <td style="cursor: pointer;" onClick="ls('${file.path}')">📂${file.name}</td>
                    <td>${file.modifiedOn}</td>
                    <td>${file.size}</td>
                    </tr>`
                }
                else {
                    return `<tr>
                    <td style="cursor: pointer;" onClick="readFile('${file.path}')">📄${file.name}</td>
                    <td>${file.modifiedOn}</td>
                    <td>${file.size}</td>
                    </tr>`
                }
            })
        })
}


document.getElementById("uploadForm").addEventListener("submit", function(event) {
    event.preventDefault();

    const fileInput = document.getElementById("file");
    const file = fileInput.files[0];
    if (!file) {
        alert("Please select a file to upload.");
        return;
    }

    const formData = new FormData();
    formData.append("file", file);

    fetch(`${window.API_URL}/fileserver/write?path=${encodeURIComponent(wd)}`, {
        method: "POST",
        credentials: "include",
        body: formData
    })
    .then(response => response.json())
    .then(result => {
        ls(wd);
    })
    .catch(err => {
        console.error(err);
        alert("Upload failed.");
    })
    .finally( () => {
        fileInput.value = ""
    })
});

