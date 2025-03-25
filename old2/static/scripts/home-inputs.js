let tableID = ""
let username = ""
let passcode = ""

function saveData() {
    localStorage.setItem("tableID", tableID);
    localStorage.setItem("passcode", passcode);
    localStorage.setItem("username", username);
    localStorage.setItem(username + ':' + tableID, passcode);
}

function updateData(elm) {
    if (elm.name == "tableID") {
        tableID = elm.value;
        console.log(tableID);
    }
    if (elm.name == "passcode") {
        passcode = elm.value;
        console.log(passcode);
    }
    if (elm.name == "username") {
        username = elm.value;
        console.log(username);
    }
}