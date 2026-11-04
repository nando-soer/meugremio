export class Modal {
    constructor() {
        this.html = `
            <div class="modal-overlay" id="modalOverlay">
                <div class="modal-content" id="modalContent">
                </div>
            </div>
        `;
        this.element = ''
        this.close = ''
    }

    gatherSelf() {
        this.element = document.querySelector('#modalOverlay');
    }

    gatherClose() {
        this.close = this.element.querySelector('#modalClose');
        if (this.close) {
            this.close.addEventListener('click', (_) => {
                this.element.classList.remove('open');
            });
        }
    }

    reload() {
        this.gatherSelf();
        this.gatherClose();
    }

    Create() {
        document.body.insertAdjacentHTML('afterbegin', this.html);
        this.reload();
    }

    Show() {
        this.element.classList.add('open');
    }

    UpdateContent(html) {
        this.html = html;
        this.element.querySelector('#modalContent').innerHTML = html;
        this.reload();
    }
}
