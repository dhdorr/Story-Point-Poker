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
    console.log("After headers: " + evt.detail.cfg.body.get("sessionID"));
})

function logEvt(evt) {
    console.log(evt.type, evt.target, evt.detail?.cfg || "", evt.detail)

    console.log("my test: ",   evt.detail.cfg.response.status)
}

function handleErrorResponse(evt) {
    if (evt.detail.cfg.response.status == 404) {
        console.log("there was a 404 error...")
        // document.head.remove();
        let error_dialog = document.getElementById("errorDialog");
        evt.detail.cfg.target = error_dialog;
        evt.detail.cfg.swap = 'innerHTML';
    }
}

document.addEventListener("fx:after", (evt) => handleErrorResponse(evt))



