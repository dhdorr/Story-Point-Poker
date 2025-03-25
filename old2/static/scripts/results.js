console.log("results page!");

// let intervalID = setInterval(doCountdown, 1000)
// let timeRemaining = 30


// what i want to do is poll the server every X miliseconds to see if a new round has started.

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
    let shouldGoToNewRound = document.getElementById("newRound");
    
    if (shouldGoToNewRound.value == 'true') {
        document.removeEventListener("fx:init", (evt)=>{
        console.log("polling turned off")
        });
        const cBtn = document.getElementById("continueBtn");
        cBtn.click();
    }
})

// function doCountdown() {
//     timeRemaining -= 1;
//     if (timeRemaining <= 0) {
//         const cBtn = document.getElementById("continueBtn");
//         clearInterval(intervalID);
//         intervalID = null;

//         cBtn.click();
//     }
// }