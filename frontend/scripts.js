document.addEventListener('keydown', function (event) {
  if (event.key === 'Tab') {
    event.preventDefault();
  }
});

const textarea = document.getElementById('codeArea');

textarea.addEventListener('keydown', function (event) {
  if (event.key === 'Tab') {
    event.preventDefault();

    const start = this.selectionStart;
    const end = this.selectionEnd;

    this.value = this.value.substring(0, start) + '    ' + this.value.substring(end);

    this.selectionStart = this.selectionEnd = start + 4;
  }
});