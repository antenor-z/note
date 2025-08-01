function try_login() {
    fetch(`${window.API_URL}/isLogged`, { credentials: "include" })
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
    document.getElementById("path").innerText = wd
    const parseDate = dateStr => ({
        year: dateStr.substring(0, 4),
        month: dateStr.substring(5, 7),
        day: dateStr.substring(8, 10),
        time: dateStr.substring(11, 16)
    });
    fetch(`${window.API_URL}/fileserver/ls?path=${wd}`, {
        method: "GET",
        credentials: "include"
    })
        .then(response => response.json())
        .then(data => {
            let up = `
                <tr>
                <th>name</th>
                <th class="colapse">modified</th>
                <th class="colapse">size</th>
                <th></th>
            `
            if (wd != "/") {
                let parent = wd.split("/").slice(0, -1).join("/")
                if (parent === "") parent = "/"
                up += `<tr style="cursor: pointer">
                <td onClick="ls('${parent}')">⬆️..</td>
                <td onClick="ls('${parent}')"></td>
                <td onClick="ls('${parent}')"></td>
                <td onClick="ls('${parent}')"></td>
                </tr>`
            }
            browser.innerHTML = up + data.data.map(file => {
                let m = parseDate(file.modifiedOn)

                if (file.isDirectory) {
                    return `<tr style="cursor: pointer">
                    <td onClick="ls('${file.path}')">📂${file.name}</td>
                    <td onClick="ls('${file.path}')" class="fit colapse">${m.day}/${m.month}/${m.year} ${m.time}</td>
                    <td onClick="ls('${file.path}')" class="fit colapse">${file.size}</td>
                    <td onClick="rm('${file.path}')">del</td>
                    </tr>`
                }
                else {
                    return `<tr style="cursor: pointer">
                    <td onClick="readFile('${file.path}')">📄${file.name}</td>
                    <td onClick="readFile('${file.path}')" class="fit colapse">${m.day}/${m.month}/${m.year} ${m.time}</td>
                    <td onClick="readFile('${file.path}')" class="fit colapse">${file.size}</td>
                    <td class="fit" onClick="rm('${file.path}')">del</td>
                    </tr>`
                }
            }).join("")
        })
}


document.getElementById("uploadForm").addEventListener("submit", function (event) {
    event.preventDefault();

    const fileInput = document.getElementById("file");
    const file = fileInput.files[0];
    if (!file) {
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
        .finally(() => {
            fileInput.value = ""
            uploadFile.close()
        })
});

document.addEventListener("dragenter", function (event) {
    uploadFile.showModal();
});


async function mkdir() {
    const dirName = await confirmation("New folder name");
    if (dirName === null) {
        return;
    }
    fetch(`${window.API_URL}/fileserver/mkdir?path=${wd}/${dirName}`, {
        method: "POST",
        credentials: "include"
    })
    .then(response => response.json())
    .then(data => { 
        ls(wd)
    })
}

async function rm(directory) {
    const confirm = await confirmation(`Delete "${directory}". Write 'delete' to confirm.`)

    if (confirm !== "delete") {
        return
    }

    fetch(`${window.API_URL}/fileserver/rm?path=${directory}`, {
        method: "DELETE",
        credentials: "include"
    })
    .then(response => response.json())
    .then(data => {
        ls(wd)
    })
}

async function readFile(path) {
    downloadFileFromPath(path)
}

function downloadFileFromPath(path) {
    const link = document.createElement('a');
    link.href = `${window.API_URL}/fileserver/read?path=${path}`;
    link.download = path.split("/").slice(-2, -1);

    document.body.appendChild(link);

    link.click();
    document.body.removeChild(link);
}
