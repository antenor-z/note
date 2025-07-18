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
fetchNotes()
showHidden.addEventListener("click", () => fetchNotes())
priority.addEventListener("click", () => fetchNotes())
function fetchCategories() {
    const container = document.getElementById('categories-container')
    const showHidden = document.getElementById("showHidden").value == "yes"
    const categoryMap = {}
    for (var child = document.getElementById("categories-container").firstChild; child !== null; child = child.nextSibling) {
        const div = child
        if (div === null) continue
        const checkBox = div.firstChild
        if (checkBox === null) continue
        if (checkBox.checked) {
            categoryMap[checkBox.name] = true
        }
    }
    fetch(showHidden ? `${window.API_URL}/category/hidden`:`${window.API_URL}/category`, {credentials: "include"})
        .then(response => response.json())
        .then(data => {
            container.innerHTML = ''
            data["data"].forEach(element => {
                const catDiv = document.createElement('div')
                const checkBox = document.createElement('input')
                checkBox.type = "checkbox"
                checkBox.id = element
                checkBox.name = element
                checkBox.addEventListener("click", () => fetchNotes())
                if (categoryMap[element] === true) {
                    checkBox.checked = true
                }

                const label = document.createElement('label')
                label.htmlFor = element
                label.innerText = element
                label.classList.add("checkbox-label")

                catDiv.appendChild(checkBox)
                catDiv.appendChild(label)
                container.appendChild(catDiv)
            })
        })
}

function fetchNotes() {
    fetchCategories();
    document.getElementById("search").value = "";

    const catList = [];
    const categoriesContainer = document.getElementById("categories-container");
    for (let child = categoriesContainer.firstChild; child !== null; child = child.nextSibling) {
        const div = child;
        if (!div) continue;
        const checkBox = div.firstChild;
        if (!checkBox) continue;
        if (checkBox.checked) catList.push(checkBox.name);
    }

    const showHidden = document.getElementById("showHidden").value == "yes";
    const priorityFilter = parseInt(document.getElementById("priority").value, 10);

    fetch(`${window.API_URL}/note/category`, {
        method: "POST",
        body: JSON.stringify({
            categories: catList,
            showHidden: showHidden,
            priority: priorityFilter
        }),
        credentials: "include"
    })
    .then(response => response.json())
    .then(data => {
        const container = document.getElementById('notes-container');
        const escapeHtml = unsafe => {
            if (!unsafe) return '';
            return unsafe.toString()
                .replace(/&/g, "&amp;")
                .replace(/</g, "&lt;")
                .replace(/>/g, "&gt;")
                .replace(/"/g, "&quot;")
                .replace(/'/g, "&#039;");
        };

        container.innerHTML = data.data.map(element => {
            const parseDate = dateStr => ({
                year: dateStr.substring(0, 4),
                month: dateStr.substring(5, 7),
                day: dateStr.substring(8, 10),
                time: dateStr.substring(11, 16)
            });

            const created = parseDate(element.createdAt);
            const updated = parseDate(element.updatedAt);

            const deadlineInputValue = element.deadline
                ? element.deadline.substring(0, 16)
                : "";

            const priorityOptions = [0, 1, 2, 3].map(p => {
                const labels = ["None", "Low", "Medium", "High"];
                const selected = element.priority === p ? ' selected' : '';
                return `<option value="${p}"${selected}>${labels[p]}</option>`;
            }).join("");

            const hiddenOptions = ["no", "yes"].map(val => {
                const isYes = val === "yes";
                const selected = element.isHidden === isYes ? ' selected' : '';
                const label = isYes ? "Yes" : "No (default)";
                return `<option value="${val}"${selected}>${label}</option>`;
            }).join("");

            const title = escapeHtml(element.title);
            const categories = element.categories.map(cat => escapeHtml(cat.name)).join(', ');
            const editCategories = element.categories.map(cat => escapeHtml(cat.name)).join(',');
            const content = DOMPurify.sanitize(marked.parse(element.content));

            let priorityBadge = "";
            if (element.priority === 1) priorityBadge = '<span style="background-color: gray; font-size: 14px">LOW PRIORITY</span>';
            if (element.priority === 2) priorityBadge = '<span style="background-color: orange; font-size: 14px">MEDIUM PRIORITY</span>';
            if (element.priority === 3) priorityBadge = '<span style="background-color: red; font-size: 14px">HIGH PRIORITY</span>';
            const hiddenBadge = element.isHidden
                ? '<span style="background-color: gray; font-size: 14px">Hidden note</span>'
                : "";

            const displayDeadline = element.deadline
                ? `| ${parseDate(element.deadline).day}/${parseDate(element.deadline).month}/${parseDate(element.deadline).year} ${parseDate(element.deadline).time}`
                : "";

            return `
                <div class="note box">
                    <div class="grid padding20" id="note${element.id}" style="display:none">
                        <label for="editTitle${element.id}">Title</label>
                        <input value="${title}" id="editTitle${element.id}">
                        <label for="editContent${element.id}">Content</label>
                        <textarea id="editContent${element.id}">${escapeHtml(element.content)}</textarea>
                        <label for="editCategories${element.id}">Categories (comma separated)</label>
                        <input value="${editCategories}" id="editCategories${element.id}">

                        <label for="fileInput${element.id}">Attach file</label>
                        <div class="upload-section">
                            <input type="file" id="fileInput${element.id}">
                        </div>

                        <label for="editDeadline${element.id}">Deadline (optional)</label>
                        <input
                            id="editDeadline${element.id}"
                            type="datetime-local"
                            value="${deadlineInputValue}"
                        />

                        <label for="editPriority${element.id}">Priority (optional)</label>
                        <select id="editPriority${element.id}">
                            ${priorityOptions}
                        </select>

                        <label for="editIsHidden${element.id}">Hidden note</label>
                        <select id="editIsHidden${element.id}">
                            ${hiddenOptions}
                        </select>

                        <div class="attachments">
                            ${element.attachments.length ? `
                                Attachments:
                                ${element.attachments.map(att => `
                                    <div class="attachment">
                                        <button class="delete-attachment" onclick="deleteAttachment(${element.id}, ${att.id})">delete</button>
                                        <a href="${window.API_URL}/note/${element.id}/attachment/${att.id}/file" target="_blank">
                                            ${escapeHtml(att.name)}
                                        </a>
                                    </div>
                                `).join('')}` : ''}
                        </div>

                        <div class="edit-action-container">
                            <button onclick="updateNote(${element.id})">Update</button>
                            <button onclick="deleteNote(${element.id})">Delete</button>
                            <button onclick="editNoteToggle(${element.id})">Close</button>
                        </div>
                    </div>

                    <div id="innerNote${element.id}">
                        <div class="box-title">
                            <h2>${title} ${priorityBadge} ${hiddenBadge}</h2>
                            <h3>[ ${categories} ]</h3>
                        </div>
                        <div class="box-content" id="noteContent${element.id}">
                            ${content}
                        </div>
                        <div class="box-bottom">
                            <h4 style="opacity: 0.4">
                            Created ${created.day}/${created.month}/${created.year} ${created.time} |
                            Updated ${updated.day}/${updated.month}/${updated.year} ${updated.time}
                            ${displayDeadline}
                            </h4>

                            <div class="attachments">
                            ${element.attachments.length ? `
                                Attachments:
                                ${element.attachments.map(att => `
                                    <div class="attachment">
                                        <a href="${window.API_URL}/note/${element.id}/attachment/${att.id}/file" target="_blank">
                                            ${escapeHtml(att.name)}
                                        </a>
                                    </div>
                                `).join('')}` : ''}
                            </div>
                            
                            
                            <div class="edit-action-container">
                                <button onclick="editNoteToggle(${element.id})" id="editButtonNote${element.id}">Edit</button>
                                <button onclick="copyNote(${element.id})" id="copyButtonNote${element.id}">Copy</button>
                                <button onclick="deleteNote(${element.id})">Delete</button>
                            </div>
                        </div>
                    </div>
                </div>
            `;
        }).join('');
    })
    .catch(error => console.error('Error fetching:', error));
}


function editNoteToggle(noteId) {
    const noteDiv = document.getElementById("note" + noteId)
    const noteButton = document.getElementById("editButtonNote" + noteId)
    const noteCopyButton = document.getElementById("copyButtonNote" + noteId)
    const innerNote = document.getElementById("innerNote" + noteId)
    const elements = document.querySelectorAll('.delete-attachment')

    if (noteDiv.style.display === "grid") {
        noteDiv.style.display = "none"
        noteButton.style.display = "inline"
        noteCopyButton.style.display = "inline"
        innerNote.style.display = "unset"
        elements.forEach(element => {
            element.style.display = 'none'
        })
    }
    else {
        noteDiv.style.display = "grid"
        noteButton.style.display = "none"
        noteCopyButton.style.display = "none"
        innerNote.style.display = "none"
        document.getElementById("editTitle" + noteId).focus()
        elements.forEach(element => {
            element.style.display = 'unset'
        })
    }
}

function copyNote(noteId) {
    const noteContent = document.getElementById("noteContent" + noteId).innerText
    navigator.clipboard.writeText(noteContent)
}

function updateNote(noteId) {
    const editTitle = document.getElementById("editTitle" + noteId).value;
    const editContent = document.getElementById("editContent" + noteId).value;
    const editCategories = document.getElementById("editCategories" + noteId).value.split(",");
    const editDeadline = document.getElementById("editDeadline" + noteId).value;
    const editPriority = document.getElementById("editPriority" + noteId).value;
    const editIsHidden = document.getElementById("editIsHidden" + noteId).value == "yes";

    const input = document.getElementById(`fileInput${noteId}`)
    const noteEdit = document.getElementById(`note${noteId}`)
    const loading = document.createElement("div")
    loading.innerHTML = `<img style='width: 20px; opacity: 50%' src='/img/spinner.gif'> Loading...`
    noteEdit.appendChild(loading)
    if (input.files.length > 0) {
        uploadFile(noteId, input.files[0])
    }

    fetch(`${window.API_URL}/note/${noteId}`, {
        method: "PUT",
        body: JSON.stringify({ 
            title: editTitle,
            content: editContent, 
            categories: editCategories,
            deadline: editDeadline != "" ? editDeadline+getLocalTimezoneOffset() : null,
            priority: parseInt(editPriority),
            isHidden: editIsHidden
        }),
        credentials: "include"
    })
        .then(response => response.json())
        .then(data => {
            document.getElementById("note" + noteId).style.display = "none";
            fetchNotes();
        })
        .catch(error => {
            console.error('Error fetching:', error);
        });
    loading.remove()
    window.scrollTo({top: 0, behavior: 'smooth'});
}

function sendNote() {
    const noteTitle = document.getElementById("noteTitle").value
    const noteContent = document.getElementById("noteContent").value
    const categories = document.getElementById("noteCategories").value.split(",")
    const isHidden = document.getElementById("noteIsHidden").value == "yes"
    let priority = parseInt(document.getElementById("notePriority").value)
    let deadline = document.getElementById("noteDeadline").value
    if (priority == 0 ) priority = null
    deadline = deadline ? deadline+getLocalTimezoneOffset() : null
    document.getElementById("noteTitle").value = ""
    document.getElementById("noteContent").value = ""
    document.getElementById("noteCategories").value = ""
    fetch(`${window.API_URL}/note`, {
        method: "POST",
        body: JSON.stringify({ title: noteTitle, content: noteContent, categories: categories, isHidden: isHidden, priority: priority, deadline: deadline }),
        credentials: "include"
    })
        .then(response => response.json())
        .then(data => {
            fetchNotes()
        })
        .catch(error => {
            console.error('Error fetching:', error)
        })
    closeSendNote()
    window.scrollTo({top: 0, behavior: 'smooth'});
}
function closeSendNote() {
    const dialog = document.getElementById("addNote")
    dialog.classList.add('closing');
    dialog.addEventListener('animationend', () => {
        dialog.close();
        dialog.classList.remove('closing');
        confirmBtn.removeEventListener("click", onConfirm)
        closeBtn.removeEventListener("click", onClose)
        document.body.style.overflow = 'auto';
    }, { once: true });
}
async function deleteNote(noteId) {
    const answer = await confirmation("Delete this note? Write 'delete' if you are sure.")
    if (answer !== "delete") {
        return;
    }
    fetch(`${window.API_URL}/note/${noteId}`, { method: "DELETE", credentials: "include" })
        .then(response => response.json())
        .then(data => {
            fetchNotes()
        })
        .catch(error => {
            console.error('Error fetching:', error)
        })
    window.scrollTo({top: 0, behavior: 'smooth'});
}
async function deleteAttachment(noteId, attachmentId) {
    const answer = await confirmation("Delete this attachment? Write 'delete' if you are sure.")
    if (answer !== "delete") {
        return;
    }
    
    fetch(`${window.API_URL}/note/${noteId}/attachment/${attachmentId}`, {
        method: 'DELETE',
        credentials: 'include'
    })
    .then(response => {
        if (response.ok) fetchNotes()
    })
}

function uploadFile(noteId, file) {
    const formData = new FormData()
    formData.append('file', file)

    fetch(`${window.API_URL}/note/${noteId}/attachment`, {
        method: 'POST',
        body: formData,
        credentials: 'include'
    })
    .then(response => {
        if (response.ok) fetchNotes()
    })
}

function getLocalTimezoneOffset() {
    const date = new Date();

    const offsetMinutes = date.getTimezoneOffset();
    const offsetSign = offsetMinutes <= 0 ? '+' : '-';
    const absOffsetMinutes = Math.abs(offsetMinutes);
    const offsetHours = String(Math.floor(absOffsetMinutes / 60)).padStart(2, '0');
    const offsetMins = String(absOffsetMinutes % 60).padStart(2, '0');
    const timezoneOffset = `${offsetSign}${offsetHours}:${offsetMins}`;

    return `:00${timezoneOffset}`;
}
