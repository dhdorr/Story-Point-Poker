console.log("waiting for players...");

const table_id = localStorage.getItem("tableID");
const passcode = localStorage.getItem("passcode");
const username = localStorage.getItem("username");
const maxPlayers = parseInt(document.getElementById("maxPlayers").innerText)

const tl = document.getElementById("timeLimit")
localStorage.setItem("timeLimit", tl.value)

// fixi polling extension
document.addEventListener("fx:init", (evt)=>{
    let elt = evt.target
    if (elt.matches("[ext-fx-poll-interval]")){
      // wait for the non-bubbling fx:inited event on the element so the __fixi property is available
      elt.addEventListener("fx:inited", ()=>{
          // squirrel away in case we want to call clearInterval() later
          elt.__fixi.pollInterval = setInterval(()=>{
              elt.dispatchEvent(new CustomEvent("poll"))
          }, parseInt(elt.getAttribute("ext-fx-poll-interval")))
      })
    }
  })

document.addEventListener("fx:swapped", (evt)=>{
    let pcount = parseInt(document.getElementById("playerCount").innerText)
    let is_ready = parseInt(document.getElementById("isReady").value)

    if (pcount >= maxPlayers || is_ready) {
      document.removeEventListener("fx:init", (evt)=>{
        console.log("polling turned off")
      });
      navigateToPokerTable()
    }
  })


function navigateToPokerTable() {
  const intervalID = setInterval(doNavigation, 3000);
}

function doNavigation() {
  let btn = document.getElementById("continueBtn")
  btn.click();
}
