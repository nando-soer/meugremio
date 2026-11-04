import { Modal } from "./modules/modals.js";
import { SideMenuToggle } from "./modules/header.js";
import { ChangePasswordStatus, ClickViewPassword } from "./modules/forms.js";

const sideMenuButton = document.querySelector('#sideMenuToggle');
const sideMenu = document.querySelector('#sideMenu');

const signupForm = document.querySelector('#credentialForm');

const passwordView = document.querySelector('#passwordView');
const confirmPasswordView = document.querySelector('#confirmPasswordView');

const userPassword = document.querySelector('#userPassword');
const userConfirmPassword = document.querySelector('#userConfirmPassword');
const confirmPasswordStatus = document.querySelector('#confirmPasswordStatus');

const passwordStrengthCheck = document.querySelector('#passwordStrengthCheck');
const strengthBar = document.querySelector('#strengthBar');

// FUNÇÕES

function checkPasswordStrength(data) {
    passwordStrengthCheck.innerHTML = data.strength;
    passwordStrengthCheck.classList.remove('strong-password', 'medium-password', 'weak-password');
    if (data.strength == 'SENHA FRACA') {
        passwordStrengthCheck.classList.add('weak-password');
        strengthBar.children[0].classList = 'active-weak';
        strengthBar.children[1].className = '';
        strengthBar.children[2].className = '';
    } else if (data.strength == 'SENHA MÉDIA') {
        passwordStrengthCheck.classList.add('medium-password');
        strengthBar.children[0].classList = 'active-medium';
        strengthBar.children[1].classList = 'active-medium';
        strengthBar.children[2].className = '';
    } else if (data.strength == 'SENHA FORTE') {
        passwordStrengthCheck.classList.add('strong-password');
        strengthBar.children[0].classList = 'active-strong';
        strengthBar.children[1].classList = 'active-strong';
        strengthBar.children[2].classList = 'active-strong';
    }
}

// INTERATIVIDADE

let modal = new Modal();
modal.Create();

sideMenuButton.addEventListener('click', (_) => SideMenuToggle(sideMenuButton, sideMenu));

passwordView.addEventListener('click', ClickViewPassword);
confirmPasswordView.addEventListener('click', ClickViewPassword);

userPassword.addEventListener('input', async (_) => {
    const response = await fetch('/api/teste/senha', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/plain',
        },
        body: userPassword.value,
    });

    if (!response.ok) { }

    const data = await response.json();

    checkPasswordStrength(data);
    ChangePasswordStatus(userPassword, userConfirmPassword, confirmPasswordStatus);
});

userConfirmPassword.addEventListener('input', (_) => ChangePasswordStatus(userPassword, userConfirmPassword, confirmPasswordStatus));

signupForm.addEventListener('submit', async (event) => {
    event.preventDefault();

    const data = {
        name: document.querySelector('#userName').value,
        email: document.querySelector('#userEmail').value,
        class: document.querySelector('#userClass').value,
        password: userPassword.value
    };

    let response = await fetch('/api/usuario/cadastrar', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(data)
    });

    let responseData = await response.json()

    modal.UpdateContent(responseData.html)
    modal.Show()

    let confirmRegisterAction = async (_) => {
        response = await fetch('/api/usuario/confirmar', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/plain',
            },
            body: document.querySelector('#confirmCode').value,
        });

        responseData = await response.json()

        modal.UpdateContent(responseData.html)
        modal.Show()

        let confirmRegister = document.querySelector('#confirmRegister');
        if (confirmRegister) confirmRegister.addEventListener('click', confirmRegisterAction)
    }

    let confirmRegister = document.querySelector('#confirmRegister');
    if (confirmRegister) {
        confirmRegister.addEventListener('click', confirmRegisterAction);
    }
})
