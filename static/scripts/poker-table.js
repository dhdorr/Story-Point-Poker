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

// for every fixi request, add the tableID, passcode, and username to the header
document.addEventListener("fx:config", (evt)=>{
    evt.detail.cfg.headers.tableID = localStorage.getItem("tableID");
    evt.detail.cfg.headers.passcode = localStorage.getItem("passcode");
    evt.detail.cfg.headers.username = localStorage.getItem("username");
    console.log(evt.detail.cfg.headers)
  })