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
    fetch('http://localhost:5000/category')

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
fetchCategories()

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
        }
    )
        .then(response => response.json())
        .then(data => {
            const container = document.getElementById('notes-container')
            container.innerHTML = ''
            data["data"].forEach(element => {
                const noteDiv = document.createElement('div')
                noteDiv.classList.add('note')

                const title = document.createElement('h2')
                title.textContent = element.title

                const btnEdit = document.createElement('button')
                btnEdit.innerText = 'Edit'
                btnEdit.addEventListener("click", () => editNote(element.id))

                const createdAt = document.createElement('h3')
                createdAt.textContent = `Created At: ${element.createdAt}`

                const categories = document.createElement('h4')
                categories.textContent = `Categories: ${element.categories.map(cat => cat.name).join(', ')}`

                // Begin edit element
                const editDiv = document.createElement('div')
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
                const btnUpdate = document.createElement('button')
                btnUpdate.innerText = 'Update'
                btnUpdate.addEventListener("click", () => updateNote(element.id))
                const btnDelete = document.createElement('button')
                btnDelete.innerText = 'Delete'
                btnDelete.addEventListener("click", () => deleteNote(element.id))
                // End edit.action buttons
                editDiv.appendChild(editTitle)
                editDiv.appendChild(editContent)
                editDiv.appendChild(editCategories)
                editDiv.appendChild(btnUpdate)
                editDiv.appendChild(btnDelete)
                // End edit element

                const content = document.createElement('p')
                content.textContent = element.content;

                noteDiv.appendChild(title)
                noteDiv.appendChild(createdAt)
                noteDiv.appendChild(categories)
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
fetchNotes()

function editNote(noteId) {
    document.getElementById("note" + noteId).style.display = "block"
}

function updateNote(noteId) {
    const editTitle = document.getElementById("editTitle" + noteId).value
    const editContent = document.getElementById("editContent" + noteId).value
    const editCategories = document.getElementById("editCategories" + noteId).value.split(",")

    console.log(editTitle)
    console.log(editContent)
    console.log(editCategories)
    fetch(`http://localhost:5000/note/${noteId}`, {
        method: "PUT",
        body: JSON.stringify({ title: editTitle, content: editContent, categories: editCategories }),
    })
        .then(response => response.json())
        .then(data => {
            document.getElementById("note" + noteId).style.display = "none"
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
    })
        .then(response => response.json())
        .then(data => {
            console.log(data)
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
    })
        .then(response => response.json())
        .then(data => {
            console.log(data)
            window.location = "index.html"
        })
        .catch(error => {
            console.error('Error fetching:', error)
        });
}
function deleteNote(noteId) {
    fetch(`http://localhost:5000/note/${noteId}`, { method: "DELETE" })
        .then(response => response.json())
        .then(data => {
            console.log(data)
            fetchNotes()
        })
        .catch(error => {
            console.error('Error fetching:', error)
        })
}