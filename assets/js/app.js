// Calculate today's word
const today = new Date();
const start = new Date(2021, 5, 19, 0, 0, 0, 0);
const t = today.setHours(0, 0, 0, 0) - start.setHours(0, 0, 0, 0);
const f = Math.round(t / 864e5);
const currentWord = words[f % words.length];

// Modal behavior
const modal = document.getElementById("help-modal");
const btn = document.getElementById("help-button");
const span = document.getElementsByClassName("close")[0];
btn.onclick = function () {
  modal.style.display = "block";
};
span.onclick = function () {
  modal.style.display = "none";
};
window.onclick = function (event) {
  if (event.target == modal) {
    modal.style.display = "none";
  }
};

// generate letter divs
const createBox = (idx) => {
  const outerDiv = document.createElement("div");
  const innerDiv = document.createElement("div");

  innerDiv.setAttribute("id", "letter-" + idx);
  innerDiv.setAttribute("class", "tile");
  innerDiv.setAttribute("data-state", "empty");
  innerDiv.setAttribute("data-animation", "idle");
  innerDiv.addEventListener("click", () => reveal(idx));

  outerDiv.appendChild(innerDiv);

  return outerDiv;
};

document.addEventListener("DOMContentLoaded", () => {
  const row = document.querySelector("#board .row");
  for (let idx = 0; idx < 5; idx++) {
    row.appendChild(createBox(idx));
  }
});

// reveals the assigned letter to the given div
function reveal(index) {
  // set div to `correct`, add inner text of the character, and set animation to
  // construct div name -> letter-0
  div = document.getElementById(`letter-${index}`);
  div.innerText = currentWord[index];
  div.setAttribute("data-state", "correct");
  div.setAttribute("data-animation", "pop");
}

const explanation = `https://api.dictionaryapi.dev/api/v2/entries/en/${currentWord}`;
const synonyms = `https://api.datamuse.com/words?ml=${currentWord}`;

function fetchURL(url) {
  return fetch(url)
  .then((response) => {
    return response.json();
  }).catch(error => {
    console.log("failed to fetch url: ", error);
  })
}


// reveal the fetched hints for the word.
async function revealHints() {
  p = document.getElementById("hints-content");
  // make sure we only fetch once after button click and then only toggle hide.
  if (p.innerHTML == "") {
    const ex = await fetchURL(explanation)
    const syn = await fetchURL(synonyms)
    p = document.getElementById("hints-content");
    p.innerHTML += "<p>" + ex[0].meanings[0].definitions[0].definition + "</p>";
    p.innerHTML += "<p>synonym: " + syn[0].word + "</p>";
  }

  const targetDiv = document.getElementById("hints-show");
  if (targetDiv.style.display !== "none") {
    targetDiv.style.display = "none";
  } else {
    targetDiv.style.display = "block";
  }
}
