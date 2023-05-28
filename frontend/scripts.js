let text = "";

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
  text = this.value;
});

const clearCode = () => {
  document.getElementById('codeArea').value = '';
}

const submitCode = async () => {
  const codeResponse = await fetch(`http://localhost:8080/api/submit/`, {
    method: "POST",
    body: JSON.stringify({
      lang_slug: "cpp",
      content: text
    })
  });

  if (codeResponse.status === 200) {
    const serverAnswer = await codeResponse.json();
    let result = serverAnswer.Result;
    const serverResponseArea = document.getElementById('serverResponse');
    serverResponseArea.innerText = serverAnswer.Result;
  } else if (codeResponse.status === 422) {
    const serverAnswer = await codeResponse.json();
    const serverResponseArea = document.getElementById('serverResponse');
    serverResponseArea.innerText = serverAnswer.Result;
    serverResponseArea.classList.add('if-server-error-area');
  }
}

let clearButton = document.getElementById('clr');
clearButton.addEventListener('click', clearCode);

let submitButton = document.getElementById('sbmt');
submitButton.addEventListener('click', submitCode);
