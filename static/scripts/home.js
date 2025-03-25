// use this to remove polling on an element by adding 'fx-delete="<id>"' to a separate element
document.addEventListener("fx:config", (evt) => {
    console.log("event: ", evt);

    if (evt.srcElement.hasAttribute("fx-delete")) {
        let tmp = document.getElementById(evt.srcElement.getAttribute("fx-delete"));
        // tmp.removeAttribute("ext-fx-poll-interval");
        console.log(tmp)
        clearInterval(tmp.__fixi.pollInterval);
    }
})