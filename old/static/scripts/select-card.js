let bg_skin = "grey";
let bg_text = "black";
let username = "guest";
let selected_card_value = -1;
let sessionID = "";
let passcode = "";

let round_Over = false;

function choose(card) {
    console.log("clicked...");
    if (round_Over) {
        return;
    }
    let cards = document.querySelectorAll(".card");
    cards.forEach(c => {
        if (c.classList.contains("selected")) {
            c.classList.remove("selected");
        }
        if (c.classList.contains("awaiting-selected")) {
            c.classList.remove("awaiting-selected");
        }
    });
    card.classList.add("awaiting-selected");
    selected_card_value = card.getAttribute("value");
}

document.addEventListener("fx:config", (evt) => {
    evt.detail.cfg.headers.testme = "don't panic";
})

function logEvt(evt) {
    console.log(evt.type, evt.target, evt.detail?.cfg || "", evt.detail)

    console.log("my test: ",   evt.detail.cfg.response.status)
}

function handleErrorResponse(evt) {
    console.log("is admin: ", evt.detail.cfg);
    // console.log("is admin? ", evt.detail.cfg.response.headers.forEach((value, key) => {console.log(`${key}: ${value}`);}));
    let error_arr = [401, 404];
    if (error_arr.includes(evt.detail.cfg.response.status)) {
        console.log("there was a 404 error...")
        // document.head.remove();
        let error_dialog = document.getElementById("errorDialog");
        evt.detail.cfg.target = error_dialog;
        evt.detail.cfg.swap = 'innerHTML';
    }
    if (evt.detail.cfg.action == "/results") {
        let tg = document.getElementById("card-wrapper")
        evt.detail.cfg.target = tg
        evt.detail.cfg.swap = 'innerHTML'
    }
}

function handleRequest(evt) {
    if (evt.detail.cfg.action == "/join" || evt.detail.cfg.action == "/create") {
        evt.detail.cfg.headers.bg_skin = bg_skin;
        evt.detail.cfg.headers.bg_text = bg_text;
        evt.detail.cfg.headers.username = username;
        evt.detail.cfg.headers.selected_card_value = selected_card_value;

        sessionID = evt.detail.cfg.body.get("sessionID");
        passcode = evt.detail.cfg.body.get("passcode");
    }

    if (evt.detail.cfg.action == "/choose" || evt.detail.cfg.action == "/playerBox") {
        if (round_Over) {
            evt.detail.cfg.abort();
        }
        evt.detail.cfg.headers.sessionID = sessionID
        evt.detail.cfg.headers.passcode = passcode
        evt.detail.cfg.headers.username = username;
    }

    if (evt.detail.cfg.action == "/results") {
        evt.detail.cfg.headers.sessionID = sessionID
        evt.detail.cfg.headers.passcode = passcode
        evt.detail.cfg.headers.username = username;
    }

    console.log("testing before req: ", evt.detail.cfg);
}

document.addEventListener("fx:after", (evt) => handleErrorResponse(evt))
document.addEventListener("fx:before", (evt) => handleRequest(evt))

function swap_skin(radio) {
    console.log("changing color to " + radio.value);
    let bg = document.getElementById("funBG");
    if (radio.value == "rainbow") {
        bg.style.background = "linear-gradient(90deg, rgba(255,0,0,1) 0%, rgba(255,154,0,1) 10%, rgba(208,222,33,1) 20%, rgba(79,220,74,1) 30%, rgba(63,218,216,1) 40%, rgba(47,201,226,1) 50%, rgba(28,127,238,1) 60%, rgba(95,21,242,1) 70%, rgba(186,12,248,1) 80%, rgba(251,7,217,1) 90%, rgba(255,0,0,1) 100%)";
        bg.style.color = "black";
    } else {
        bg.style.background = null;
        bg.style.backgroundColor = radio.value;
    }
    bg.style.color = radio.getAttribute("data-color");

    bg_skin = radio.value;
    bg_text = radio.getAttribute("data-color");
}

const input = document.querySelector(".username");

input.addEventListener("input", updateUsername);

function updateUsername(evt) {
    // console.log(evt.target.value);
    username = evt.target.value;
}

function udpateTimer() {
    // let timeLeft = 1000; // 100 seconds
    const countdownInterval = setInterval(() => {
        // timeLeft--;
        let pb = document.getElementById("timer")
        if (pb == null) {
            return
        }
        // pb.value = timeLeft
        let tmp = pb.value;
        tmp-= 0.1;
        pb.value = tmp;

        if (pb.value <= 0) {
            round_Over = true;
            clearInterval(countdownInterval);
            console.log("Time's up!");
        }
    }, 1000 / 10); // Update every second
}

function getResults() {
    let tmp = 1
    const countdownInterval = setInterval(() => {
        tmp -= 1;

        if (tmp < 0) {
            clearInterval(countdownInterval);
            getData();
        }
    }, 1000); // Update every second
}

// fixi polling extension
document.addEventListener("fx:init", (evt)=>{
    let elt = evt.target
    console.log("element? ", elt);
    if (elt.matches("[ext-fx-poll-interval]")){
      // wait for the non-bubbling fx:inited event on the element so the __fixi property is available
      elt.addEventListener("fx:inited", ()=>{
          // squirrel away in case we want to call clearInterval() later
          elt.__fixi.pollInterval = setInterval(()=>{
              elt.dispatchEvent(new CustomEvent("poll"))
              if (round_Over && elt.id == "playerBox") {
                clearInterval(elt.__fixi.pollInterval);
                getResults();
              }
          }, parseInt(elt.getAttribute("ext-fx-poll-interval")))
      })
    }
  })


// async function pollServer() {
//     const pollInterval = setInterval(() => {
//         getData();
//         console.log("server polled");
//     }, 1000); // Update every second
// }

function getData() {
    let rb = document.getElementById("resultsBtn");
    rb.click();
  }
  

udpateTimer();
// await pollServer();