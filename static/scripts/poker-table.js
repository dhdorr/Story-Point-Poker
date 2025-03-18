const t_id = localStorage.getItem("tableID");
const pc = localStorage.getItem("passcode");
const un = localStorage.getItem("username");

console.log("help")

document.addEventListener("fx:init", (evt)=>{
    console.log("init: ", evt);
    if (evt.target.matches("[ext-fx-disable]")){
        console.log("found a disable...");
        var disableSelector = evt.target.getAttribute('ext-fx-disable')
        evt.target.addEventListener('fx:before', ()=>{
            let disableTarget = disableSelector == "" ? evt.target : document.querySelector(disableSelector)
            disableTarget.disabled = true
            disableTarget.classList.add("pre-selected")
            evt.target.addEventListener('fx:after', (afterEvt)=>{
                if (afterEvt.target == evt.target){
                    disableTarget.disabled = false
                    disableTarget.classList.add("selected")
                }
            })
        })
    }
  })

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

// for every fixi request, add the tableID, passcode, and username to the header
document.addEventListener("fx:config", (evt)=>{
    evt.detail.cfg.headers.tableID = t_id
    evt.detail.cfg.headers.passcode = pc
    evt.detail.cfg.headers.username = un
    console.log(evt.detail.cfg.headers)
  })


function UpdateSelectedCard(elm) {
    let cards = document.querySelectorAll(".selected");
    cards.forEach(element => {
        element.classList.remove("selected");
    });
}