function confirmation(question, showConfirmInput = true) {
    return new Promise(resolve => {
        const dialog = document.createElement('dialog');
        dialog.classList.add('confirmation-modal');
        let inputText = ""
        if (showConfirmInput) {
            inputText = '<input id="confirm-input" type="text">';
        }
        dialog.innerHTML = `
        <div class="box grid padding20">
            <p>${question}</p>
            ${inputText}
            <div>
                <button id="btn-close">Cancel</button>
                <button id="btn-confirm">Yes</button>
            </div>
        </div>
        `;

        document.body.append(dialog);
        document.body.style.overflow = 'hidden';

        const input = dialog.querySelector('#confirm-input');
        const btnClose = dialog.querySelector('#btn-close');
        const btnConfirm = dialog.querySelector('#btn-confirm');

        dialog.showModal();
        if (input) input.focus();

        function cleanup() {
            dialog.classList.add('closing');
            dialog.addEventListener('animationend', () => {
                dialog.close();
                dialog.remove();
                document.body.style.overflow = '';
            }, { once: true });
        }

        btnConfirm.addEventListener('click', e => {
            e.preventDefault();
            cleanup();
            if (input)
                resolve(input.value);
            else
                resolve("yes");
        });
        btnClose.addEventListener('click', e => {
            e.preventDefault();
            cleanup();
            resolve(null);
        });

        dialog.addEventListener('cancel', e => {
            e.preventDefault();
            cleanup();
            resolve(null);
        });

        dialog.addEventListener('keydown', e => {
            if (e.key === 'Enter') {
                e.preventDefault();
                btnConfirm.click();
            }
        });
    });
}
