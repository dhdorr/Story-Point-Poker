console.log("waiting for players...");

const table_id = localStorage.getItem("tableID");
const passcode = localStorage.getItem("passcode");
const username = localStorage.getItem("username");


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
    let plist = document.getElementById("playersList")
    let pcount = document.getElementById("playerCount")
    let pmax = document.getElementById("playerMax")
    let c = plist.childElementCount

    pcount.innerText = c
    if (c >= pmax.innerText) {
      navigateToPokerTable()
    }
  })

  
function navigateToPokerTable() {
  const intervalID = setInterval(myCallback, 3000);
}

function doNavigation() {
  let btn = document.getElementById("continueBtn")
  btn.click();
}
