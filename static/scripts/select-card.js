let bg_skin = "grey";
let bg_text = "black";

function choose(card) {
    console.log("clicked...");
    if (card.classList.contains("selected")) {
        card.classList.remove("selected");
    } else {
        card.classList.add("selected");
    }
}

document.addEventListener("fx:config", (evt) => {
    evt.detail.cfg.headers.testme = "don't panic";
    console.log("After headers: " + evt.detail.cfg.body);
})

function logEvt(evt) {
    console.log(evt.type, evt.target, evt.detail?.cfg || "", evt.detail)

    console.log("my test: ",   evt.detail.cfg.response.status)
}

function handleErrorResponse(evt) {
    let error_arr = [401, 404];
    if (error_arr.includes(evt.detail.cfg.response.status)) {
        console.log("there was a 404 error...")
        // document.head.remove();
        let error_dialog = document.getElementById("errorDialog");
        evt.detail.cfg.target = error_dialog;
        evt.detail.cfg.swap = 'innerHTML';
    }
}

function handleRequest(evt) {
    evt.detail.cfg.headers.bg_skin = bg_skin;
    evt.detail.cfg.headers.bg_text = bg_text;
    console.log("testing before req: ", evt.detail.cfg);
}

document.addEventListener("fx:after", (evt) => handleErrorResponse(evt))
document.addEventListener("fx:before", (evt) => handleRequest(evt))

// let radios = document.getElementsByName("skin");

// radios.forEach(rad => {
//     rad.addEventListener("change", function() {
//         console.log("changing color to " + this.value);
//         if (this.checked) {
//             console.log("changing color to " + this.value);
//         }
//     })
// })

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