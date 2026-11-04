function arePasswordsEqual(password, confirmPassword) {
    return (password.value === confirmPassword.value) && (password.value && confirmPassword)
}

export function ChangePasswordStatus(password, confirmPassword, confirmPasswordStatus) {
    confirmPasswordStatus.classList.remove('incorrect-tag', 'correct-tag');
    if (arePasswordsEqual(password, confirmPassword)) {
        confirmPasswordStatus.classList.add('correct-tag');
        confirmPasswordStatus.innerHTML = '* SENHAS IGUAIS :)';
    } else {
        confirmPasswordStatus.classList.add('incorrect-tag');
        confirmPasswordStatus.innerHTML = '* SENHAS NÃO BATEM!';
    }
}

export function ClickViewPassword(event) {
    const button = event.currentTarget;

    const passwordWrap = button.closest('.password-wrap');
    const passwordInput = passwordWrap.querySelector('.text-input');

    const eyeOpen = button.querySelector('.eye-open');
    const eyeClosed = button.querySelector('.eye-closed');

    if (passwordInput) {
        const isPassword = passwordInput.type === 'password';
        passwordInput.type = isPassword ? 'text' : 'password';

        if (eyeOpen && eyeClosed) {
            eyeOpen.style.display = isPassword ? 'none' : 'block';
            eyeClosed.style.display = isPassword ? 'block' : 'none';
        }
    }
}
