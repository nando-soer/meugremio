import { Modal } from "./modules/modals.js";
import { SideMenuToggle } from "./modules/header.js";
import { ChangePasswordStatus, ClickViewPassword } from "./modules/forms.js";

const sideMenuButton = document.querySelector('#sideMenuToggle');
const sideMenu = document.querySelector('#sideMenu');

const userPassword = document.querySelector('#userPassword');
const passwordView = document.querySelector('#passwordView');
const confirmPasswordView = document.querySelector('#confirmPasswordView');
const userConfirmPassword = document.querySelector('#userConfirmPassword');

let modal = new Modal();
modal.Create();

sideMenuButton.addEventListener('click', (_) => SideMenuToggle(sideMenuButton, sideMenu));

userPassword.addEventListener('input', (_) => ChangePasswordStatus(userPassword, userConfirmPassword, confirmPasswordStatus));
userConfirmPassword.addEventListener('input', (_) => ChangePasswordStatus(userPassword, userConfirmPassword, confirmPasswordStatus));

passwordView.addEventListener('click', ClickViewPassword);
confirmPasswordView.addEventListener('click', ClickViewPassword);
