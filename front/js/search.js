document.getElementById("search").addEventListener("input", function () { search(this) })

function search(inputElement) {
    const searchText = inputElement.value.toLowerCase().trim();
    const notes = document.getElementById("notes-container").children
    for (let i = 0; i < notes.length; i += 1) {
        const note = notes[i]
        const noteInner = note.children[1]
        const noteTitle = noteInner.children[0].innerText.toLowerCase()
        const noteContent = noteInner.children[1].innerText.toLowerCase()
        if (noteTitle.includes(searchText) || noteContent.includes(searchText)) {
            note.style.display = "block"
        }
        else {
            note.style.display = "none"
        }
    }
}