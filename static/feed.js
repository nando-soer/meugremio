import { Modal } from "./modules/modals.js";
import { SideMenuToggle } from "./modules/header.js";

const sideMenuButton = document.querySelector('#sideMenuToggle');
const sideMenu = document.querySelector('#sideMenu');

let modal = new Modal();
modal.Create();

sideMenuButton.addEventListener('click', (_) => SideMenuToggle(sideMenuButton, sideMenu));
