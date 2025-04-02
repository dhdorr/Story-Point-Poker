let ready = document.getElementById("isReady");
let startBtn = document.getElementById("startGameBtn");

let ready_interval = setInterval(CheckIfReady, 1000);

function CheckIfReady() {
    ready = document.getElementById("isReady");
    console.log("is ready? ", ready);
    if (ready.getAttribute("value") == "1" || ready.getAttribute("value") == "true" ) {
        console.log("ready to start the game!");
        startBtn.click();
        // clearInterval(ready_interval);
        // let pc = document.getElementById("playerCountPoll");
        // clearInterval(pc.__fixi.pollInterval);
    }
}

function TestMe() {
    clearInterval(ready_interval);
    let pc = document.getElementById("playerCountPoll");
    clearInterval(pc.__fixi.pollInterval);
}
