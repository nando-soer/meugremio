export function SideMenuToggle(button, sideMenu) {
    if (sideMenu.classList.contains('open')) {
        sideMenu.classList.remove('open');
        button.classList.remove('is-open');
    } else {
        sideMenu.classList.add('open');
        button.classList.add('is-open');
    }
}
