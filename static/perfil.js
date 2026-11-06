import { SideMenuToggle } from "./modules/header.js";

const sideMenuButton = document.querySelector('#sideMenuToggle');
const sideMenu = document.querySelector('#sideMenu');

sideMenuButton.addEventListener('click', (_) => SideMenuToggle(sideMenuButton, sideMenu));
