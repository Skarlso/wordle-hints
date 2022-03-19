const today = new Date();
const start = new Date(2021,5,19,0,0,0,0);
const t = today.setHours(0,0,0,0) - start.setHours(0,0,0,0);
const f = Math.round(t/864e5);
const currentWord = words[f%words.length]

// reveals the assigned letter to this div
function reveal(index) {
    // set div to `correct`, add inner text of the character, and set animation to
    // construct div name -> letter-0
    div = document.getElementById(`letter-${index}`)
    div.innerText = currentWord[index]
    div.setAttribute("data-state", "correct")
    div.setAttribute("data-animation", "pop")
}