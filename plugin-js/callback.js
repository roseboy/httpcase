var myDate = new Date()

function rwork(result) {
    var res = JSON.parse(result);
    var num = 0;
    var details = new Array();
    for (var i = 0; i < res.TestCases.length; i++) {
        var cas = res.TestCases[i];
        if (cas.IsGlobal) {
            continue;
        }
        num++;

        var msg = ""
        if (cas.Pass) {
            try {
                var obj = JSON.parse(cas.Response.Body)
                msg = "Msg:" + obj.msg;
            } catch (e) {

            }
        }
        details.push({
            "testcaseName": cas.Name,
            "requestUrl": cas.Request.Url,
            "totalNum": 1,
            "successNum": cas.Pass?1:0,
            "failNum": cas.Pass?0:1,
            "errorNum": 0,
            "successRate": cas.Pass?'100%':'0%',
            "apiAverageTime": cas.Time+'ms',
            "apiMinTime":  cas.Time+'ms',
            "apiMaxTime":  cas.Time+'ms',
            "responseCode": cas.Response.Status,
            "requestInfo": cas.Request.Body,
            "errorInfo": ""
        })
    }
    var ret = {
        "taskName": "ai党e家",
        "buildId": "",
        "totalNum": res.Total,
        "successNum": res.Passed,
        "failNum": res.Failed,
        "errorNum": 0,
        "successRate":( parseFloat(res.Passed) / parseFloat(res.Total)*100).toFixed(2)+"%",
        "useTime": res.Duration+"ms",
        "reportTime": now(),
        "details": details
    }
    console.log("\n\n\n\n*********************************************")
    console.log(JSON.stringify(ret))
    console.log("*********************************************\n\n\n\n")
    return JSON.stringify(ret);
}

function now() {
    var day = new Date();
    return day.getFullYear() + "-" + fullZero(day.getMonth() + 1) + "-" + fullZero(day.getDate())+
        " "+fullZero(day.getHours())+":"+fullZero(day.getMinutes())+":"+fullZero(day.getSeconds());
}

function fullZero(n){
    if(n<=9){
        return "0"+n
    }else{
        return n+""
    }
}

//微信通知
function WeChatNotify(result) {
    var res = JSON.parse(result)
    var title = "【" + res.Passed + "通过," + res.Failed + "失败】";
    var body = ""
    var num = 0;
    for (var i = 0; i < res.TestCases.length; i++) {
        var cas = res.TestCases[i];
        if (cas.IsGlobal) {
            continue;
        }
        num++;

        var msg = ""
        if (cas.Pass) {
            try {
                var obj = JSON.parse(cas.Response.Body)
                msg = "Msg:" + obj.msg;
            } catch (e) {

            }
        }
        body += ("\n[" + num + "] " + cas.Name + (cas.Pass ? "[通过]" : "[失败]") + "\n    " + msg)
    }

    return "?a=1&title=" + encodeURI(title) + "&desp=" + encodeURI(body)
}