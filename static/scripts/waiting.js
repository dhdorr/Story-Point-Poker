// let myInterval = setInterval(printMe, 1000)
// let count = 0;

// function TestMe(el) {
//     clearInterval(myInterval);
// }

// function printMe() {
//     count += 1;
//     console.log("doin stuff...", count);
// }

function TestMe(el) {
    let pc = document.getElementById("playerCountPoll");

    clearInterval(pc.__fixi.pollInterval);
}
