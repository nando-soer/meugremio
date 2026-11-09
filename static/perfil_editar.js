import { Modal } from "./modules/modals.js";
import { SideMenuToggle } from "./modules/header.js";

const sideMenuButton = document.querySelector('#sideMenuToggle');
const sideMenu = document.querySelector('#sideMenu');

const profileForm = document.querySelector('.credential-form');
const studentName = document.querySelector('#studentName');
const studentClass = document.querySelector('#studentClass');
const studentDescription = document.querySelector('#studentDescription');
const studentBanner = document.querySelector('#studentBanner');
const studentAvatar = document.querySelector('#studentAvatar');

const openBanners = document.querySelector('#openBanners');
const modalBanners = document.querySelector('#modalBanners');
const bannerOption = document.querySelectorAll('.banner-option');
const closeBanners = document.querySelector('#closeBanners');

const openAvatarEdit = document.querySelector('#openAvatarEdit');
const saveAvatarBtn = document.querySelector('.btn-save-avatar');

const modalClose = document.querySelector('#modalClose');
const avatarModal = document.querySelector('#avatarModal');

const avatarOptions = document.querySelectorAll('.option-switch');
const avatarPreview = document.querySelector('#preview');

const AVATAR_FOLDERS = ['base', 'camisas', 'calcas', 'acessorios'];
const previewTable = { 'base': 0, 'camisas': 1, 'calcas': 2, 'acessorios': 3 };

// FUNÇÕES

function getCookie(name) {
    const match = document.cookie.match(new RegExp('(?:^|; )' + name + '=([^;]*)'));
    return match ? decodeURIComponent(match[1]) : null;
}

function reloadSelection(previewSelection, previewTable) {
    avatarOptions.forEach((option) => {
        const optionSelection = option.dataset.selection.split(',');
        if (optionSelection[1] == previewSelection[previewTable[optionSelection[0]]]) {
            option.classList.add('selected');
        } else {
            option.classList.remove('selected');
        }
    });
}

function reloadPreview(selection, previewTable, avatarPreview) {
    const avatarElement = avatarPreview.children[previewTable[selection[0]]];
    if (selection[1] == 0) {
        avatarElement.classList.add('inactive');
    } else {
        avatarElement.classList.remove('inactive');
        avatarElement.src = `/static/ilustracoes/avatar/${selection[0]}/${selection[1]}.png`;
    }
}

function changeToOption(option, avatarPreview, previewSelection, previewTable) {
    const optionSelection = option.dataset.selection.split(',');
    previewSelection[previewTable[optionSelection[0]]] = optionSelection[1];

    const selectionString = previewSelection.join(',');
    avatarPreview.dataset.selection = selectionString;
    localStorage.setItem("selection", selectionString);

    reloadSelection(previewSelection, previewTable);
    reloadPreview(optionSelection, previewTable, avatarPreview);
}

function applySelectionToStudentAvatar(selectionString) {
    if (!studentAvatar) return;
    const values = selectionString.split(',');
    values.forEach((value, index) => {
        const el = studentAvatar.children[index];
        if (!el) return;
        if (parseInt(value) === 0) {
            el.classList.add('inactive');
        } else {
            el.classList.remove('inactive');
            el.src = `/static/ilustracoes/avatar/${AVATAR_FOLDERS[index]}/${value}.png`;
        }
    });
}

// INTERATIVIDADE

const code = getCookie('id');

let modal = new Modal();
modal.Create();

try {
    let response = await fetch('/api/usuario', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: code,
    });

    const userData = await response.json()

    if (!response.ok) {
        throw userData.html
    }

    if (studentBanner) {
        studentBanner.src = `/static/ilustracoes/banners/${userData.banner || 1}.png`;
        studentBanner.dataset.name = userData.banner || 1;
    }

    if (studentName) studentName.placeholder = userData.name ?? '';
    if (studentDescription) studentDescription.placeholder = userData.description ?? '';

    if (userData.avatar) {
        if (!localStorage.getItem("selection")) {
            avatarPreview.dataset.selection = [
                userData.avatar.Body,
                userData.avatar.Shirt,
                userData.avatar.Pants,
                userData.avatar.Accessory,
            ].join(',');
        } else {
            avatarPreview.dataset.selection = localStorage.getItem("selection");
        }

        if (studentAvatar) {
            const avatarValues = [userData.avatar.Body, userData.avatar.Shirt, userData.avatar.Pants, userData.avatar.Accessory];
            avatarValues.forEach((value, index) => {
                const el = studentAvatar.children[index];
                if (!el) return;
                if (value > 0) {
                    el.classList.remove('inactive');
                    el.src = `/static/ilustracoes/avatar/${AVATAR_FOLDERS[index]}/${value}.png`;
                } else {
                    el.classList.add('inactive');
                }
            });
        }
    }

    response = await fetch('/api/usuario/classes', {
        method: 'GET'
    });

    let classData = await response.json();

    if (!response.ok) {
        throw classData.html
    }

    classData.forEach((classObject) => {
        const selectClass = new Option(classObject.code, classObject.name)

        if (classObject.name == userData.class) {
            selectClass.selected = true
        }
    
        studentClass.add(selectClass)
    });
} catch (error) {
    modal.UpdateContent(error)
    modal.Show()
}

let previewSelection = avatarPreview.dataset.selection.split(',');

sideMenuButton.addEventListener('click', (_) => SideMenuToggle(sideMenuButton, sideMenu));

openBanners.addEventListener('click', (_) => {
    modalBanners.classList.add('active');
});

closeBanners.addEventListener('click', (_) => {
    modalBanners.classList.remove('active');
});

openAvatarEdit.addEventListener('click', (_) => {
    avatarModal.classList.add('active');
});

modalClose.addEventListener('click', (_) => {
    avatarModal.classList.remove('active');
});

reloadSelection(previewSelection, previewTable);
for (const option in previewTable) {
    reloadPreview([option, previewSelection[previewTable[option]]], previewTable, avatarPreview);
}

bannerOption.forEach((option) => {
    option.addEventListener('click', (event) => {
        let selection = event.target
        studentBanner.src = selection.src
        studentBanner.dataset.name = selection.dataset.name
        modalBanners.classList.remove('active');
    })
})

avatarOptions.forEach((option) => {
    option.addEventListener('click', (event) =>
        changeToOption(event.currentTarget, avatarPreview, previewSelection, previewTable)
    );
});

saveAvatarBtn.addEventListener('click', (_) => {
    applySelectionToStudentAvatar(avatarPreview.dataset.selection);
    avatarModal.classList.remove('active');
});

profileForm.addEventListener('submit', async (event) => {
    event.preventDefault();

    const avatarSelection = (avatarPreview.dataset.selection || '').split(',');

    const payload = {
        name: (studentName.value.trim() || studentName.placeholder || '').trim(),
        class: studentClass.value,
        description: (studentDescription.value.trim() || studentDescription.placeholder || '').trim(),
        banner: parseInt(studentBanner.dataset.name),
        avatar: {
            Body: parseInt(avatarSelection[0]) || 0,
            Shirt: parseInt(avatarSelection[1]) || 0,
            Pants: parseInt(avatarSelection[2]) || 0,
            Accessory: parseInt(avatarSelection[3]) || 0,
        },
    };

    const response = await fetch('/api/usuario/editar', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
    });

    let responseData = await response.json()

    modal.UpdateContent(responseData.html)
    modal.Show()

    localStorage.removeItem('selection');
});