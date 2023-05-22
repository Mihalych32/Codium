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


textarea.addEventListener('input', function (event) {
  let text = this.value;
});

const submitCode = async () => {
  const codeResponse = await fetch(`${process.env.BACKEND_URL}/api/submit/`, {
    method: "POST",
    body: {
      lang_slug: "cpp",
      content: text
    }
  });

  if (codeResponse.status === 200) {
    const serverAnswer = await codeResponse.json();
    const serverResponseArea = document.getElementById('serverResponse');
    serverResponseArea.innerHTML = serverAnswer;
  } else {
    const serverResponseArea = document.getElementById('serverResponse');
    serverResponseArea = "Compilation error!"
  }
}

let submitButton = document.querySelector('submit-button');
submitButton = addEventListener('click', submitCode);