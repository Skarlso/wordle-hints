// Calculate today's word
const today = new Date();
const start = new Date(2021,5,19,0,0,0,0);
const t = today.setHours(0,0,0,0) - start.setHours(0,0,0,0);
const f = Math.round(t/864e5);
const currentWord = words[f%words.length]

// Modal behavior
const modal = document.getElementById("help-modal");
const btn = document.getElementById("help-button");
const span = document.getElementsByClassName("close")[0];
btn.onclick = function() {
  modal.style.display = "block";
}
span.onclick = function() {
  modal.style.display = "none";
}
window.onclick = function(event) {
  if (event.target == modal) {
    modal.style.display = "none";
  }
}

// reveals the assigned letter to the given div
function reveal(index) {
    // set div to `correct`, add inner text of the character, and set animation to
    // construct div name -> letter-0
    div = document.getElementById(`letter-${index}`)
    div.innerText = currentWord[index]
    div.setAttribute("data-state", "correct")
    div.setAttribute("data-animation", "pop")
}