function confirmation(question) {
    return new Promise((resolve) => {
        document.body.style.overflow = 'hidden'
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
            resolve(input.value);
        }
        function onClose(e) {
            e.preventDefault()
            cleanup();
            resolve(null);
        }
        function cleanup() {
            dialog.classList.add('closing');
            dialog.addEventListener('animationend', () => {
                dialog.close();
                dialog.classList.remove('closing');
                confirmBtn.removeEventListener("click", onConfirm)
                closeBtn.removeEventListener("click", onClose)
                document.body.style.overflow = 'auto';
            }, { once: true });
        }

        confirmBtn.addEventListener("click", onConfirm)
        closeBtn.addEventListener("click", onClose)
    })
}