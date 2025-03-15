let tableID = ""
let username = ""
let passcode = ""

function saveData() {
    
    localStorage.setItem("tableID", tableID);
    localStorage.setItem("passcode", passcode);
    localStorage.setItem("username", username);
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