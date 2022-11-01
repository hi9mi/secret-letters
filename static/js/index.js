const keyInput = document.getElementById('key-input');
const keyButton = document.getElementById('key-button');
const keyError = document.getElementById('key-error');

let timeoutId;

const keyButtonHandler = () => {
  const url = keyInput.value;

  if (url.trim().length) {
    navigator.clipboard
      .writeText(url)
      .then(() => {
        keyError.classList.add('hidden');
        keyButton.textContent = 'Copied!';
        clearTimeout(timeoutId);
        timeoutId = setTimeout(() => {
          keyButton.textContent = 'Copy';
        }, 2500);
      })
      .catch(() => {
        keyError.classList.remove('hidden');
      });
  }
};

window.addEventListener('load', () => {
  if (window.location.pathname.match(/key/)) {
    keyButton.addEventListener('click', keyButtonHandler);
  }
});
