let update_interval = setInterval(GetUpdate, 1000);

function GetUpdate() {
    let round_over = document.getElementById("gameOver");
    if (round_over != null) {
        if (round_over.getAttribute("value") == "true" || round_over.getAttribute("value") == "1") {
            // clearInterval(update_interval);
            let pc = document.getElementById("gameUpdatePoll");
            clearInterval(pc.__fixi.pollInterval);
    
            let btn = document.getElementById("resultsBtn");
            btn.click();
        }
    }

    let proceed_next_round = document.getElementById("nextRound");
    if (proceed_next_round != null) {
        if (proceed_next_round.getAttribute("value") == "true" || proceed_next_round.getAttribute("value") == "1") {
            // clearInterval(update_interval);
            let pc = document.getElementById("nextRoundPoll");
            clearInterval(pc.__fixi.pollInterval);
    
            let btn = document.getElementById("nextRoundBtn");
            btn.click();
        }
    }
}

function ProceedToNextRound() {
    let pc = document.getElementById("nextRoundPoll");
    clearInterval(pc.__fixi.pollInterval);
}