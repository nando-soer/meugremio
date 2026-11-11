import { Modal } from "./modules/modals.js";
import { SideMenuToggle } from "./modules/header.js";
import { ChangePasswordStatus, ClickViewPassword } from "./modules/forms.js";

const sideMenuButton = document.querySelector('#sideMenuToggle');
const sideMenu = document.querySelector('#sideMenu');

const userPassword = document.querySelector('#userPassword');
const passwordView = document.querySelector('#passwordView');
const confirmPasswordView = document.querySelector('#confirmPasswordView');
const userConfirmPassword = document.querySelector('#userConfirmPassword');

const passwordChangeForm = document.querySelector('#credentialForm');

let modal = new Modal();
modal.Create();

const queryString = window.location.search;
const urlParams = new URLSearchParams(queryString);
const token = urlParams.get('token');

try {
    let data = {
        token: token
    };

    console.log(JSON.stringify(data))

    let response = await fetch('/api/usuario/alterar', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    })

    if (!response.ok) {
        let responseData = await response.json();
        throw responseData.html
    }
} catch (error) {
    modal.UpdateContent(error);
    modal.Show();

    document.addEventListener('click', (_) => {
        window.location.href = "/usuario/entrar";
    });
}

sideMenuButton.addEventListener('click', (_) => SideMenuToggle(sideMenuButton, sideMenu));

userPassword.addEventListener('input', (_) => ChangePasswordStatus(userPassword, userConfirmPassword, confirmPasswordStatus));
userConfirmPassword.addEventListener('input', (_) => ChangePasswordStatus(userPassword, userConfirmPassword, confirmPasswordStatus));

passwordView.addEventListener('click', ClickViewPassword);
confirmPasswordView.addEventListener('click', ClickViewPassword);

passwordChangeForm.addEventListener('submit', async (event) => {
    event.preventDefault();

    let data = {
        token: token,
        password: userPassword.value
    };

    let response = await fetch('/api/usuario/alterar', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    })

    let responseData = await response.json();

    modal.UpdateContent(responseData.html);
    modal.Show();
})
