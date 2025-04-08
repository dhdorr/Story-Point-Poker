let update_interval = setInterval(GetUpdate, 1000);

function GetUpdate() {
    let temp = document.getElementById("gameOver");
    if (temp == null) {
        return
    }

    if (temp.getAttribute("value") == "true" || temp.getAttribute("value") == "1") {
        clearInterval(update_interval);
        let pc = document.getElementById("gameUpdatePoll");
        clearInterval(pc.__fixi.pollInterval);

        let btn = document.getElementById("resultsBtn");
        btn.click();
    }
}