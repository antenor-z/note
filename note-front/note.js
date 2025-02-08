function toggleAddNote() {
    const addNote = document.getElementById('addNote')
    if (addNote.style.display === "none") {
        addNote.style.display = "grid"
    }
    else {
        addNote.style.display = "none"
    }
}

function try_login() {
    fetch('http://localhost:5000/isLogged', {credentials: "include" })
        .then(response => {
            if (response.status !== 200) {
                window.location = "login.html"
            }
        })
}
try_login()
function fetchCategories() {
    const container = document.getElementById('categories-container')
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
    console.log(categoryMap)
    fetch('http://localhost:5000/category', {credentials: "include"})

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

                catDiv.appendChild(checkBox)
                catDiv.appendChild(label)
                container.appendChild(catDiv)
            })
        })
}

function fetchNotes() {
    fetchCategories()
    const catList = []
    for (var child = document.getElementById("categories-container").firstChild; child !== null; child = child.nextSibling) {
        const div = child
        if (div === null) continue
        const checkBox = div.firstChild
        if (checkBox === null) continue
        if (checkBox.checked) {
            catList.push(checkBox.name)
        }
    }

    fetch('http://localhost:5000/noteCat',
        {
            method: "POST",
            body: JSON.stringify({ categories: catList }),
            credentials: "include"
        }
    )
        .then(response => response.json())
        .then(data => {
            const container = document.getElementById('notes-container')
            container.innerHTML = ''
            data["data"].forEach(element => {
                const noteDiv = document.createElement('div')
                noteDiv.classList.add('note')
                noteDiv.classList.add('box')

                const title = document.createElement('h2')
                title.textContent = element.title

                const btnEdit = document.createElement('button')
                btnEdit.innerText = 'Edit'
                btnEdit.id = "editButtonNote" + element.id
                btnEdit.addEventListener("click", () => editNoteToggle(element.id))

                const categories = document.createElement('h3')
                categories.textContent = `[ ${element.categories.map(cat => cat.name).join(', ')} ]`

                const createdAt = document.createElement('h4')
                const date = element.createdAt
                const year = date.substring(0, 4)
                // 2025-02-04T12:00:37.753623325-03:00 
                const month = date.substring(5, 7)
                const day = date.substring(8, 10)
                const time = date.substring(11, 16)
                const dateUpdated = element.updatedAt
                const yearU = dateUpdated.substring(0, 4)
                // 2025-02-04T12:00:37.753623325-03:00 
                const monthU = dateUpdated.substring(5, 7)
                const dayU = dateUpdated.substring(8, 10)
                const timeU = dateUpdated.substring(11, 16)
                createdAt.textContent = `Created ${day}/${month}/${year} ${time} | Updated ${dayU}/${monthU}/${yearU} ${timeU}`

                // Begin edit element
                const editDiv = document.createElement('div')
                editDiv.classList.add("grid")
                editDiv.id = "note" + element.id
                editDiv.style.display = "none"
                const editTitle = document.createElement('input')
                editTitle.value = element.title
                editTitle.id = "editTitle" + element.id
                const editContent = document.createElement('textarea')
                editContent.value = element.content
                editContent.id = "editContent" + element.id
                const editCategories = document.createElement('input')
                editCategories.value = element.categories.map(cat => cat.name).join(',')
                editCategories.id = "editCategories" + element.id

                // Begin edit.action buttons
                const editActionContainer = document.createElement('div')
                const btnUpdate = document.createElement('button')
                btnUpdate.innerText = 'Update'
                btnUpdate.addEventListener("click", () => updateNote(element.id))
                editActionContainer.appendChild(btnUpdate)
                const btnDelete = document.createElement('button')
                btnDelete.innerText = 'Delete'
                btnDelete.addEventListener("click", () => deleteNote(element.id))
                editActionContainer.appendChild(btnDelete)
                const btnClose = document.createElement('button')
                btnClose.innerText = 'Close'
                btnClose.addEventListener("click", () => editNoteToggle(element.id))
                editActionContainer.appendChild(btnClose)
                // End edit.action buttons
                editDiv.appendChild(editTitle)
                editDiv.appendChild(editContent)
                editDiv.appendChild(editCategories)
                editDiv.appendChild(editActionContainer)
                // End edit element

                const content = document.createElement('p')
                content.textContent = element.content;

                noteDiv.appendChild(title)
                noteDiv.appendChild(categories)
                noteDiv.appendChild(createdAt)
                noteDiv.appendChild(content)
                noteDiv.appendChild(editDiv)
                noteDiv.appendChild(btnEdit)

                container.appendChild(noteDiv)
            })
        })
        .catch(error => {
            console.error('Error fetching:', error)
        })
}

function editNoteToggle(noteId) {
    const noteDiv = document.getElementById("note" + noteId)
    const noteButton = document.getElementById("editButtonNote" + noteId)
    if (noteDiv.style.display === "grid") {
        noteDiv.style.display = "none"
        noteButton.style.display = "block"
    }
    else {
        noteDiv.style.display = "grid"
        noteButton.style.display = "none"
    }
}

function updateNote(noteId) {
    const editTitle = document.getElementById("editTitle" + noteId).value;
    const editContent = document.getElementById("editContent" + noteId).value;
    const editCategories = document.getElementById("editCategories" + noteId).value.split(",");

    fetch(`http://localhost:5000/note/${noteId}`, {
        method: "PUT",
        body: JSON.stringify({ title: editTitle, content: editContent, categories: editCategories }),
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
}

function sendNote() {
    const noteTitle = document.getElementById("noteTitle").value
    const noteContent = document.getElementById("noteContent").value
    const categories = document.getElementById("noteCategories").value.split(",")
    document.getElementById("noteTitle").value = ""
    document.getElementById("noteContent").value = ""
    document.getElementById("noteCategories").value = ""
    fetch(`http://localhost:5000/note`, {
        method: "POST",
        body: JSON.stringify({ title: noteTitle, content: noteContent, categories: categories }),
        credentials: "include"
    })
        .then(response => response.json())
        .then(data => {
            fetchNotes()
        })
        .catch(error => {
            console.error('Error fetching:', error)
        })
}
function sendNote() {
    const noteTitle = document.getElementById("noteTitle").value
    const noteContent = document.getElementById("noteContent").value
    const categories = document.getElementById("noteCategories").value.split(",")
    document.getElementById("noteTitle").value = ""
    document.getElementById("noteContent").value = ""
    document.getElementById("noteCategories").value = ""
    fetch(`http://localhost:5000/note`, {
        method: "POST",
        body: JSON.stringify({ title: noteTitle, content: noteContent, categories: categories }),
        credentials: "include"
    })
        .then(response => response.json())
        .then(data => {
            window.location = "index.html"
        })
        .catch(error => {
            console.error('Error fetching:', error)
        });
}
function deleteNote(noteId) {
    const person = prompt("Delete this note? Write 'delete' if you are sure.");
    if (person !== "delete") {
        return
    }
    fetch(`http://localhost:5000/note/${noteId}`, { method: "DELETE", credentials: "include" })
        .then(response => response.json())
        .then(data => {
            fetchNotes()
        })
        .catch(error => {
            console.error('Error fetching:', error)
        })
}

