function confirmation(question) {
    return new Promise((resolve) => {
        document.body.style.overflow = 'hidden';
        const dialog = document.getElementById("deleteModal")
        const input = document.getElementById("answer")
        const confirmBtn = document.getElementById("confirmBtn")
        const message = dialog.querySelector("p")

        message.textContent = question

        input.value = ""

        dialog.showModal()

        function onConfirm(e) {
            e.preventDefault()
            cleanup()
            dialog.close()
            resolve(input.value);
        }
        function onClose() {
            cleanup();
            resolve(null);
        }
        function cleanup() {
            confirmBtn.removeEventListener("click", onConfirm)
            dialog.removeEventListener("close", onClose)
            document.body.style.overflow = 'auto';
        }

        confirmBtn.addEventListener("click", onConfirm)
        dialog.addEventListener("close", onClose)
    })
}