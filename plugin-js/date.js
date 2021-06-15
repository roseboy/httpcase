
function addDate(year, month, day) {
    var date = new Date();
    date.setDate(date.getDate() + parseInt(day));
    date.setMonth(date.getMonth() + parseInt(month));
    date.setFullYear(date.getFullYear() + parseInt(year));
    var year = date.getFullYear();
    var month = change(date.getMonth() + 1);
    var day = change(date.getDate());
    var hour = change(date.getHours());
    var minute = change(date.getMinutes());
    var second = change(date.getSeconds());
    var time = year + '-' + month + '-' + day + ' ' + hour + ':' + minute + ':' + second;
    return time;
}

function min(date) {
    return date.substr(0, 16);
}

function day(date) {
    return date.substr(0, 10);
}

function month(date) {
    return date.substr(0, 7);
}

function year(date) {
    return date.substr(0, 4);
}

function getLocalTime() {
    var date = new Date();
    return date.getTime();
}


function redomdata() {
    var date = new Date();
    return date.getTime() % 1000000000;
}


parseInt(new Date().getTime() / 1000) * 1000

function change(t) {
    if (t < 10) {
        return "0" + t;
    } else {
        return t;
    }
}

function getTomorrowDate() {
    return addDate(0,0,1);
}

function test(v1,v2,v3,v4) {
    return "hello "+v1+" "+v2+" "+v3+" "+v4+" ";
}