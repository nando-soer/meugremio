import { Modal } from "./modules/modals.js";
import { SideMenuToggle } from "./modules/header.js";
import { ChangePasswordStatus, ClickViewPassword } from "./modules/forms.js";

const sideMenuButton = document.querySelector('#sideMenuToggle');
const sideMenu = document.querySelector('#sideMenu');

const loginForm = document.querySelector('#credentialForm');

const userPassword = document.querySelector('#userPassword');
const passwordView = document.querySelector('#passwordView');
const confirmPasswordView = document.querySelector('#confirmPasswordView');
const userConfirmPassword = document.querySelector('#userConfirmPassword');
const passwordRecovery = document.querySelector('#passwordRecovery');

let modal = new Modal();
modal.Create();

loginForm.addEventListener('submit', async (event) => {
    event.preventDefault();

    const data = {
        email: document.querySelector('#userEmail').value,
        password: userPassword.value
    };

    let response = await fetch('/api/usuario/entrar', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(data)
    });

    let responseData = await response.json();

    modal.UpdateContent(responseData.html);
    modal.Show();
});

sideMenuButton.addEventListener('click', (_) => SideMenuToggle(sideMenuButton, sideMenu));

userPassword.addEventListener('input', (_) => ChangePasswordStatus(userPassword, userConfirmPassword, confirmPasswordStatus));
userConfirmPassword.addEventListener('input', (_) => ChangePasswordStatus(userPassword, userConfirmPassword, confirmPasswordStatus));

passwordView.addEventListener('click', ClickViewPassword);
confirmPasswordView.addEventListener('click', ClickViewPassword);

passwordRecovery.addEventListener('click', async (event) => {
    event.preventDefault();

    let response = await fetch('/api/usuario/esqueci', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/plain',
        },
        body: document.querySelector('#userEmail').value
    });

    let responseData = await response.json();

    modal.UpdateContent(responseData.html);
    modal.Show();
});
