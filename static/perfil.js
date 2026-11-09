import { Modal } from "./modules/modals.js";
import { SideMenuToggle } from "./modules/header.js";

const sideMenuButton = document.querySelector('#sideMenuToggle');
const sideMenu = document.querySelector('#sideMenu');

const profileBanner = document.querySelector('#profileBanner');
const profileAvatar = document.querySelector('#profileAvatar')
const profileName = document.querySelector('#profileName');
const profileClass = document.querySelector('#profileClass');
const profileDescription = document.querySelector('#profileDescription');

sideMenuButton.addEventListener('click', (_) => SideMenuToggle(sideMenuButton, sideMenu));

function getCookie(name) {
    const match = document.cookie.match(new RegExp('(?:^|; )' + name + '=([^;]*)'));
    return match ? decodeURIComponent(match[1]) : null;
}

function getCodeFromUrl() {
    const params = new URLSearchParams(window.location.search);
    for (const [, value] of params) {
        if (value) return value;
    }
    return null;
}

let modal = new Modal();
modal.Create();

const code = getCodeFromUrl() ?? getCookie('id');

try {
    let response = await fetch('/api/usuario', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: code,
    });

    let data = await response.json();

    if (!response.ok) {
        throw data.html;
    }

    profileBanner.src = `/static/ilustracoes/banners/${data.banner || 1}.png`;
    profileName.textContent = data.name;
    profileDescription.textContent = data.description;

    if (data.avatar.Body > 0) {
        profileAvatar.innerHTML += `<img src="/static/ilustracoes/avatar/base/${data.avatar.Body}.png" alt="Base">`
    }
    if (data.avatar.Shirt > 0) {
        profileAvatar.innerHTML += `<img src="/static/ilustracoes/avatar/camisas/${data.avatar.Shirt}.png" alt="Camisa">`
    }
    if (data.avatar.Pants > 0) {
        profileAvatar.innerHTML += `<img src="/static/ilustracoes/avatar/calcas/${data.avatar.Pants}.png" alt="Calça">`
    }
    if (data.avatar.Accessory > 0) {
        profileAvatar.innerHTML += `<img src="/static/ilustracoes/avatar/acessorios/${data.avatar.Accessory}.png" alt="Acessorio">`
    }

    response = await fetch('/api/usuario/classes', {
        method: 'POST',
        headers: { 'Content-Type': 'application/plain' },
        body: data.class,
    });

    data = await response.json();

    if (!response.ok) {
        throw data.html
    }

    profileClass.textContent = data.name
} catch (err) {
    console.error('Falha ao carregar perfil:', err);
}